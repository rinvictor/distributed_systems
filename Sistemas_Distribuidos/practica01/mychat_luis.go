// Information and user guide:
// Enter name (a-z A-Z 0-9 _) (18 max). If duplicate, will ask again.
// Active clients information will be displayed every time a new client joins.
// Clients in private chats will be preceded by a #.
// To request a private chat enter >>:name where name is the name of
// the other client whom must be connected.
// Only one private chat is allowed at a time. To exit a private chat
// enter >>:#exit and you will return to the general chat.
// If the person you want to start a private session with is already
// in one you will be notified and will have to wait for it to finish.
// Private chats must be accepted by the other client.
// There is a 10 second span where it can be accepted, after which it
// will expire and thus be considered declined.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"regexp"
    "time"
)

//!+broadcaster
type client struct {
    ch chan<- string // an outgoing message channel
    name string
}
type privChat struct {
	client1 client
    client2 client
}

var (
	maxLen = 18
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
	request = make(chan string) // request and accept private chats
)

func broadcaster() {
	clients := make(map[client]bool) // true if not in a private chat
	private := make(map[privChat]bool) // private channels. true if accepted
	for {
		sel:
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			nick := strings.Split(msg, ":")[0]
			// check if client is in a private channel
			for p, accepted := range private {
				if (nick == p.client1.name) && accepted {
					for cli := range clients {
						if cli.name == p.client2.name {
							cli.ch <- msg
							break sel
						}
					}
				}else if (nick == p.client2.name) && accepted {
					for cli := range clients {
						if cli.name == p.client1.name {
							cli.ch <- msg
							break sel
						}
					}
				}
			}
			// skip whom has sent the message and private channels
			for cli := range clients {
                found := false
				for p, accepted := range private {
					if ((cli.name == p.client1.name) || (cli.name == p.client2.name)) && accepted {
                        found = true
					}
				}
                if (cli.name != nick) && !found {
                    cli.ch <- msg
                }
			}
		case msg := <-request:
			// check if request or accept
			nick := strings.Split(msg, "-")[0]
			nickDest := strings.Split(msg, "-")[1]
            if nickDest == "#exit" {
                // remove entry from private
                for p := range private {
                    if p.client1.name == nick {
                        nickDest = p.client2.name
                        delete(private, p)
                    }else if p.client2.name == nick {
                        nickDest = p.client1.name
                        delete(private, p)
                    }
                }
                for cli := range clients {
                    if (cli.name == nick) || (cli.name == nickDest) {
                        clients[cli] = true // back to general
                        cli.ch <- ">>Private channel closed. Returned to general."
                    }
                }
                break sel
            }
			for p, accepted := range private {
				if (nick == p.client2.name) && (nickDest == p.client1.name) && (!accepted) {
					private[p] = true
                    notify := ">>Private channel created (" + nick + "-" + nickDest + ").\n"
                    notify += ">>Enter [>>:#exit] to go back to general."
                    p.client1.ch <- notify
                    p.client2.ch <- notify
					// uncheck them from general
                    for cli := range clients {
                        if (cli.name == nick) || (cli.name == nickDest) {
                            clients[cli] = false
                        }
                    }
					break sel
				}
			}
			// check that dest exists
			found := false
			var dest chan<- string
			for cli := range clients {
				if cli.name == nickDest {
					dest = cli.ch
					found = true
					break
				}
			}
            // if theres no such client: inform requester
			if !found {
                abort := ">>No such client is active at the moment."
				for cli := range clients {
					if cli.name == nick {
						cli.ch <- abort
						break
					}
				}
				break
			}
			//check if client or dest is already in a private chat
			for p := range private {
				if (p.client1.name == nick) || (p.client2.name == nick) {
					for cli := range clients {
						if cli.name == nick {
							warn := "(!) Cannot be in more than 1 chat at a time.\n"
							warn += "Exit current chat first: (>>:#exit)"
							cli.ch <- warn
							break
						}
					}
					break sel
				}else if (p.client1.name == nickDest) || (p.client2.name == nickDest) {
					for cli := range clients {
						if cli.name == nick {
							warn := "(!) That person is already in a private chat.\n"
							warn += "Try again later."
							cli.ch <- warn
							break
						}
					}
					break sel
				}
			}
            // build private chat structure
            var nickch chan<- string
            for cli := range clients {
                if cli.name == nick {
                    nickch = cli.ch
                }
            }
            c1 := client{ch: nickch, name: nick}
            c2 := client{ch: dest, name: nickDest}
            pchat := privChat{client1: c1, client2: c2}
			private[pchat] = false // not yet accepted by the other person

            // ask dest to accept private channel
			ask := ">>:" + nick + " wants to start a private chat with you.\n"
			ask += ">>Enter [>>:"+nick+"] to accept (will expire in 10 seconds): "
			dest <- ask // until dest accepts, keep working!
            // wait 10s for reply.
            go func(){
                time.Sleep(10 * time.Second)
                if !private[pchat] {
                    delete(private, pchat)
                    // inform clients of expiration
                    expired := "(!) Request expired."
                    pchat.client1.ch <- expired
                    pchat.client2.ch <- expired
                }
            }()

		case cli := <-entering:
			// save client name and broadcast client list
            for c := range clients {
                if c.name == cli.name {
                    cli.ch <- "yes" // duplicate
                    break sel
                }
            }
            cli.ch <- "no"
            clients[cli] = true
			list := "-Clients currently active: "
			for c, general := range clients {
                if !general {
                    list += "#" // meaning that client is in a private chat
                }
				list += c.name + ", "
			}
			for c, general := range clients {
                if general { // send only to those in general chat
                    c.ch <- list
                }
			}
		case cli := <-leaving:
            // delete private chat
			for p := range private {
				if (cli.name == p.client1.name) || (cli.name == p.client2.name) {
					delete(private, p)
				}
			}
            // delete from client list
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages

    // get name and notify broadcaster
	input := bufio.NewScanner(conn)
    who := ""
    var cli client
    for {
        who = getName(conn, *input)
        cli = client{ch: ch, name: who}
        entering <- cli
        dup := <- ch
        if dup == "no" {
            break
        }
        fmt.Fprintln(conn, "(!) That name is taken. Try a different one.")
    }

    go clientWriter(conn, ch)
    ch <- "You are " + who
	messages <- who + ": just arrived!"
    // read and send until client disconnects
	text := ""
	for input.Scan() {
		// check for private chat indicator [>>:]
		text = input.Text()
		if strings.Split(text, ":")[0] == ">>" {
			request <- who + "-" + strings.Split(text, ":")[1]
		}else{
			messages <- who + ": " + text
		}
	}
	// NOTE: ignoring potential errors from input.Err()
    // client has disconected
	leaving <- cli
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

func getName(conn net.Conn, input bufio.Scanner) string{
	// accepted characters for a username:
	isValid := regexp.MustCompile(`^[A-Za-z0-9_]+$`).MatchString
	for {
		fmt.Fprintf(conn, "Please enter a valid name: ")
		input.Scan()

		if len(input.Text()) > 18 {
			fmt.Fprintf(conn, "(!) Name is too long (18 characters max.)\n")
		}else if isValid(input.Text()) {
			break
		}else{
			fmt.Fprintf(conn, "(!) Invalid character.")
			fmt.Fprintf(conn, " Only a-z, A-Z, 0-9, _ are accepted.\n")
		}
	}
	return input.Text()
}

//!-handleConn

//!+main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

//!-main

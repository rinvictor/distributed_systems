// Víctor Rincón Yepes

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type client chan<- string // An outgoing message channel

var (
	entering  = make(chan client)
	leaving   = make(chan client)
	messages  = make(chan string) // All incoming client messages
	usernames = make(chan string) // A channel for usernames
)

func compareIP(v string, msg string) bool {
	if strings.Split(v, "-")[1] == strings.Split(msg, "-")[1] {
		return false
	}
	return true
}

func addUsers(clients map[client]string) string {
	count := 0
	list := "USUARIOS ACTUALES CON IP: \n"
	for _, value := range clients {
		count = count + 1
		if len(clients) > count { // Cleaning the output
			list += "\t" + value + "\n"
		} else {
			list += "\t" + value
		}
	}
	return list
}

//!+broadcaster
func broadcaster() {
	clients := make(map[client]string) // Active clients
	for {
		select {
		case msg := <-messages: // Receiving messages
			if strings.Index(msg, "!list") > -1 {
				list := addUsers(clients) // Adding all usernames+ip to a string

				for cli, v := range clients {
					if !compareIP(v, msg) {
						cli <- list
					}
				}

			} else {
				for cli, v := range clients { // Not for the sender
					if compareIP(v, msg) {
						cli <- strings.Split(msg, "-")[0]
					}
				}
			}

		case cli := <-entering: // New client
			clients[cli] = <-usernames // Receiving usernames
			list := addUsers(clients)  // Adding all usernames+ip to a string

			for c, _ := range clients {
				c <- list
			}

		case cli := <-leaving: // Client leaving
			delete(clients, cli)
			close(cli)
		}
	}
}

//!-broadcaster

func isValid(s string) bool { // '-' character is forbidden
	i := strings.Index(s, "-")
	if i > -1 {
		return false
	}
	return true
}

func getName(conn net.Conn) string {
	inputname := bufio.NewScanner(conn)

	for {
		fmt.Fprintf(conn, "Introduce tu nombre: ")
		inputname.Scan()

		if isValid(inputname.Text()) {
			break
		} else {
			fmt.Fprintf(conn, "Caracteres no validos.\n")
		}
	}
	return inputname.Text()
}

//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string) // Sending all messages

	go clientWriter(conn, ch)

	name := getName(conn)

	who := conn.RemoteAddr().String() // ip to string
	welcome := "comandos:\n\t!exit --> abandonar el chat\n\t!list --> lista los usuarios conectados"
	ch <- "Eres " + who + ", con nick: " + name
	ch <- welcome
	messages <- name + " se ha unido" + "-" + who
	entering <- ch
	usernames <- name + "-" + who // Channel for usernames

	// Sending messages
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- name + ": " + input.Text() + "-" + who // Canal de lo mensajes
		if input.Text() == "!exit" {
			break
		}
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- ch // Client leaving
	messages <- name + " abandona el chat" + "-" + who
	conn.Close()
}

//!-handleConn

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors, para cada mensaje imprimo conn y el mensaje
	}
}

//!+main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000") // Listening on port 8000
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

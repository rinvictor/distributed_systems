// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

//!+broadcaster
type client chan<- string // an outgoing message channel

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
	//necesito otro canal para poder pasar la info en entering, porque es es de tipo cliente
	usernames = make(chan string)
)

func CompareIP(v string, msg string) bool{
	if strings.Split(v, "-")[1] == strings.Split(msg, "-")[1]{
		return false
	}
	return true	
}

func broadcaster() {
	clients := make(map[client]string) // mapa de todos los clientes coenctados
	for {
		select {
		case msg := <-messages: //recupero el mensaje que va para todos los clientes
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			
			//fmt.Println(msg)
			//ip := strings.Split(msg, "-")[1]
			//fmt.Println(ip)
			//fmt.Println(msg)
			//fmt.Println(strings.Split(msg, "-")[0])
			for cli, v := range clients { //le paso el mensaje a todos los clientes
				if CompareIP(v, msg){ //No me lo envio a mi mimso
					cli <- strings.Split(msg, "-")[0] //Este canal envia los mansajes a todos los clientes incluido el mismo
				}
				
				
				/*if v == strings.Split(msg, "-")[1]{
					fmt.Println("ww", v, strings.Split(msg, "-")[1])
				}*/
			}
		
		case cli := <-entering: //aqui digo quien es quien
			clients[cli] = <- usernames
			//fmt.Println(cli)
			//escribo los clientes que están unidos
			
			cli<-"USUARIOS ACTUALES CON IP: "
			for _, v := range clients {
				cli <- v 
			}
			


 



			//for cli := range clients { //digo que un nuevo cliente se ha conectado
				//continue //Esto se lo envío a todos los clientes incluido el mismo
			//}
		
		case cli := <-leaving: //aqui dices que clientes se han ido
			delete(clients, cli) //quito a los clientes
			close(cli) //cierro el canal
		}
	}
}

//!-broadcaster
func getName(conn net.Conn) string{
	
	inputname := bufio.NewScanner(conn)
	inputname.Scan()
	name := inputname.Text()
	return name
}
//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string) //mensajes de clientes de salida

	go clientWriter(conn,ch)

	//preguntar quien es
	ch <- "Nombre: "
 	name := getName(conn)

	who := conn.RemoteAddr().String() //ip del cliente en string
	ch <- "Eres " + who + ", con nick: " + name
	messages <- name + " se ha unido" +  "-" + who 
	//entering <- ch //por este canal le paso quien dice que
	entering <- ch
	usernames<-name + "-" + who

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- name + ": " + input.Text() + "-" + who //canal de lo mensajes
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- ch //canal para cuando alguien se va del chat
	messages <- name + " abandona el chat" + "-" + who
	conn.Close() //cierro la conexion
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors, para cada mensaje imprimo conn y el mensaje
	}
}

//!-handleConn

//!+main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000") //crea un servidor
	if err != nil { //manejo el error
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept() //servidor acepta un cliente
		if err != nil { //manejo el error
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

//!-main

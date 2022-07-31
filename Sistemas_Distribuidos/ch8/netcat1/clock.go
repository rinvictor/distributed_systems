package main

import (
	"io"
	"log"
	"net"
	"time"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1{
		fmt.Println("usage: clock [port] [timezone]")
	}
	port := os.Args[1]
	port = "localhost:"+port

	listener, err := net.Listen("tcp",port) //
	//manejar error
	if err != nil {
		log.Fatal(err)
	}
	//
	for {
		conn, err := listener.Accept()//

		if err != nil {
			log.Print(err) // e.g., connection aborted, manejo de error
			continue
		}

		go handleConn(conn) //una conexion
	}
}

func handleConn(c net.Conn) {
	defer c.Close() //defer se usa para asegurarte que va lo ultimo en el orden de ejecución
	for {
		loc, _ := time.LoadLocation(os.Args[2])
		t:=time.Now().In(loc)
		_, err := io.WriteString(c, t.Format("15:04:05\n"))

		//manejo de error
		if err != nil {
			fmt.Println("disconnected")
			return // e.g., client disconnected
		}

		time.Sleep(1 * time.Second)
	}
}

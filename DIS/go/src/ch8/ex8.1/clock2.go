// Clock is a TCP server that periodically writes the time.
package main

import (
	"fmt"
	"io"
	"flag"
	"log"
	"net"
	"time"
	"strconv"
)

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	port := flag.Int("port", 8000, "port to listen at")
	flag.Parse()

	listener, err := net.Listen("tcp", "localhost:" + strconv.Itoa(*port))
	if err != nil {
		log.Fatal(err)
	}
	//!+
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn) // handle connections concurrently
	}
	//!-
}

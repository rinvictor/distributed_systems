// ClockWall acts as a client for different clock servers at once.

package main

import (
    "fmt"
    "io"
	"log"
	"net"
    "os"
    "strings"
    "strconv"
)

var host = "localhost"

func main() {
    // leo y proceso los argumentos
    args := os.Args[1:]
    timetable := make([][]string, 0)
    for i := range args {
        entry := strings.Split(args[i], "=")
        timetable = append(timetable, entry)
    }
    //necesito un canal para esperar a las gorutinas
    c := make(chan struct{})
    // pensaba mostrar las horas en columnas pero ha resultado inviable
    for i := range args {
        fmt.Print(timetable[i][0] + "\t\t")
    }
    fmt.Println()
    //lanzo una gorutina por cada servidor con el que conectarse
    for i := range args {
        p, err := strconv.Atoi(timetable[i][1])
        if err != nil {
            log.Fatal(err)
        }
        go getTime(host, p, c)
    }
    <-c //espero a las subrutinas (nunca van a volver)
}

func getTime(h string, p int, c chan struct{}) {
    conn, err := net.Dial("tcp", h + ":" + strconv.Itoa(p))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(os.Stdout, conn)
    c<-struct{}{} //aqui ni va a llegar
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

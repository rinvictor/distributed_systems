package main

import(
  "fmt"
  "sync"
  "math/rand"
  "time"
)

const Nrations = 5
const NCannibals = 9

var explorers = 7

var mu sync.Mutex //Mutex general
var rations int //Variable compartida

var vacia sync.Mutex
var llena sync.Mutex

func eat(i int){
    mu.Lock()
    if(rations == 0){
        fmt.Println("Vaya!, está vacía...", i)
        vacia.Unlock()
        fmt.Println("0") //aqui llega
        llena.Lock()
        fmt.Println("1") //aqui no, se bloquea
        //PROBLEMA: con esto no llamo al cocinero

    }else{
        time.Sleep(time.Duration(rand.Float32() + 1.0) * time.Second)
        fmt.Println("Estoy comiendo...", i)
        rations--
        fmt.Println("Quedan ", rations, " raciones")
    }
    mu.Unlock()
}

func work(i int){
    time.Sleep(time.Duration(rand.Float32() + 1.0) * time.Second)
    fmt.Println("Estoy trabajando...", i)
}


func sleep(){
    fmt.Println("Estoy durmiendo...")
}

func cannibal(i int){
    eat(i)
    work(i)
}

func cook(){
    mu.Lock()
    fmt.Println("Estoy cocinando...")
    for i:=0; i<Nrations; i++{
        time.Sleep(time.Duration(rand.Float32() + 1.0) * time.Second)
        rations++
        fmt.Println("Una ración más! llevo ",rations)
    }

    explorers--
    fmt.Println("Quedan ", explorers," exploradores")
    mu.Unlock()
}

func cooker(){
    vacia.Lock()//cuando el cocinero es llamado es porque está vacía
    fmt.Println("cocinero llamado")
    cook()
    if (explorers == 0){
        fmt.Println("FIN")
        //aviso que se acaba el programa
    }
    llena.Lock() //ya está llena
    vacia.Lock()
}


func main(){
    go cooker()
    time.Sleep(1.0 * time.Second)

    for i := 0; i < NCannibals; i++ {
        go cannibal(i)
    }

    //Implementar espera
    time.Sleep(time.Second * 20)



}

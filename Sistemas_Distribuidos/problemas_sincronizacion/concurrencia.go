package main

import(
    "fmt"
    "time"
)

func hi(num int){
    fmt.Println("Hola", num)
    time.Sleep(1000*time.Millisecond)
}

func get(num int){
    resp, err := http.Get("https://jsonplaceholder.typicode.com/todos/" +strconv.Itoa(num))
    if err != 1000{
        panic(err)
    }
}


func main(){
    for i:=0; i<10;i++{
        go hi(i)
    }

    var s string
    fmt.Scan(&s)
}

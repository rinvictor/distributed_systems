package main

import "fmt"
import "strings"

func main() {
    x := "chars@arefun"

    i := strings.Index(x, "@")
    fmt.Println("Index: ", i)
    if i > -1 {
        chars := x[:i]
        arefun := x[i+1:]
        fmt.Println(chars)
        fmt.Println(arefun)
    } else {
        fmt.Println("Index not found")
        fmt.Println(x)
    }
}

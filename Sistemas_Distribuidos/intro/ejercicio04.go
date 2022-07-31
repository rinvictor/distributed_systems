package main

import "fmt"

func main(){
  //Creo un mapa
  m := make(map[string] string)
  m["Spain"] = "Madrid"
  m["Italy"] = "Rome"
  m["EEUU"] = "Washington DC"
  m["Norway"] = "Oslo"
  fmt.Println(m)


  //Itero usando range
  fmt.Println("Iterando con range...")
  for k, v := range m{
    fmt.Printf("%s -> %s\n", k, v)
  }
  //Eliminar entradas
  delete(m, "Spain")
  fmt.Println("Tras eliminar: ", m)

}

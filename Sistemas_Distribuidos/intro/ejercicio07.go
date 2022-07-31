package main

import (
  "fmt"
  "math"
)
const Pi = 3.141592

type circle struct{
  radio float64
}


func main(){
  var r_user float64
  fmt.Print("Intro radio: ")
  fmt.Scanln(&r_user)
  c1 := circle{radio:r_user}
  fmt.Println("area: ", area(c1))
}

func area(c circle) float64{
  return Pi*math.Pow(c.radio,2)
}

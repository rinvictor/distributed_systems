package main

import "fmt"

func main(){
  var myarray[10] int
  myarray=initArray(myarray,10)
  fmt.Println(myarray)
  var mean = mean(myarray)
  fmt.Println("mean: ", mean)

  var otroarray[5] int
  otroarray[0] = 200
  otroarray[1] = 4
  otroarray[2] = 6
  otroarray[3] = 2
  otroarray[4] = 15
  fmt.Println("Iterando con range...")
  for i := range otroarray {
      fmt.Println("index:", i, "->", otroarray[i])
  }

  var max = max (otroarray)
  fmt.Println("max: ", max)

  var min = min(otroarray)
  fmt.Println("min: ", min)

  var position = position(otroarray, 200)
  fmt.Println("pos: ", position)
}

func printArray(myarray [10] int){
  fmt.Println(myarray)
}

func mean(myarray [10] int) int{
  var lenght = len(myarray)
  var sum = 0
  for i:=0;i<lenght;i++{
    sum = sum + myarray[i]
  }
  return sum/lenght
}

func max(myarray[5] int) int{
  var lenght = len(myarray)
  var max = 0
  for i:=0;i<lenght;i++{
    if myarray[i] > max {
      max = myarray[i]
    }
  }
  return max
}

func min(myarray[5] int) int{
  var lenght = len(myarray)
  var min = myarray[0]
  for i:=0;i<lenght;i++{
    if myarray[i] < min {
      min = myarray[i]
    }
  }
  return min
}

func position(myarray[5] int, n int) int{ //Funcion que dado un valor devuelve la posicion del array
  var pos = -1
  for i:=0;i<len(myarray);i++{
    if myarray[i] == n{
      pos = i
    }
  }

  return pos
}

func initArray(myarray [10] int, initvalue int) [10] int{ //No se por que no puedo pasarle al array npos directamente
  var npos = len(myarray)
  for i:=0;i<npos;i++{
    myarray[i]=initvalue
  }
  //printArray(myarray)
  return myarray
}

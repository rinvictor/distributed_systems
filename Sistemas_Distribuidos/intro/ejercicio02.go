package main
import "fmt"

func main(){
  var myslice = make([]int,2,30) //Esto es un array de capacidad 30 pero que de inicio solo usa dos posiciones, su len es de 2
  sliceInfo(myslice)
  myslice[0]=0
  myslice[1]=1
  sliceInfo(myslice)
  myslice = append(myslice, 2, 3) //Asi se añaden numeros una vez hemos superado la longitud asignada al principio
  sliceInfo(myslice)

  var myslice2 = make([]int,len(myslice),60) //DUDA: si pongo 2 como longitud inicial: [0 1], solo copia los dos primeros, si pongo 4: [0 1 2 3]

  sliceInfo(myslice2)
  copy(myslice2, myslice)
  sliceInfo(myslice2)
  fmt.Println("Iterando con range...")
  for _,myslice2 := range myslice2{
    fmt.Println(myslice2)
  }
}

func sliceInfo(myslice []int){
  fmt.Println(myslice)
  fmt.Println("capacity: ", cap(myslice))
  fmt.Println("lenght: ", len(myslice))
}

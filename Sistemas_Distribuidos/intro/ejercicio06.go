package main
import "fmt"

type book struct{
  title string
  author string
  edition string
  year int
  ISBN int
}

func main(){
  b1 := book{title:"title1", author: "author1"}
  b1.edition = "edition1"
  b1.year = 2020
  b1.ISBN = 1234567890
  printBook(&b1)

  b2 := book{}
  b2.title = "title2"
  b2.author = "author2"
  b2.edition = "edition2"
  b2.year = 2021
  b2.ISBN = 1234567891
  printBook(&b2)

}

func printBook(mybook *book){ //Paso por referencia
  fmt.Println(*mybook)
}

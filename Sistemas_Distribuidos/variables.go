package main
import "fmt"

func main() {
	
	var a = "initial"
	fmt.Println(a)
	
	var b, c = 1, 2
	fmt.Println(b, c)
	
	var d = true
	fmt.Println(d)
	
	var e int /*Por defecto 0*/
	fmt.Println(e)
	
	f := "apple" /*No es recomendable*/
	fmt.Println(f)
}

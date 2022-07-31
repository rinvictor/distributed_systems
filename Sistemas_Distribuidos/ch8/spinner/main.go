// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 218.

// Spinner displays an animation while computing the 45th Fibonacci number.
package main

import (
	"fmt"
	"time"
)

//!+
//Se lanza la gouroutine que se ejecuta hasta que termine lo que va a continuación
func main() {
	go spinner(100 * time.Millisecond) //crea una nueva gouroutine
	const n = 20
	fibN := fib(n) // slow, calculo los n terminos de fibonacci
	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
}

//Esta funcion imprime esos valores que hacen parecer que da vueltas con ese delay, un numero más pequeño da la impresion de ir mas rapido
func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

//calculo de la sucesion de Fibonacci del n valor
func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}

//!-

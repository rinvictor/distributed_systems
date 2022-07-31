package main

import (
	"math/rand"
	"fmt"
	"time"
)

var maxRenos = 9
var maxElfos = 12
var minElfos = 3
var minRenos = 9
var tRenos = 5

//Funcion que define el funcionamiento de Santa
func santa(problemasE chan bool, elfosListos chan bool, renosVuelta chan bool, renosListos chan bool, santaInicio chan bool) {
	fmt.Println("Iniciando a Santa!")

	elfos := 0
	renos := 0

	//Ponemos a Santa a trabajar con un bucle infinito
	for {
		select {
		case <-renosVuelta:
			renos++
			if renos == minRenos {
				fmt.Println("Santa está preparando el trineo!")
				for i := 0; i < renos; i++ {
					renosListos <- true
					<-santaInicio
				}
				fmt.Println("santa está repartiendo regalos!")
				renos = 0
			}
		case <-problemasE:
			elfos++
			if elfos == minElfos {
				//Calculamos el tiempo de espera de Santa y esperamos
				tiempoAyuda := time.Duration(rand.Intn(4) + 2)
				time.Sleep(tiempoAyuda * time.Second)
				fmt.Println("Santa esta ayudando a los elfos!")
				for i := 0; i < elfos; i++ {
					elfosListos <- true
				}
				elfos = 0
			}

		}
	}

}

func elfo(n int, problemasE chan bool, elfosListos chan bool) {
	for {
		//Genero un numero entre 0 y 2 si es 0, fallo
		problemas := rand.Intn(3) == 0
		if problemas{
			fmt.Println("Elfo número ", n, "esperando la ayuda de Santa!")
			//Lo añadimos al canal de problemas
			problemasE <-  true
			<-elfosListos //Esperamos a que Santa acabe
			fmt.Println("Elfo número ", n, "vuelve a la carga!")
		}
		fmt.Println("Elfo número ", n, "trabajando!")
		tiempoAyuda := time.Duration(rand.Intn(4) + 2)
		time.Sleep(tiempoAyuda * time.Second)

	}
}

func reno(n int, renosVuelta chan bool, renosListos chan bool, santaInicio chan bool)  {
	for{
		
		renosVuelta <- true
		//Esperamos a que tengamos el mínimo de renos
		<-renosListos
		//Nos sincronizamos con Santa
		santaInicio <- true
		fmt.Println("Reno número ", n, "repartiendo regalos!")
		tiempoEspera := time.Duration(5 + rand.Intn(4))
		time.Sleep(tiempoEspera* time.Second) //Esperamos de 5 a 8 segundos para lanzar el reno
	}
}

func main() {

	santaInicio := make(chan bool)
	problemasE := make(chan bool)
	elfosListos := make(chan bool)
	renosVuelta := make(chan bool)
	renosListos := make(chan bool)

	go santa(problemasE, elfosListos, renosVuelta, renosListos, santaInicio)

	//Creamos las gorutinas necesarias para los elfos
	for i := 0; i < maxElfos; i++ {
		fmt.Println("Creando elfo ", i)
		go elfo(i, problemasE, elfosListos)
	}

	//Creamos las gorutinas necesarias para los renos

	for i := 0; i < maxRenos; i++ {
		tiempoEspera := time.Duration(5 + rand.Intn(4))
		time.Sleep(tiempoEspera* time.Second) //Esperamos de 5 a 8 segundos para lanzar el reno
		fmt.Println("Creando reno ", i)
		go reno(i, renosVuelta, renosListos, santaInicio)
	}
	for{}
}

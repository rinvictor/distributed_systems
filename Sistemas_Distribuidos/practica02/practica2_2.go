//Víctor Rincón Yepes

package main

import (
	"fmt"
	"time"
	"sync"
	"math/rand"
	"strconv"
)

	var santasBed sync.Mutex //Mutex para despertar a Santa.
	var factory sync.Mutex //Fábrica de juguetes. Controla nelvesintrouble.
	var nelvesintrouble = 0 //Número de elfos que estan esperando por la ayuda de santa.
	var elfWait sync.Mutex //Esperan a que Santa los ayude.
	var sledge sync.Mutex //Trineo. Mutex para nreindeers
	var nreindeers = 0 //Variable compartida que representa los renos que han llegado
	var wg sync.WaitGroup //Esperar a que acaben todas las gorutinas lanzadas

func santa(R int){
	//Santa Claus invierte un tiempo aleatorio de entre 2 y 5 
	//segundos en ayudar a un equipo de 3 elfos con problemas
	mintime := 2.0
	maxtime := 5.0
	for{
		//Santa duerme hasta que esten los renos o deba ayudar
		fmt.Println("Santa duerme...")
		santasBed.Lock()
		fmt.Println("Santa: ¡Alguien me ha despertado!")
		sledge.Lock()
		//alguien le despierta, los renos tienen prioridad
		if nreindeers == R {
			fmt.Println("Santa: Ya están todos los renos aquí, ¡creía que no llegarían nunca!")
			break
		}else{
			fmt.Println("Santa: Voy a ayudar a estos duendecillos")
			time.Sleep(time.Duration(rand.Float64()*(maxtime-mintime)+mintime)* time.Second) //aleatorio entre 2 y 5 seg
			elfWait.Unlock() //termino de ayudarlos
		}
		sledge.Unlock()

	}
	wg.Done()

}

func elf(){
	mintime := 3.0 //tiempo para fabricar un juguete
	troublerate := 3 //uno de cada 3 juguetes tienen problemas
	minelvesintrouble := 3 //numero de elfos para despertar a santa

	time.Sleep(time.Duration(rand.Float64()*60.0 + mintime) * time.Second)
	//uno de cada 3 juguetes dara  problemas
	if rand.Intn(troublerate) == 0 {
		
		fmt.Println("Elfo: Tengo problemas con este juguete")
		factory.Lock() //paro para comprobar
		nelvesintrouble++
		if nelvesintrouble == minelvesintrouble{
			fmt.Println("Elfo: A despertar al jefe para que nos ayude")
			santasBed.Unlock() //Despierto
			elfWait.Lock() //esperan a que termine de ayudar a estos
			//reseteo el conteo
			nelvesintrouble = 0
		}
		factory.Unlock()

	}else{
		fmt.Println("juguete terminado")
	}
	wg.Done()
}

func reindeer(R int){
	sledge.Lock() //bloqueo para comprobar
	//compruebo cuantos renos han llegado
	nreindeers++
	if nreindeers == R{
		fmt.Println("Rudolf: Soy el reno que faltaba, perdón por llegar tarde...¡A repartir!")
		santasBed.Unlock()
	}else{
		//nreindeers++
		fmt.Println("un reno llega, van "+strconv.Itoa(nreindeers))
	}
	sledge.Unlock()
	wg.Done()
}

func launchElves(D int) {
	for i := 0; i < D; i++ {
		go elf()
	}
}

func launchReindeers(R int) {
	mintime := 5.0 //tiempo minimo entre llegada de renos
	
	for i := 0; i < R; i++ {
		time.Sleep(time.Duration(rand.Float64()*2.0 + mintime) * time.Second)
		go reindeer(R)
	}
}

func main() {
	D := 12 //número de elfos
	R := 9 //número de renos
	wg.Add(D+R+1) //número de gorutinas
	rand.Seed(time.Now().UTC().UnixNano()) //aleatoreidad, necesario, si no no serían 
											//los mismos números siempre en las funciones rand

	//init
	santasBed.Lock()
	elfWait.Lock()

	go santa(R)
	go launchReindeers(R)
	go launchElves(D)

	wg.Wait() //Espero que acaben todas las gorutinas lanzadas
	fmt.Println("Fin. ¡Feliz Navidad!")
}

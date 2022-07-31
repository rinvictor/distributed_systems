//Víctor Rincón Yepes

package main

import (
	"fmt"
	"time"
	"sync"
	"math/rand"
	"strconv"
)

const (
	D = 15 //número de elfos, sin el lider
	R = 9 //número de renos
	NBattalions = 2 //numero de batallones
	Rows = 3
	Spaces = 30
)

	var santasBed sync.Mutex //Mutex para despertar a Santa.
	var wait sync.Mutex
	var factory sync.Mutex //Fábrica de juguetes. Controla nelvesintrouble.
	var nelvesintrouble = 0 //Número de elfos que estan esperando por la ayuda de santa.
	var elfWait sync.Mutex //Esperan a que Santa los ayude.
	var sledge sync.Mutex //Trineo. Mutex para nreindeers
	var nreindeers = 0 //Variable compartida que representa los renos que han llegado
	var wg sync.WaitGroup //Esperar a que acaben todas las gorutinas lanzadas
	var ntoysmade = 0
	var note = 0 
	var storeMu sync.Mutex //Mutex para protejer el almacén y las cachés

func santa(reindeersChan chan int, elvesChan chan int, leadersChan chan int){
	//Santa Claus invierte un tiempo aleatorio de entre 2 y 5 
	//segundos en ayudar a un equipo de 3 elfos con problemas
	mintime := 2.0
	maxtime := 5.0
	
	Loop:
	for{
		fmt.Println("Santa duerme...")
		//Recibe una llamada por uno de los canales
		select {

		case <- reindeersChan:
			fmt.Println("Santa: Ya están todos los renos aquí, ¡creía que no llegarían nunca!")
			fmt.Println("Santa: Antes de irme dejaré una nota para los elfos...")
			note = 1 //protected by santasRoom in reindeer
			wait.Unlock()
			break Loop
			
		case <- elvesChan:
			//Ayudar a los elfos
			fmt.Println("Santa: Voy a ayudar a estos duendecillos")
			time.Sleep(time.Duration(rand.Float64()*(maxtime-mintime)+mintime)* time.Second) //aleatorio entre 2 y 5 seg
			wait.Unlock()
		case <- leadersChan:
			fmt.Println("Santa: Ya habéis terminado los regalos, genial, este año tendréis cesta de Navidad.")
			break Loop
		}
	}
	wg.Done()
}

func initToys(toy *[Rows]bool) {
	for i := 0; i < Rows; i++ {
		(*toy)[i] = false
	}
}

func printoys(toy *[Rows]bool){
	for i := 0; i < Rows; i++ {
		fmt.Println((*toy)[i])
	}

}



func leader(tag int, caches *[NBattalions][Rows][Spaces]bool,
	toyStore *[Rows][Spaces]bool, elves, end chan int) {

	//coffeeBreak := 1 //coffee break time (seconds)
	full := 0 //number of full shelfs. when full==shelves: stop
	var reg [Rows]bool //registry for shelves already chosen

	//choose present, update caches, send present to elves.
	Loop:
	//for ; full < Rows; {
	for {
		initToys(&reg)
		for i := 0; i < Rows; {
			if note == 1 { //stop working and signal elves to finish
				fmt.Println("LíderB" + strconv.Itoa(tag) +
				": Santa se ha ido, terminad lo que estéis haciendo y a casa")
				break Loop
			}
			row := rand.Intn(Rows) //choose a random shelf
			storeMu.Lock()
			//fmt.Println(row)
			if !reg[row] { //if that shelf has not been chosen yet
				//fmt.Print("def: ")
				//fmt.Println(row)
				reg[row] = true //mark it as chosen
				n := rand.Intn(Spaces) //choose a random present
				for j := 0; j < Spaces; j++{ //Hasta que un regalo que aún esté disponible se elija 
					pos := (n + j) % Spaces //Paso al siguiente cuando esté vacío
					if !(*caches)[tag][row][pos] { //Compruebo en mi caché si el regalo está disponible
						toyStore[row][pos] = false //Lo retiro del almacén
						//Actualizo ambas cachés
						for k := 0; k < NBattalions; k++ {
							(*caches)[k][row][pos] = true
						}
						fmt.Println("LíderB" + strconv.Itoa(tag) + ": ¡Va un regalo!")
						elves <- 1 //give present to elves
						break
					}else if j == (Spaces-1){
						full++ //shelf is full
					}
				}
				i++ //misses do not count
			}
			storeMu.Unlock()
			//Paramos un poco la ejecución
			time.Sleep(time.Duration(1) * time.Second)
		}

		
		//condicion de finalización
		if full >= Rows{
			break
		}
		
		
	}

	for i := 0; i < D; i++ {
		elves <- 0 //tell elves that we have finished
	}
	santasBed.Lock()
	if note == 0 {
		note = 2 //let reindeers know that we've finished
		end <- 0 //signal Santa that we've finished
	}
	santasBed.Unlock()
	wg.Done()
}

func elf(toysChan chan int, elvesChan chan int){
	mintime := 3.0 //tiempo para fabricar un juguete
	troublerate := 3 //uno de cada 3 juguetes tienen problemas
	minelvesintrouble := 3 //numero de elfos para despertar a santa


	for {
		if (<-toysChan) == 0{ //el lider dice que no quedan juguetes
			break
		}
		fmt.Println("Elfo: ¡Ya tengo trabajo!, recibí un juguete.")
		time.Sleep(time.Duration(rand.Float64()*5.0+ mintime) * time.Second)
		if rand.Intn(troublerate) == 0 {
			fmt.Println("Elfo: Tengo problemas con este juguete")
			factory.Lock() //paro para comprobar
			nelvesintrouble++
			if nelvesintrouble == minelvesintrouble{
				fmt.Println("Elfo: A despertar al jefe para que nos ayude")
				santasBed.Lock() //Despierto
				if note == 0{
					elvesChan <- 0 //despierto a santa
					wait.Lock() //esperamos hasta que termine de ayudar
					nelvesintrouble = 0 //reseteo el conteo
					ntoysmade += 3 //actualizo el número de juguetes
					fmt.Println("Elfo: Ahí van " + strconv.Itoa(minelvesintrouble) + 
					" juguetes más. " + strconv.Itoa(ntoysmade) + " en total.")
				}else{
					fmt.Println("Elfo: Aquí pone que se ha ido...")
				}
				santasBed.Unlock() //Despierto
			}
			factory.Unlock()
		}else{
			fmt.Println("Juguete "+ strconv.Itoa(ntoysmade) +" terminado")
			ntoysmade++
		}
	}
	wg.Done()
}

func reindeer(reindeersChan chan int){
	sledge.Lock() //bloqueo para comprobar
	//compruebo cuantos renos han llegado
	nreindeers++
	if nreindeers == R{
		santasBed.Lock()
		fmt.Println("Rudolf: Soy el reno que faltaba, perdón por llegar tarde...¡A repartir!")
		reindeersChan <- 1 //Notificamos a Santa
		wait.Lock() //Esperamos a que Santa coja el trineo
		santasBed.Unlock()
	}else{
		//nreindeers++
		fmt.Println("un reno llega, van "+strconv.Itoa(nreindeers))
	}
	sledge.Unlock()
	//wg.Done()
}

func launchElves(toyStore *[Rows][Spaces]bool,
	caches *[NBattalions][Rows][Spaces]bool,
	elvesChan chan int,
	leadersChan chan int) {
	
	var toyChans [NBattalions] chan int //los líderes dan juguetes a los elfos
	for i := 0; i < NBattalions; i++ { //un lider asociado a su canal
		toyChans[i] = make(chan int)
		go leader(i, caches, toyStore, toyChans[i], leadersChan)
		for j := 0; j < D; j++{ //con sus elfos correspondientes
			go elf(toyChans[i], elvesChan)
		}	
	}
}

func launchReindeers(reindeersChan chan int) {
	mintime := 5.0 //tiempo minimo entre llegada de renos
	
	for i := 0; i < R; i++ {
		time.Sleep(time.Duration(rand.Float64()*2.0 + mintime) * time.Second)
		if note == 2{
			break
		}
		go reindeer(reindeersChan)
	}
	wg.Done()
}

func initStore(store *[Rows][Spaces]bool, isToyStore bool){
	for i := 0; i < Rows; i++ {
		for j := 0; j < Spaces; j++{
			if isToyStore {
				(*store)[i][j] = true //ya hay regalo
			}else{
				(*store)[i][j] = false //no hay regalo aun
			}
		}
	}
}

//'0' si está vacío, 'X' si está lleno
func printStatus(input *[Rows][Spaces]bool) {
	for i := 0; i < Rows; i++ {
		for j:= 0; j < Spaces; j++ {
			if input[i][j] {
				fmt.Print("X")
			}else{
				fmt.Print("0")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	
	var toyStore [Rows][Spaces]bool //Almacén de juguetes
	var cache [NBattalions][Rows][Spaces]bool //Una caché por batallon
	reindeersChan := make(chan int)
	elvesChan := make(chan int)
	leadersChan := make(chan int)

	//Al principio el almacen principal esta lleno
	//Las caches estan vacias
	//Segundo parametro isToyStore
	//Si isToyStore=true es el almacen

	initStore(&toyStore, true)
	for i := 0; i < NBattalions; i++ {
		initStore(&cache[i], false)
	} 
	fmt.Println("--ESTADO DEL ALMACÉN--")
	printStatus(&toyStore)
	 
	//wg.Add(NBattalions*(D+1)+R+1) //gorutinas a esperar
	wg.Add(NBattalions*(D+1) + 2) //goroutines to wait for
	rand.Seed(time.Now().UTC().UnixNano()) //aleatoreidad, necesario, si no no serían 
											//los mismos números siempre en las funciones rand
	wait.Lock()
	go santa(reindeersChan, elvesChan, leadersChan)
	go launchReindeers(reindeersChan)
	go launchElves(&toyStore, &cache, elvesChan, leadersChan)
	
	wg.Wait() //Espero que acaben todas las gorutinas lanzadas
	fmt.Println("Fin. ¡Feliz Navidad!")

	fmt.Println("--ESTADO DEL ALMACÉN Y LAS CACHÉS--")
	printStatus(&toyStore)
	for i := 0; i<NBattalions; i++{
		printStatus(&cache[i])
	}
}

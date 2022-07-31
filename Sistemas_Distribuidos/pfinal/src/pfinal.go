//Víctor Rincón Yepes

/*
Nota:He añadido un segundo de espera tras el envío de cada regalo esto provoca que en prácticamente
la mayoría de las ejecuciones el programa finalice por la llegada de los renos,
si se quiere ver el comportamiento en el caso de que la condición de finalización del programa
sea que los regalos del almacén se han acabado basta con comentar la línea 122: time.Sleep(time.Duration(1) * time.Second)
*/

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
	var alert = 0 //Para controlar las acciones conjuntas de renos santa y elfos
	var storeMu sync.Mutex //Mutex para protejer el almacén y las cachés

func santa(reindeersChan chan int, elvesChan chan int, leadersChan chan int){
	//Santa Claus invierte un tiempo aleatorio de entre 2 y 5 
	//segundos en ayudar a un equipo de 3 elfos con problemas
	mintime := 2.0
	maxtime := 5.0
	
	for{
		fmt.Println("Santa duerme...")
		//Recibe una llamada por uno de los canales
		select {

		case <- reindeersChan:
			fmt.Println("Santa: Ya están todos los renos aquí, ¡creía que no llegarían nunca!")
			fmt.Println("Santa: Aviso a los elfos de que me voy (alert=1)")
			alert = 1
			wait.Unlock()
			wg.Done()
			
		case <- elvesChan:
			//Ayudar a los elfos
			fmt.Println("Santa: Voy a ayudar a estos duendecillos")
			time.Sleep(time.Duration(rand.Float64()*(maxtime-mintime)+mintime)* time.Second) //aleatorio entre 2 y 5 seg
			wait.Unlock()
		case <- leadersChan:
			fmt.Println("Santa: Ya habéis terminado los regalos, genial, este año tendréis cesta de Navidad.")
			wg.Done()
		}
	}
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
	toyStore *[Rows][Spaces]bool, elvesChan chan int, leadersChan chan int) {

	full := 0 //número de filas llenas, cuando full sea igual a rows paramos.
	var toys [Rows]bool //registro para las filas ya elegidas

	//elegimos un juguete, actualizamos las cachés y se lo pasamos a los elfos usando un canal.
	Loop:
	for {
		initToys(&toys)
		for i := 0; i < Rows; {
			if alert == 1 { //Se para de trabajar y se indica a los elfos que acaben
				fmt.Println("LíderB" + strconv.Itoa(tag) +
				": Santa se ha ido, terminad lo que estéis haciendo y a casa")
				break Loop
			}
			row := rand.Intn(Rows) //Una fila al azar
			storeMu.Lock()
			if !toys[row] { //Si esa fila no se ha elegido aun
				toys[row] = true //La marcamos como elegida
				n := rand.Intn(Spaces) //Un regalo al azar
				for j := 0; j < Spaces; j++{ //Hasta que un regalo que aún esté disponible se elija 
					pos := (n + j) % Spaces //Paso al siguiente cuando esté vacío
					if !(*caches)[tag][row][pos] { //Compruebo en mi caché si el regalo está disponible
						toyStore[row][pos] = false //Lo retiro del almacén
						//Actualizo ambas cachés
						for k := 0; k < NBattalions; k++ {
							(*caches)[k][row][pos] = true
						}
						fmt.Println("LíderB" + strconv.Itoa(tag) + ": ¡Va un regalo!")
						elvesChan <- 1 //Le mando un regalo a los elfos
						break
					}else if j == (Spaces-1){
						full++ //La fila está llena
					}
				}
				i++
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
		elvesChan <- 0 //indicamos a los elfos que hemos terminado
	}

	santasBed.Lock()
	if alert == 0 {
		alert = 2 //si alert == 2 launchReindeers acaba
		leadersChan <- 0 //señal a santa de que hemos acabado
	}
	santasBed.Unlock()
	wg.Done()
	
}

func elf(toysChan chan int, elvesChan chan int){
	mintime := 3.0 //tiempo para fabricar un juguete
	troublerate := 3 //uno de cada 3 juguetes tienen problemas
	minelvesintrouble := 3 //numero de elfos para despertar a santa

	for {
		if (<-toysChan) == 0{ //el lider dice que no quedan juguetes, terminamos
			break
		}
		fmt.Println("Elfo: ¡Ya tengo trabajo!, recibí un juguete.")
		time.Sleep(time.Duration(rand.Float64()*5.0+ mintime) * time.Second)
		if rand.Intn(troublerate) == 0 { //Problemas
			fmt.Println("Elfo: Tengo problemas con este juguete")
			factory.Lock() //paro para comprobar
			if nelvesintrouble == minelvesintrouble - 1{
				fmt.Println("Elfo: A despertar al jefe para que nos ayude")
				santasBed.Lock()
				if alert == 0{
					elvesChan <- 0 //despierto a santa
					wait.Lock() //esperamos hasta que termine de ayudar
					nelvesintrouble = 0 //reseteo el conteo
					ntoysmade += 3 //actualizo el número de juguetes
					fmt.Println("Elfo: Ahí van " + strconv.Itoa(minelvesintrouble) + 
					" juguetes más.")
				}else{
					fmt.Println("Elfo: Santa ya se se ha ido...")
				}
				santasBed.Unlock()
			}else{
				nelvesintrouble++
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
		fmt.Println("un reno llega, van "+strconv.Itoa(nreindeers))
	}
	sledge.Unlock()
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
		if alert == 2{
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
				(*store)[i][j] = true //ya hay juguete
			}else{
				(*store)[i][j] = false //no hay juguete aun
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

	initStore(&toyStore, true)
	for i := 0; i < NBattalions; i++ {
		initStore(&cache[i], false)
	} 
	fmt.Println("--ESTADO DEL ALMACÉN--")
	printStatus(&toyStore)
	
	wg.Add(NBattalions*(D+1) + 2) //gorutinas
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
	for i := 0; i < NBattalions; i++{
		printStatus(&cache[i])
	}
}

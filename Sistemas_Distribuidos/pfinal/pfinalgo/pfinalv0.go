//Con respecto al almacén, los valores booleanos de las cachés están
//invertidos ya que en el almacén representan si hay un regalo en esa
//posición, mientras que en las cachés representan si está disponible
//esa posición (no hay ni ha habido regalo ahí)

//Cuando un elfo tiene problemas con un regalo, lo deja hasta que Santa
//lo arregle y mientras tanto continua con otros regalos. Es decir, la
//gorrutina no se queda esperando a Santa (a excepción del tercer elfo
//el cual necesariamente espera a Santa para acto seguido actualizar la
//cuenta de elfos en problemas).

//Santa deja una nota en su puerta cuando se va a entregar regalos para
//que el elfo que va a despertarlo no se quede esperando para siempre.
package main

import (
	"fmt"
	"strconv"
    "time"
	"sync"
	"math/rand"
)

const (
	Shelves = 3 //shelves in the stash
	Positions = 30 //presents per shelf
	D = 15 //number of elves not including leader
	B = 2 //number of battalions
	R = 9 //number of reindeers
)
var (
	elfWait = &sync.Mutex{} //wait while santa helps you (dont unlock factory yet)
	santasBed = &sync.Mutex{} //where Santa sleeps and is awakened
	sledge = &sync.Mutex{} //protect arrived
	factory = &sync.Mutex{} //protect introuble
	arrived = 0 //reindeers that have arrived so far. shared variable.
	introuble = 0 //elves that are in trouble and need santa's help. shared.
	totalMade = 0 //presents made so far
	wg sync.WaitGroup
	stashKey = &sync.Mutex{}
)

func main() {
	var stash [Shelves][Positions]bool //true if available (ALMACËN)

	initPresents(&stash, true) //isStash = true
	wg.Add(B*(D+1)+R+1) //total number of goroutines
	rand.Seed(time.Now().UTC().UnixNano()) //random seeding
	santasBed.Lock()
	elfWait.Lock()

	go santa()
	go sendoffRdeers()
	go sendoffElves(&stash)
	wg.Wait()
	fmt.Println("Y colorín colorado, el main() ha terminado.")
}

func santa() {
	min := 2.0 //minimum time for santa to help out a group of elves
	max := 5.0 //same but max

	for{
		fmt.Println("Santa: A dormir ¡Ho-ho-ho!")
		//block at lock (go to sleep)
		santasBed.Lock()

		fmt.Println("Santa: Ya voy, no grites...")
		sledge.Lock()
		//check who woke santa up. reindeers have priority
		if arrived == R -1 { //it was the last reindeer
			//fmt.Println("Santa: Antes de irme dejaré una nota para los elfos...")
			fmt.Println("Santa: Han llegado todos los renos, ¡a repartir!")
			break
		}else { //it was an elf
			//help elves (sleep)
			fmt.Println("Santa: *ayuda a los elfos*")
			time.Sleep(time.Duration(rand.Float64()*(max-min)+min)* time.Second)
			elfWait.Unlock()
		}
		sledge.Unlock()
	}
	fmt.Println("Santa: *termina de repartir y muere (la gorutina, no Santa)*")
	wg.Done()
}

func leader(tag int, stash *[Shelves][Positions]bool, give chan int) {
	coffeeBreak := 1 //coffee break time (seconds)
	full := 0 //number of full shelfs. when full==shelves: stop
	var cache [Shelves][Positions]bool //true if there is/was a present (CACHÉ)
	var reg [Shelves]bool //registry for shelves already chosen

	initPresents(&cache, false) //isStash = false
	//choose present, update caches, send present to elves.
	for ; full < Shelves; {
		initReg(&reg)
		for i := 0; i < Shelves; {
			row := rand.Intn(Shelves) //choose a random shelf
			stashKey.Lock()
			if !reg[row] { //if that shelf had not been chosen yet
				reg[row] = true //mark it as chosen
				n := rand.Intn(Positions) //choose a random present
				for j := 0; j < Positions; j++{ //until an available present is randomly selected
					pos := (n + j) % Positions //when miss, choose next present
					if stash[row][pos] { //if the present is available
						stash[row][pos] = false //take it from stash
						cache[row][pos] = true //place it in cache
						//actualizar la otra caché
						fmt.Println("Líd" + strconv.Itoa(tag) + ": ¡Va un regalo!")
						give <- 1 //give present to elves
						break
					}else if j == (Positions-1){
						full++ //shelf is full
					}
				}
				i++ //misses do not count
			}
			stashKey.Unlock()
			//coffee break after working hard (or hardly working...)
			time.Sleep(time.Duration(coffeeBreak) * time.Second)
		}
	}
	for i := 0; i < D; i++ {
		give <- 0 //tell elves that we have finished
	}
	wg.Done()
}

func initPresents(x *[Shelves][Positions]bool, isStash bool) {
	for i := 0; i < Shelves; i++ {
		for j := 0; j < Positions; j++ {
			if isStash {
				(*x)[i][j] = true //present is available
			}else {
				(*x)[i][j] = false //there is no present yet
			}
		}
	}
}

func initReg(reg *[Shelves]bool) {
	for i := 0; i < Shelves; i++ {
		(*reg)[i] = false
	}
}

func elf(receive chan int) {
	min := 2.0 // minimum time to finish working
	v := 10.0 //time span factor
	P := 3 //1 in P elves will be in trouble (randomly)
	G := 3 //number of elves in trouble in order to wake santa

	for {
		if (<- receive) == 0 { //leader said we have finished already
			break
		}
		fmt.Println("Elfo: Regalo recibido, ¡a fabricarlo!")
		time.Sleep(time.Duration(rand.Float64()*v + min) * time.Second)

		if rand.Intn(P) == 0 { //1 in every P will need help
			//elf in trouble!
			factory.Lock()
			fmt.Println("Elfo: ¡Ayuda!")
			if introuble == G-1 {
				//wake santa
				fmt.Println("Elfo: Ya somos " + strconv.Itoa(G) + " elfos con problemas")
				fmt.Println("Elfo: ¡Santa! ¡Despierta!")

				santasBed.Unlock()
				elfWait.Lock() //wait for santa to finish helping us
				//reset introuble count
				introuble = 0
				//update total presents made
				totalMade += 3
				fmt.Println("Elfo: Ahí van " + strconv.Itoa(G) + " juguetes más. " + strconv.Itoa(totalMade) + " en total.")

			}else {
				introuble++
			}
			factory.Unlock()
		}else {
			totalMade++
			fmt.Println("Elfo: Juguete " + strconv.Itoa(totalMade) + " terminado")
		}
	}
	wg.Done()
}

func reindeer() {
	sledge.Lock()
	//check how many reindeers have arrived so far
	if arrived == R -1 {
		arrived++
		//im last. must wake santa.
		fmt.Println("Reno: Ya estoy. Conmigo somos " + strconv.Itoa(R) + " y estamos todos")
		fmt.Println("Reno: ¡Santa! ¡Hora de irse, que no nos da tiempo!")
		santasBed.Unlock() //wake
	}else {
		arrived++
		fmt.Println("Reno: Ya he llegado. Ya somos " + strconv.Itoa(arrived) + " renos")
	}
	sledge.Unlock()
	wg.Done()
}

func sendoffRdeers() {
	min := 12.0 //minimum time between reindeer's arrivals
	v := 2.0 //time span factor

	for i := 0; i < R; i++ {
		time.Sleep(time.Duration(rand.Float64()*v + min) * time.Second)
		go reindeer()
	}
}

func sendoffElves(stash *[Shelves][Positions]bool) {
	var channels [B]chan int //where leaders give presents to elves

	for i := 0; i < B; i++ {
		channels[i] = make(chan int) //where leader gives presents to elves
		go leader(i, stash, channels[i])
		for j := 0; j < D; j++ {
			go elf(channels[i])
		}
	}
}

func printCache(cache *[Shelves][Positions]bool) {
	for i := 0; i < Shelves; i++ {
		for j:= 0; j < Positions; j++ {
			if cache[i][j] {
				fmt.Print("O")
			}else{
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

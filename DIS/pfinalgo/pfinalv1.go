
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
	wait = &sync.Mutex{} //wait for Santa
	santasRoom = &sync.Mutex{}
	sledge = &sync.Mutex{} //protects arrived
	factory = &sync.Mutex{} //protects introuble
	door = &sync.Mutex{} //where the note is posted
	stashKey = &sync.Mutex{} //protects stash and caches
	arrived = 0 //reindeers that have arrived so far. shared variable.
	introuble = 0 //elves that are in trouble and need santa's help. shared.
	totalMade = 0 //presents made so far
	note = 0 //written by Santa upon leaving
	wg sync.WaitGroup
)

func main() {
	var stash [Shelves][Positions]bool //true if available (ALMACËN)
	var caches [B][Shelves][Positions]bool //true if no present placed (CACHÉS)
	reinWake := make(chan int) //channel for the reindeers to wake Santa up
	elfWake := make(chan int) //channel for the elves to wake Santa up
	end := make(chan int) //channel for the leaders to notify upon ending

	initPresents(&stash, true) //isStash = true
	for i := 0; i < B; i++ {
		initPresents(&caches[i], false) //isStash = false
	}
	wg.Add(B*(D+1)+2) //goroutines to wait for
	rand.Seed(time.Now().UTC().UnixNano()) //random seeding
	wait.Lock()

	go santa(reinWake, elfWake, end)
	go sendoffRdeers(reinWake)
	go sendoffElves(&stash, &caches, elfWake, end)
	wg.Wait()
	fmt.Println("Y colorín colorado, el main() ha terminado.")
}

func santa(reindeers, help, end chan int) {
	min := 2.0 //minimum time for santa to help out a group of elves
	max := 5.0 //same but max

	Loop:
	for{
		fmt.Println("Santa: A dormir ¡Ho-ho-ho!")
		//receive a wake-up call from someone
		select {
		case <- reindeers:
			fmt.Println("Santa: Ya voy, no grites...")
			fmt.Println("Santa: Antes de irme dejaré una nota para los elfos...")
			note = 1 //protected by santasRoom in reindeer
			fmt.Println("Santa: Han llegado todos los renos, ¡a repartir!")
			wait.Unlock()
			break Loop
		case <- help:
			fmt.Println("Santa: Ya voy, no grites...")
			//help elves (sleep)
			fmt.Println("Santa: *ayuda a los elfos*")
			time.Sleep(time.Duration(rand.Float64()*(max-min)+min)* time.Second)
			wait.Unlock()
		case <- end:
			fmt.Println("Santa: Ya habéis terminado los regalos, genial.")
			break Loop
		}
	}
	wg.Done()
}

func leader(tag int, caches *[B][Shelves][Positions]bool,
	stash *[Shelves][Positions]bool, elves, end chan int) {

	coffeeBreak := 1 //coffee break time (seconds)
	full := 0 //number of full shelfs. when full==shelves: stop
	var reg [Shelves]bool //registry for shelves already chosen

	//choose present, update caches, send present to elves.
	Loop:
	for ; full < Shelves; {
		initReg(&reg)
		for i := 0; i < Shelves; {
			if note == 1 { //stop working and signal elves to finish
				fmt.Println("Líd" + strconv.Itoa(tag) +
				": Santa se ha ido, terminad lo que estéis haciendo y a casa")
				break Loop
			}
			row := rand.Intn(Shelves) //choose a random shelf
			stashKey.Lock()
			if !reg[row] { //if that shelf has not been chosen yet
				reg[row] = true //mark it as chosen
				n := rand.Intn(Positions) //choose a random present
				for j := 0; j < Positions; j++{ //until an available present is randomly selected
					pos := (n + j) % Positions //when miss, choose next present
					if !(*caches)[tag][row][pos] { //if the present is available
						stash[row][pos] = false //take it from stash
						//update all caches
						for k := 0; k < B; k++ {
							(*caches)[k][row][pos] = true
						}
						fmt.Println("Líd" + strconv.Itoa(tag) + ": ¡Va un regalo!")
						elves <- 1 //give present to elves
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
		elves <- 0 //tell elves that we have finished
	}
	santasRoom.Lock()
	if note == 0 {
		note = 2 //let reindeers know that we've finished
		end <- 0 //signal Santa that we've finished
	}
	santasRoom.Unlock()
	wg.Done()
}

func elf(leader, wakeSanta chan int) {
	min := 1.0 // minimum time to finish working
	v := 5.0 //time span factor
	P := 3 //1 in P elves will be in trouble (randomly)
	G := 3 //number of elves in trouble in order to wake santa

	for {
		if (<- leader) == 0 { //leader said we have finished already
			break
		}
		fmt.Println("Elfo: Regalo recibido, ¡a fabricarlo!")
		time.Sleep(time.Duration(rand.Float64()*v + min) * time.Second)
		if rand.Intn(P) == 0 { //1 in every P will need help
			//elf in trouble!
			factory.Lock()
			fmt.Println("Elfo: ¡Ayuda!")
			if introuble == (G-1) {
				//wake santa
				santasRoom.Lock()
				fmt.Println("Elfo: Ya somos " + strconv.Itoa(G) +
				" elfos con problemas")
				fmt.Println("Elfo: ¡Santa! ¡Despierta!")
				if note == 0 { //Santa is here
					wakeSanta <- 0 //wake Santa
					wait.Lock() //wait for santa to finish helping us
					introuble = 0 //reset introuble count
					totalMade += 3 //update total presents made
					fmt.Println("Elfo: Ahí van " + strconv.Itoa(G) +
					" juguetes más. " + strconv.Itoa(totalMade) + " en total.")
				}else {
					fmt.Println("Elfo: Vaya, parece que Santa se ha ido...")
				}
				santasRoom.Unlock()
			}else {
				introuble++
			}
			factory.Unlock()
		}else { //present made without difficulties
			totalMade++
			fmt.Println("Elfo: Juguete " + strconv.Itoa(totalMade) + " terminado")
		}
	}
	wg.Done()
}

func reindeer(wakeSanta chan int) {
	sledge.Lock()
	//check how many reindeers have arrived so far
	if arrived == R -1 {
		santasRoom.Lock()
		arrived++
		sledge.Unlock()
		//im last. must wake santa.
		fmt.Println("Reno: Ya estoy. Conmigo somos " + strconv.Itoa(R) +
		" y estamos todos")
		fmt.Println("Reno: ¡Santa! ¡Hora de irse, que no nos da tiempo!")
		wakeSanta <- 1 //wake Santa
		wait.Lock() //wait for Santa to get on the sledge
		santasRoom.Unlock()
	}else {
		arrived++
		sledge.Unlock()
		fmt.Println("Reno: Ya he llegado. Ya somos " + strconv.Itoa(arrived) + " renos")
	}
}

func sendoffElves(stash *[Shelves][Positions]bool,
	caches *[B][Shelves][Positions]bool, wakeSanta, end chan int) {
	var channels [B]chan int //where leaders give presents to elves

	for i := 0; i < B; i++ {
		channels[i] = make(chan int) //where leader gives presents to elves
		go leader(i, caches, stash, channels[i], end)
		for j := 0; j < D; j++ {
			go elf(channels[i], wakeSanta)
		}
	}
}

func sendoffRdeers(wakeSanta chan int) {
	min := 5.0 //minimum time between reindeer's arrivals
	v := 10.0 //time span factor

	for i := 0; i < R; i++ {
		time.Sleep(time.Duration(rand.Float64()*v + min) * time.Second)
		if note == 2 {
			break
		}
		go reindeer(wakeSanta)
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

//for debugging purposes
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

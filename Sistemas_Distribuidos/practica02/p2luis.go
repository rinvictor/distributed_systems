package main

import (
	"fmt"
	"strconv"
    "time"
	"sync"
	"math/rand"
)

var (
	elfWait = &sync.Mutex{} //wait while santa helps you (dont unlock factory yet)
	santasBed = &sync.Mutex{} //where Santa sleeps and is awakened
	sledge = &sync.Mutex{} //protect arrived
	factory = &sync.Mutex{} //protect introuble
	arrived = 0 //reindeers that have arrived so far. shared variable.
	introuble = 0 //elves that are in trouble and need santa's help. shared.
	wg sync.WaitGroup
)

func main() {
	D := 12 //number of elves
	R := 9 //number of reindeers
	wg.Add(D+R+1) //total number of goroutines
	rand.Seed(time.Now().UTC().UnixNano()) //random seeding
	santasBed.Lock()
	elfWait.Lock()

	go santa(R)
	go sendoffRdeers(R)
	go sendoffElves(D)
	wg.Wait()
	fmt.Println("Y colorín colorado, el main() ha terminado.")
}

func santa(R int) {
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

func elf() {
	min := 2.0 // minimum time to finish working
	v := 40.0 //time span factor
	P := 1 //1 in P elves will be in trouble (randomly)
	G := 1 //number of elves in trouble in order to wake santa
	//work (sleep. oh, the irony...)
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
		}else {
			introuble++
		}
		factory.Unlock()
	}else {
		fmt.Println("Elfo: Juguete terminado")
	}
	wg.Done()
}

func reindeer(R int) {
	sledge.Lock()
	//check how many reindeers have arrived so far
	if arrived == R -1 {
		//im last. must wake santa.
		fmt.Println("Reno: Ya estoy. Conmigo somos " + strconv.Itoa(R) + " y estamos todos")
		fmt.Println("Reno: ¡Santa! ¡Hora de irse, que no nos da tiempo!")
		santasBed.Unlock()
	}else {
		arrived++
		fmt.Println("Reno: Ya he llegado. Ya somos " + strconv.Itoa(arrived) + " renos")
	}
	sledge.Unlock()
	wg.Done()
}

func sendoffRdeers(R int) {
	min := 5.0 //minimum time between reindeer's arrivals
	v := 2.0 //time span factor
	for i := 0; i < R; i++ {
		time.Sleep(time.Duration(rand.Float64()*v + min) * time.Second)
		go reindeer(R)
	}
}

func sendoffElves(D int) {
	for i := 0; i < D; i++ {
		go elf()
	}
}

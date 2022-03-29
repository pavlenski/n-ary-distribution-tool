package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var arity = 2

func main() {
	// the key of the map for the bag of words could be a string
	// the checking of its existence can either be done with contains (checking if the words of the key are the same)
	// or the words being

	// log2arity - sorting the words to check the key (using O(log2n) sort)
	// arity^2 - contains

	times, _ := time.ParseDuration("5000ms")
	time.Sleep(times)
	fmt.Printf("wtf")

	//
	//fmt.Println("beginning state:", State)
	//sleepChan := make(chan struct{})
	//go sleep(sleepChan)
	////
	////time.Sleep(time.Second)
	////sleepChan <- struct{}{}
	////State = 420
	//time.Sleep(5 * time.Second)
	//
	//fmt.Printf("state: %d\n", State)
}

var State = 1
var sleepDur = 3 * time.Second

func sleep(sleepChan <-chan struct{}) {
	select {
	case <-sleepChan:
		return
	case <-time.Tick(sleepDur):
		fmt.Printf("nigga\n")
		State = 10
	}
}

func cliTest() {
	input1 := make(chan int)
	go input(1, input1)

	buffer := bufio.NewReader(os.Stdin)

	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
			fmt.Println("error scanning command")
			return
		}
		lineFmt := line[:len(line)-1]

		if lineFmt == "pause" {
			fmt.Println("pausing")
			input1 <- 2
		}

		if lineFmt == "resume" {
			fmt.Println("resuming")
			input1 <- 1
		}
	}
}

func input(inputID int, ws <-chan int) {
	fmt.Printf("input [%d] created\n", inputID)
	state := 2

	for {
		select {
		case state = <-ws:
			switch state {
			case 1:
				fmt.Printf("input [%d] resuming\n", inputID)
			case 2:
				fmt.Printf("input [%d] pausing\n", inputID)
			}
		default:
			//runtime.Gosched()
			if state == 2 {
				break
			}
			fmt.Printf("input [%d] working..\n", inputID)
			time.Sleep(time.Second)
		}
	}
}

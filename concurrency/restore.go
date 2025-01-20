package main

import (
	"fmt"
	"math/rand"
	"time"
)

func boring(msg string) <-chan string{ // returns receive only channel of string
	c := make(chan string)
	waitForIt := make(chan bool)
	go func() {
		for i := 0; ; i++ {
			c <- Message{fmt.Sprintf("%s %d", msg, i), waitForIt}
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			<- waitForIt
		}
	}()
	return c;
}

func fanIn(inpu1, input2 <- chan string) <- chan string {
	c := make(chan string)
	go func() { for { c <- <- inpu1 } }()
	go func() { for { c <- <- input2 } }()
	// just shoving things into one channel
	return c
}

type Message struct {
	str string
	wait chan bool
}

func main() {
	c:= fanIn(boring("Joe"), boring("Ann"))
	for i := 0; i < 10; i++ {
		msg1 := <-c; fmt.Println(msg1.str)
		msg2 := <-c; fmt.Println(msg2.str)
		msg1.wait <- true
		msg2.wait <- true
	}
	fmt.Println("You're boring; I'm leaving.")
}
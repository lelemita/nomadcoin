package main

import (
	"fmt"
	"time"
)

func countToTen(c chan<- int) {
	for i := range [5]int{} {
		time.Sleep(500 * time.Millisecond)
		c <- i
	}
	close(c)
}

func receive(label string, c <-chan int) {
	for {
		num, isOpen := <-c
		if !isOpen {
			fmt.Println("We are done.")
			break
		}
		fmt.Println(label, "receive: ", num)
	}
}

func main() {
	c := make(chan int)
	go countToTen(c)
	go receive("A", c)
	d := make(chan int)
	go countToTen(d)
	receive("BBB", d)
}

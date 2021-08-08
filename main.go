package main

import (
	"fmt"
	"time"
)

func countToTen(c chan<- int) {
	for i := range [5]int{} {
		c <- i
		fmt.Printf(">> sent : %d <<\n", i)
	}
	close(c)
}

func receive(label string, c <-chan int) {
	for {
		time.Sleep(500 * time.Millisecond)
		num, isOpen := <-c
		if !isOpen {
			fmt.Println("We are done.")
			break
		}
		fmt.Printf("%s || receive %d ||\n", label, num)
	}
}

func main() {
	c := make(chan int, 3)
	go countToTen(c)
	receive("A", c)
}

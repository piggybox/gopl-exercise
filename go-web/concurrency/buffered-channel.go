package main

import "fmt"

func main() {
	c := make(chan int, 2) // if using 1, it will fail
	c <- 1
	c <- 2
	fmt.Println(<-c)
	fmt.Println(<-c)
}

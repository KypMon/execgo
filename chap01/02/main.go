package main

import (
	"math/rand"
	"time"
)

func main () {
	c := make(chan int, 10)

	go func() {
		for {
			time.Sleep(time.Second)
			c <- rand.Int()
		}
	}()

	for v := range c {
		time.Sleep(time.Second)
		println(v)
	}
}
package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(ch chan int, wg *sync.WaitGroup) {

	go func() {
		i := 1
		for ;; i++ {
			ch <- i
			fmt.Println("p_in", i)
			time.Sleep(time.Second)

			if i == 20 {
				close(ch)
				fmt.Println("close")
				break
			}
			//i++
		}

		fmt.Println("producer done")
		wg.Done();
	}()

}

func consumer(ch chan int, wg *sync.WaitGroup) {
	go func() {
		for {
			if val, ok := <- ch; ok {
				fmt.Println("consumer", val)

				if val != 20 {
					time.Sleep(time.Second)
				} else {
					break
				}
			}
		}
		fmt.Println("consumer done")
		wg.Done()
	}()
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan int, 10)

	producer(ch, &wg)
	wg.Add(1)
	consumer(ch, &wg)
	wg.Wait()
}

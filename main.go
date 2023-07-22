package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

type CirclBuffer struct {
	array []int
	start int
	size  int
	m     sync.Mutex
}

func NewCircl(size int) *CirclBuffer {
	return &CirclBuffer{make([]int, size), -1, size, sync.Mutex{}}
}

func (c *CirclBuffer) Push(el int) {
	c.m.Lock()
	defer c.m.Unlock()
	if c.start == c.size-1 {
		for i := 1; i <= c.size-1; i++ {
			c.array[i-1] = c.array[i]
		}
		c.array[c.start] = el
	} else {
		c.start++
		c.array[c.start] = el
	}
}

func (c *CirclBuffer) Get() []int {
	if c.start <= 0 {
		return nil
	}
	c.m.Lock()
	var output []int = c.array[:c.start]
	c.start = 0
	return output
}

func read(input chan<- int) {
	for {
		var u int
		_, err := fmt.Scanf("%d\n", &u)
		if err != nil {
			fmt.Println("Это не цифра")
		} else {
			input <-u
		}
	}
}

func removeNegatives(currentChan <-chan int, nextChan chan<- int) {

	for number := range currentChan {
		if number >= 0 {
			nextChan <- number
		}
	}
}

func removedivTree(currentChan <-chan int, nextChan chan<- int) {

	for number := range currentChan {
		if number%3 != 0 {
			nextChan <- number
		}
	}
}

func writetoBuffer(currentChan <-chan int, c *CirclBuffer){
	for number := range currentChan{
	c.Push(number)
	}
}

func writeToConsole(c *CirclBuffer, t *time.Ticker){
for range t.C{
	buffer := c.Get()
	if len(buffer) > 0{
		fmt.Println("Буфер:", buffer)
	}
}
}

func main() {
	

	input := make(chan int)
	go read(input)
	negFilter := make(chan int)
	go removeNegatives(input, negFilter)
	divTreeCh := make(chan int)
	go removedivTree(negFilter, divTreeCh)
	size := 20
c := NewCircl(size)
	go writetoBuffer(divTreeCh, c)

	delay := 5
	ticker := time.NewTicker(time.Second * time.Duration(delay))
	go writeToConsole(c, ticker)

	y := make(chan os.Signal)
	signal.Notify(y, os.Interrupt)
	select {
	case sig := <-y:
		fmt.Printf("Сигнал выхода %s ... \n", sig)
		os.Exit(0)

	}
}

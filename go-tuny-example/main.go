package main

import (
	"log"
	"runtime"
	"time"

	"github.com/Jeffail/tunny"
)

func poller(ch chan<-int) {
	i := 1
	
	for {
		ch <- i
		i++
		time.Sleep(time.Second*10)
	}
}

func dispatcher(in <-chan int, out chan<- int, pool *tunny.Pool) {
	resultQueue := make(chan []int)
	go resultProc(resultQueue)
	for msg := range in {
		store(msg)
		for i := 0; i < 10; i++ {
			go func(item, m int) {
				result := pool.Process([]int{item, m})
				resultQueue <- result.([]int)
			}(msg, i)
		}
	}
}

func resultProc(results <-chan []int) {
	for result := range results {
		item, num := result[0], result[1]
		log.Printf("[ResultProc] %d, %d", item, num)
	}
}

func store(msg int) {
	log.Printf("[Store] %d", msg)
}

func worker(payload interface{}) interface{} {
	_payload := payload.([]int)
	time.Sleep(time.Second)
	log.Printf("[Worker] Item = %v, Payload = %v", _payload[0], _payload[1])
	return payload
}

func main() {
	numCPUs := runtime.NumCPU()
	log.Print(numCPUs)

	pool := tunny.NewFunc(numCPUs, worker)
	pool.SetSize(10)
	defer pool.Close()

	in := make(chan int)
	out := make(chan int)

	go dispatcher(in, out, pool)
	poller(in)
}

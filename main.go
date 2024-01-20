package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	wp "workerpool/workerpool"
)

var wg sync.WaitGroup

func task(id ...any) {
	defer wg.Done()
	rand := rand.Int31n(15)
	time.Sleep(time.Duration(rand) * time.Second)
	fmt.Printf("reach task end for goroutine %v after time %v \n", id, rand)
}

func main() {
	worker := wp.NewWorkerConfig(5)
	worker.Start()
	for i := 0; i < 20; i++ {
		wg.Add(1)
		worker.SubmitJob(task)
		fmt.Printf("loop interval %v \n", i)
	}
	fmt.Println("all task submitted")
	wg.Wait()
	fmt.Println("done")
}

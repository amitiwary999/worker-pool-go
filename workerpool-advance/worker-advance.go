package workerpooladvance

import "fmt"

type worker struct {
	pool  *pool
	tasks chan func()
}

func (w *worker) doJob() {
	w.pool.addRunning(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic in do job %v ", r)
			}
		}()
		for task := range w.tasks {
			fmt.Println("about to start task")
			task()
			w.pool.addRunning(-1)
			w.pool.workerPool.Put(w)
			w.pool.cond.Signal()
		}
	}()
}

func (w *worker) submitJob(fn func()) {
	w.tasks <- fn
}

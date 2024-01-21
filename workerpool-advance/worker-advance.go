package workerpooladvance

import "fmt"

type worker struct {
	pool  *pool
	tasks chan func()
}

func (w *worker) doJob() {
	go func() {
		defer func() {
			w.pool.workerPool.Put(w)
			if r := recover(); r != nil {
				fmt.Printf("panic in do job %v ", r)
			}
			w.pool.cond.Broadcast()
		}()
		for task := range w.tasks {
			task()
		}
	}()
}

func (w *worker) submitJob(fn func()) {
	w.tasks <- fn
}

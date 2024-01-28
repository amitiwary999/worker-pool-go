package workerpooladvance

import (
	"fmt"
	"time"
)

type worker struct {
	pool         *pool
	tasks        chan func()
	lastUsedTime time.Time
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
			if task == nil {
				return
			}
			fmt.Println("about to start task")
			task()
			w.returnWorkerPool()
		}
	}()
}

func (w *worker) returnWorkerPool() {
	w.lastUsedTime = time.Now()
	w.pool.workerPool.Put(w)
	w.pool.addRunning(-1)
	w.pool.cond.Signal()
	w.pool.lock.Unlock()
}

func (w *worker) submitJob(fn func()) {
	w.tasks <- fn
}

func (w *worker) finish() {
	w.tasks <- nil
}

package workerpooladvance

import (
	"fmt"
	"time"
)

type workerIntf interface {
	doJob()
	finish()
	submitJob(func())
	getLastUsedTime() time.Time
}

type worker struct {
	pool         *pool
	tasks        chan func()
	lastUsedTime time.Time
}

func (w *worker) doJob() {
	w.pool.addRunning(1)
	go func() {
		defer func() {
			w.pool.addRunning(-1)
			w.pool.workerPool.Put(w)
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
			w.pool.returnWorkerPool(w)
		}
	}()
}

func (w *worker) submitJob(fn func()) {
	w.tasks <- fn
}

func (w *worker) finish() {
	w.tasks <- nil
}

func (w *worker) getLastUsedTime() time.Time {
	return w.lastUsedTime
}

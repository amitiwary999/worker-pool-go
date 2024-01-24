package workerpooladvance

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type pool struct {
	capacity   int32
	running    int32
	lock       sync.Mutex
	cond       *sync.Cond
	workerPool sync.Pool
}

func NewPool(capacity int32) *pool {
	p := &pool{
		capacity: capacity,
	}
	p.cond = sync.NewCond(&p.lock)
	p.workerPool.New = func() interface{} {
		return &worker{
			pool:  p,
			tasks: make(chan func(), capacity),
		}
	}
	return p
}

func (p *pool) addRunning(taskCount int) {
	atomic.AddInt32(&p.running, int32(taskCount))
}

func (p *pool) Submit(fn func()) {
	w, err := p.getWorker()
	if err != nil {
		fmt.Println("failed to get worker")
	}
	fmt.Println("submit task")
	w.submitJob(fn)
}

func (p *pool) getWorker() (w *worker, err error) {
	p.lock.Lock()
retry:
	if c := p.capacity; c > p.running {
		w = p.workerPool.Get().(*worker)
		p.lock.Unlock()
		w.doJob()
		fmt.Printf("worker count %v capacity %v \n", p.running, c)
		return
	}
	fmt.Println("before the wait")
	p.cond.Wait()
	goto retry
}

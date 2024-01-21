package workerpooladvance

import (
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

func (p *pool) addRunning() {
	atomic.AddInt32(&p.running, 1)
}

func (p *pool) submit(fn func()) {

}

func (p *pool) getWorker() (w worker, err error) {
	p.lock.Lock()
retry:
	if c := p.capacity; c > p.running {
		w = p.workerPool.Get().(worker)
		p.lock.Unlock()
		w.doJob()
		return
	}
	p.cond.Wait()
	goto retry
}

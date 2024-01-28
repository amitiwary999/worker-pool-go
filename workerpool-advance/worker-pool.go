package workerpooladvance

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	expiryTime = 5 * time.Second
)

type pool struct {
	capacity   int32
	running    int32
	lock       sync.Mutex
	cond       *sync.Cond
	workerPool sync.Pool
	workers    []workerIntf
}

func NewPool(capacity int32) *pool {
	p := &pool{
		capacity: capacity,
		workers:  make([]workerIntf, 0, capacity),
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

func (p *pool) getWorker() (w workerIntf, err error) {
	p.lock.Lock()
retry:
	workerLen := len(p.workers)
	if workerLen > 0 {
		w = p.workers[workerLen-1]
		p.workers[workerLen-1] = nil
		p.workers = p.workers[:workerLen-1]
		return
	}

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

func (p *pool) purgeWorker() {
	ticker := time.NewTicker(expiryTime)

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
		}
		unusedWorkers := p.getUnusedWorker()
		for _, unusedWorker := range unusedWorkers {
			unusedWorker.finish()
		}
	}
}

func (p *pool) returnWorkerPool(w *worker) {
	w.lastUsedTime = time.Now()
	p.workers = append(p.workers, w)
	p.cond.Signal()
	p.lock.Unlock()
}

func (p *pool) getUnusedWorker() []workerIntf {
	expiredTime := time.Now().Add(-expiryTime)
	length := len(p.workers)
	last := length
	first := 0
	for first <= last {
		mid := (first + (last-first)<<1)
		if p.workers[mid].getLastUsedTime().Before(expiredTime) {
			first = mid + 1
		} else {
			last = mid - 1
		}
	}
	if last != -1 {
		expired := p.workers[:last]
		copy(p.workers, p.workers[last+1:])
		for i := last + 1; i < length; i++ {
			p.workers[i] = nil
		}
		return expired
	}
	return nil
}

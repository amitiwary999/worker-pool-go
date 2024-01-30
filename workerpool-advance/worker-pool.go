package workerpooladvance

import (
	"context"
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
	cancelFunc context.CancelFunc
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
	ctx, cancelFunc := context.WithCancel(context.Background())
	p.cancelFunc = cancelFunc
	go p.purgeWorker(ctx)
	return p
}

func (p *pool) Release() {
	p.cancelFunc()
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
		p.lock.Unlock()
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

func (p *pool) purgeWorker(ctx context.Context) {
	ticker := time.NewTicker(expiryTime)

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
		p.lock.Lock()
		unusedWorkers := p.getUnusedWorker()
		p.lock.Unlock()
		for _, unusedWorker := range unusedWorkers {
			unusedWorker.finish()
		}
	}
}

func (p *pool) returnWorkerPool(w *worker) {
	w.lastUsedTime = time.Now()
	p.workers = append(p.workers, w)
	p.cond.Signal()
}

func (p *pool) getUnusedWorker() []workerIntf {
	expiredTime := time.Now().Add(-expiryTime)
	length := len(p.workers)
	last := length - 1
	first := 0
	for first <= last {
		mid := (first + (last-first)>>1)
		if p.workers[mid].getLastUsedTime().Before(expiredTime) {
			first = mid + 1
		} else {
			last = mid - 1
		}
	}
	expired := []workerIntf{}
	if last != -1 {
		expired = append(expired, p.workers[:last+1]...)
		copiedCount := copy(p.workers, p.workers[last+1:])
		for i := copiedCount; i < length; i++ {
			p.workers[i] = nil
		}
		p.workers = p.workers[:copiedCount]
		return expired
	}
	return expired
}

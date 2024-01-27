package workerpool

import "fmt"

type workerConfig struct {
	size int
	jobs chan func(params ...any)
}

func NewWorkerConfig(size int) *workerConfig {
	return &workerConfig{
		size: size,
		jobs: make(chan func(params ...any), size),
	}
}

func (w *workerConfig) doJob(id int) {
	for job := range w.jobs {
		job(id)
	}
	fmt.Println("end of for loop")
}

func (w *workerConfig) Start() {
	for i := 0; i < w.size; i++ {
		go w.doJob(i)
	}
}

func (w *workerConfig) SubmitJob(fn func(params ...any)) {
	w.jobs <- fn
}

func (w *workerConfig) Purge() {
	close(w.jobs)
}

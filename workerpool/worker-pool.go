package workerpool

type workerConfig struct {
	size int
	jobs chan func()
}

func NewWorkerConfig(size int) *workerConfig {
	return &workerConfig{
		size: size,
		jobs: make(chan func()),
	}
}

func (w *workerConfig) doJob() {
	for job := range w.jobs {
		job()
	}
}

func (w *workerConfig) Start() {
	for i := 0; i < w.size; i++ {
		go w.doJob()
	}
}

func (w *workerConfig) SubmitJob(fn func()) {
	w.jobs <- fn
}

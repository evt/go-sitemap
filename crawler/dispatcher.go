package crawler

import "log"

// Dispatcher dispatches jobs
type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
	// max workers to run
	MaxWorkers int
	// Job Queue
	JobQueue chan Job
}

// NewWorker creates new worker
func NewWorker(workerPool chan chan Job, seq int) Worker {
	return Worker{
		Seq:        seq,
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

// NewDispatcher creates new dispatcher
func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)

	d := &Dispatcher{
		WorkerPool: pool,
		MaxWorkers: maxWorkers,
		JobQueue:   make(chan Job),
	}
	d.Run()

	return d
}

// Run runs workers and starts handling jobs
func (d *Dispatcher) Run() {
	// starting n number of workers
	log.Printf(`Starting %d workers`, d.MaxWorkers)
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, i)
		worker.Start()
	}

	go d.Dispatch()
}

// Dispatch handles dispatcher jobs
func (d *Dispatcher) Dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			// a job request has been received
			go func(job Job) {
				// Obtain an available worker job channel.
				// This will block until a worker is idle.
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}

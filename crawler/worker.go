package crawler

// Worker represents the worker that executes the job
type Worker struct {
	Seq        int
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

// Start runs loop for the worker, listening for a quit channel in case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				job.Run(w.Seq)

			case <-w.quit:
				// Stop worker
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

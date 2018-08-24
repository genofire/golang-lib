// Package worker to run functions like a cronjob in background
package worker

import (
	"sync"
	"time"
)

// Worker Struct which handles the job
type Worker struct {
	every time.Duration
	run   func()
	quit  chan struct{}
	wg    sync.WaitGroup
}

// NewWorker create a Worker with a timestamp, run, every and it's function
func NewWorker(every time.Duration, f func()) (w *Worker) {
	w = &Worker{
		every: every,
		run:   f,
		quit:  make(chan struct{}),
	}
	return
}

// Start the Worker
func (w *Worker) Start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		ticker := time.NewTicker(w.every)
		for {
			select {
			case <-ticker.C:
				w.run()
			case <-w.quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// Close stops the Worker
func (w *Worker) Close() {
	close(w.quit)
	w.wg.Wait()
}

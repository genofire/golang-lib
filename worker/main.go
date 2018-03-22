// Package with a lib for cronjobs to run in background
package worker

import (
	"sync"
	"time"
)

// Struct which handles the job
type Worker struct {
	every time.Duration
	run   func()
	quit  chan struct{}
	wg    sync.WaitGroup
}

// Function to create a new Worker with a timestamp, run, every and it's function
func NewWorker(every time.Duration, f func()) (w *Worker) {
	w = &Worker{
		every: every,
		run:   f,
		quit:  make(chan struct{}),
	}
	return
}

// Function to start the Worker
// (please us it as a go routine with go w.Start())
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

// Function to stop the Worker
func (w *Worker) Close() {
	close(w.quit)
	w.wg.Wait()
}

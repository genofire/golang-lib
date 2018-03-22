// Package with a lib for cronjobs to run in background
package worker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Function to test the Worker
func TestWorker(t *testing.T) {
	assert := assert.New(t)

	runtime := 0

	w := NewWorker(time.Duration(5)*time.Millisecond, func() {
		runtime++
	})
	w.Start()
	time.Sleep(time.Duration(18) * time.Millisecond)
	w.Close()
	time.Sleep(time.Duration(18) * time.Millisecond)
	assert.Equal(3, runtime)
	assert.Panics(func() {
		w.Close()
	})
}

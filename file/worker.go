package file

import (
	"time"

	"dev.sum7.eu/genofire/golang-lib/worker"
)

// NewSaveJSONWorker Starts a worker, which save periodly data to json file
func NewSaveJSONWorker(repeat time.Duration, path string, data interface{}) *worker.Worker {
	saveWorker := worker.NewWorker(repeat, func() {
		SaveJSON(path, data)
	})
	saveWorker.Start()
	return saveWorker
}

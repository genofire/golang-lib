package file

import (
	"time"

	"dev.sum7.eu/genofire/golang-lib/worker"
)

func NewSaveJSONWorker(repeat time.Duration, path string, data interface{}) *worker.Worker {
	saveWorker := worker.NewWorker(repeat, func() {
		SaveJSON(path, data)
	})
	go saveWorker.Start()
	return saveWorker
}

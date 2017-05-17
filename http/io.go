// Package that provides the logic of the webserver
package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// Function to read data from a http request via json format (input)
func Read(r *http.Request, to interface{}) (err error) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		err = errors.New("no json request recieved")
		return
	}
	err = json.NewDecoder(r.Body).Decode(to)
	return
}

// Function to write data as json to a http response (output)
func Write(w http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

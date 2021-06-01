package web

import (
	"encoding/json"
	"net/http"
	"time"
)

// JSONRequest easy get request for json
func JSONRequest(url string, value interface{}) error {
	var netClient = &http.Client{
		Timeout: time.Second * 20,
	}

	resp, err := netClient.Get(url)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&value)
	if err != nil {
		return err
	}
	return nil
}

package web

import (
	"encoding/json"
	"net/http"
	"time"
)

// JSONRequest issues a GET request to the specified URL and reads the returned
// JSON into value. See json.Unmarshal for the rules for converting JSON into a
// value.
func JSONRequest(url string, value interface{}) error {
	netClient := &http.Client{
		Timeout: time.Second * 20,
	}

	resp, err := netClient.Get(url)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&value)
	resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

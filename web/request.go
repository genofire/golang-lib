package web

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

// NewRequestWithFile Create a Request with file as body
func NewRequestWithFile(url, filename string) (*http.Request, error) {
	buf := bytes.NewBuffer(nil)
	bodyWriter := multipart.NewWriter(buf)

	// We need to truncate the input filename, as the server might be stupid and take the input
	// filename verbatim. Then, he will have directory parts which do not exist on the server.
	fileWriter, err := bodyWriter.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if _, err := io.Copy(fileWriter, file); err != nil {
		return nil, err
	}

	// We have all the data written to the bodyWriter.
	// Now we can infer the content type
	contentType := bodyWriter.FormDataContentType()

	// This is mandatory as it flushes the buffer.
	bodyWriter.Close()

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return req, nil
}

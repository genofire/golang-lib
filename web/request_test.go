package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONRequest(t *testing.T) {
	assert := assert.New(t)

	data := struct {
		IP string `json:"ip"`
	}{}
	err := JSONRequest("http://ip.jsontest.com/", &data)
	assert.NoError(err)
	assert.NotEqual("", data.IP)

	wrongData := ""
	err = JSONRequest("http://ip.jsontest.com/", &wrongData)
	assert.Error(err)

	wrongData = ""
	err = JSONRequest("http://no.example.org/", &wrongData)
	assert.Error(err)
}

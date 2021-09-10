package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONRequest(t *testing.T) {
	assert := assert.New(t)

	data := struct {
		IP string `json:"query"`
	}{}
	err := JSONRequest("http://ip-api.com/json/?fields=query", &data)
	assert.NoError(err)
	assert.NotEqual("", data.IP)

	wrongData := ""
	err = JSONRequest("http://ip-api.com/json/?fields=query", &wrongData)
	assert.Error(err)

	wrongData = ""
	err = JSONRequest("http://no.example.org/", &wrongData)
	assert.Error(err)
}

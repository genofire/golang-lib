package mailer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	assert := assert.New(t)

	value := struct {
		T string
	}{
		T: "<script>alert('you have been pwned')</script>",
	}

	templ := TemplateTXT("A {{ .T }")
	// invalid template
	assert.Nil(templ)

	templ = TemplateTXT("A {{ .T }}")
	// text template
	assert.Equal("", templ.Render(3))

	// text template
	assert.Equal("A <script>alert('you have been pwned')</script>", templ.Render(&value))

	templ = TemplateHTML("A {{ .T }")
	// invalid template
	assert.Nil(templ)

	templ = TemplateHTML("A {{ .T }}")
	// html template
	assert.Equal("", templ.Render(3))

	// html template
	assert.Equal("A &lt;script&gt;alert(&#39;you have been pwned&#39;)&lt;/script&gt;", templ.Render(&value))
}

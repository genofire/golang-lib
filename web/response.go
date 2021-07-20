package web

import (
	"github.com/gin-gonic/gin"
)

// MIME type strings.
const (
	ContentTypeJSON = "application/json"
	ContentTypeJS   = "application/javascript"
	ContentTypeXML  = "text/xml"
	ContentTypeYAML = "text/yaml"
	ContentTypeHTML = "text/html"
)

// Response sends an HTTP response.
//
// statusCode is the respone's status.
//
// If the request's Content-Type is JavaScript, JSON, YAML, or XML, it returns
// data serialized as JSONP, JSON, YAML, or XML, respectively. If the
// Content-Type is HTML, it returns the HTML template templateName rendered with
// data.
func Response(ctx *gin.Context, statusCode int, data interface{}, templateName string) {
	switch ctx.ContentType() {
	case ContentTypeJS:
		ctx.JSONP(statusCode, data)
		return
	case ContentTypeJSON:
		ctx.JSON(statusCode, data)
		return
	case ContentTypeYAML:
		ctx.YAML(statusCode, data)
		return
	case ContentTypeXML:
		ctx.XML(statusCode, data)
		return
	case ContentTypeHTML:
		ctx.HTML(statusCode, templateName, data)
		return
	default:
		ctx.JSON(statusCode, data)
		return
	}
}

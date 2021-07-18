package web

import (
	"github.com/gin-gonic/gin"
)

const (
	// ContentTypeJSON content type of json
	ContentTypeJSON = "application/json"
	// ContentTypeJS content type of javascript
	ContentTypeJS = "application/javascript"
	// ContentTypeXML content type of xml
	ContentTypeXML = "text/xml"
	// ContentTypeYAML content type of yaml
	ContentTypeYAML = "text/yaml"
	// ContentTypeHTML content type of html
	ContentTypeHTML = "text/html"
)

// Response give
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

package route

import (
	"github.com/IceWhaleTech/CasaOS/web"
	"github.com/gin-gonic/gin"
	"html/template"
)

func WebUIHome(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	index, _ := template.ParseFS(web.Static, "index.html")
	index.Execute(c.Writer, nil)
	return
}

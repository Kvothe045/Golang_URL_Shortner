// url_shortner/routers/router.go
package routers

import (
	"url_shortner/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db handlers.URLStore) *gin.Engine {
	r := gin.Default()

	r.POST("/encode", handlers.HandleEncodingURL(db))
	r.GET("/:code", handlers.HandleDecodingURL(db))

	return r
}

package api

import (
	v1 "app/api/v1"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApplyRoutes(app *gin.Engine) {
	api := app.Group("api")
	{
		v1.ApplyRoutes(api)
	}
	app.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})
}

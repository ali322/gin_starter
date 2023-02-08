package middleware

import (
	"app/lib"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		errs := c.Errors
		if len(errs) == 0 {
			return
		}
		last := errs.Last()
		err, ok := last.Err.(validator.ValidationErrors)
		if ok {
			c.JSON(http.StatusOK, &gin.H{
				"code": -2, "msg": lib.TranslateValidatorErrors(err),
			})
		} else {
			c.JSON(http.StatusOK, &gin.H{
				"code": -2, "msg": last.Error(),
			})
		}
		c.Abort()
	}
}

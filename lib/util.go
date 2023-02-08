package lib

import (
	"github.com/gin-gonic/gin"
)

func Reply(data any) *gin.H {
	return &gin.H{
		"code": 0, "data": data,
	}
}

func Reject(code int, data any) *gin.H {
	return &gin.H{
		"code": code, "msg": data,
	}
}

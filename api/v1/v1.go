package v1

import (
	"app/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("v1")
	{
		v1.Use(middleware.JWT(map[string]string{
			"ping":   "get",
			"public": "post|get",
		}))
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, &gin.H{
				"code": 0, "message": "pong",
			})
		})
		v1.POST("public/register", register)
		v1.POST("public/login", login)
		v1.POST("change/password", changePassword)
		v1.POST("reset/:id/password", resetPassword)
		v1.GET("public/message", messager)

		v1.GET("public/user", users)
		v1.GET("public/user/:id", user)
		v1.PUT("user/:id", updateUser)
		v1.DELETE("user/:id", deleteUser)
		v1.POST("active/user", activeUser)
		v1.POST("deactive/user", deactiveUser)
		v1.GET("me", me)

		v1.POST("follow/user", follow)
		v1.DELETE("follow/user", unfollow)
		v1.GET("user/:id/fans", fans)
		v1.GET("user/:id/following", followings)

		v1.POST("role", createRole)
		v1.GET("public/role", roles)
		v1.PUT("role/:id", updateRole)
		v1.DELETE("role/:id", deleteRole)
		v1.GET("public/role/:id", role)
		v1.POST("user/role", grantRole)
		v1.DELETE("user/role", revokeRole)
		v1.PUT("user/role", changeRole)
		v1.POST("active/role", activeRole)
		v1.DELETE("active/role", deactiveRole)

		v1.POST("action-category", createActionCategory)
		v1.PUT("action-category/:id", updateActionCategory)
		v1.GET("public/action-category/:id", actionCategory)
		v1.GET("public/action-category", actionCategories)
		v1.DELETE("action-category/:id", deleteActionCategory)

		v1.POST("action", createAction)
		v1.GET("public/action", actions)
		v1.PUT("action/:id", updateAction)
		v1.DELETE("action/:id", deleteAction)
		v1.GET("public/action/:id", action)
		v1.POST("role/action", grantAction)
		v1.DELETE("role/action", revokeAction)
		v1.PUT("role/action", changeAction)
	}
}

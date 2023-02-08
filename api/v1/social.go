package v1

import (
	"app/lib"
	"app/repository/dao"
	"app/repository/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func follow(c *gin.Context) {
	var body dto.ToggleFollow
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	auth := c.GetStringMap("auth")
	id := auth["id"].(string)
	me, err := body.Follow(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(me))
}

func unfollow(c *gin.Context) {
	var body dto.ToggleFollow
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	auth := c.GetStringMap("auth")
	id := auth["id"].(string)
	me, err := body.Unfollow(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(me))
}

func fans(c *gin.Context) {
	id := c.Param("id")
	user, err := dao.FindUser(id, map[string]interface{}{
		"preload": []string{"Fans"},
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(user))
}

func followings(c *gin.Context) {
	id := c.Param("id")
	user, err := dao.FindUser(id, map[string]interface{}{
		"preload": []string{"Followings"},
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(user))
}

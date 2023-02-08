package v1

import (
	"app/lib"
	"app/repository/dao"
	"app/repository/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func users(c *gin.Context) {
	var query dto.QueryUser
	if err := c.ShouldBind(&query); err != nil {
		_ = c.Error(err)
		return
	}
	rows, count, err := query.Find()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(map[string]interface{}{
		"count": count,
		"rows":  rows,
	}))
}

func user(c *gin.Context) {
	id := c.Param("id")
	user, err := dao.FindUser(id, map[string]interface{}{
		"preload": []string{"Group", "Role", "Role.Actions"},
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(user))
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var body dto.UpdateUser
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	updated, err := body.Save(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(updated))
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	deleted, err := dao.DeleteUser(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(deleted))
}

func activeUser(c *gin.Context) {
	var body dto.ToggleUserActive
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	err := body.Active()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(nil))
}

func deactiveUser(c *gin.Context) {
	var body dto.ToggleUserActive
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	err := body.Deactive()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(nil))
}

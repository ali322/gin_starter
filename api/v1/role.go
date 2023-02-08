package v1

import (
	"app/lib"
	"app/repository/dao"
	"app/repository/dto"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func createRole(c *gin.Context) {
	var body dto.NewRole
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	exists, _ := dao.RoleExistsByName(body.Name)
	if exists {
		_ = c.Error(errors.New("角色已存在"))
		return
	}
	created, err := body.Create()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(created))
}

func updateRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	var body dto.UpdateRole
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	saved, err := body.Save(uint(id))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(saved))
}

func deleteRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	exists, found := dao.RoleExists(uint(id))
	if !exists {
		_ = c.Error(errors.New("角色不存在"))
		return
	}
	err = found.Delete()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(nil))
}

func role(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	found, err := dao.FindRole(uint(id), map[string]interface{}{
		"preload": []string{"Users", "Actions", "Actions.Category"},
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(found))
}

func roles(c *gin.Context) {
	var query dto.QueryRole
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

func grantRole(c *gin.Context) {
	var body dto.OPRole
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	err := body.Grant()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(nil))
}

func revokeRole(c *gin.Context) {
	var body dto.OPRole
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	err := body.Revoke()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(nil))
}

func changeRole(c *gin.Context) {
	var body dto.OPRole
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	err := body.Change()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(nil))
}

func activeRole(c *gin.Context) {
	var body dto.ToggleRoleActive
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

func deactiveRole(c *gin.Context) {
	var body dto.ToggleRoleActive
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

package v1

import (
	"app/lib"
	"app/lib/config"
	"app/repository/dao"
	"app/repository/dto"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func register(c *gin.Context) {
	var body dto.RegisterUser
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	exists, _ := dao.FindByUsername(body.Username)
	if exists {
		_ = c.Error(errors.New("用户已存在"))
		return
	}
	defaultRoleID, err := strconv.Atoi(config.App.DefaultRole)
	if err != nil {
		_ = c.Error(err)
		return
	}
	created, err := body.Create(uint(defaultRoleID))
	if err != nil {
		_ = c.Error(err)
		return
	}
	token, err := lib.GenerateJWTToken(config.App.JWTSecret, map[string]interface{}{
		"id": created.ID, "username": created.Username, "roleID": defaultRoleID,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(map[string]interface{}{
		"user": created, "token": token,
	}))
}

func login(c *gin.Context) {
	var body dto.LoginUser
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	defaultRoleID, err := strconv.Atoi(config.App.DefaultRole)
	if err != nil {
		_ = c.Error(err)
		return
	}
	found, err := body.Login(uint(defaultRoleID))
	if err != nil {
		_ = c.Error(err)
		return
	}
	if !found.IsActived {
		_ = c.Error(errors.New("用户未激活"))
		return
	}
	token, err := lib.GenerateJWTToken(config.App.JWTSecret, map[string]interface{}{
		"id": found.ID, "username": found.Username, "roleID": found.RoleID,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(map[string]interface{}{
		"user": found, "token": token,
	}))
}

func changePassword(c *gin.Context) {
	var body dto.ChangePassword
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	auth := c.GetStringMap("auth")
	id := auth["id"].(string)
	updated, err := body.ChangePassword(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(updated))
}

func resetPassword(c *gin.Context) {
	var body dto.ResetPassword
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	id := c.Param("id")
	updated, err := body.ResetPassword(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(updated))
}

func me(c *gin.Context) {
	auth := c.GetStringMap("auth")
	id := auth["id"].(string)
	user, err := dao.FindUser(id, map[string]interface{}{
		"preload": []string{"Group"},
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	role, err := dao.FindRole(*user.RoleID, map[string]interface{}{
		"preload": []string{"Actions"},
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	defaultRoleID, err := strconv.Atoi(config.App.DefaultRole)
	if err != nil {
		_ = c.Error(err)
		return
	}
	defaultRole, err := dao.FindRole(uint(defaultRoleID), map[string]interface{}{
		"preload": []string{"Actions"},
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(map[string]interface{}{
		"user": user, "currentRole": role, "defaultRole": defaultRole,
	}))
}

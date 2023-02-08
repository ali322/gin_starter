package v1

import (
	"app/lib"
	"app/repository/dao"
	"app/repository/dto"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createAction(c *gin.Context) {
	var body dto.NewAction
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	created, err := body.Create()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(created))
}

func updateAction(c *gin.Context) {
	id := c.Param("id")
	var body dto.UpdateAction
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	saved, err := body.Save(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(saved))
}

func deleteAction(c *gin.Context) {
	id := c.Param("id")
	exists, found := dao.ActionExists(id)
	if !exists {
		_ = c.Error(errors.New("行为不存在"))
		return
	}
	err := found.Delete()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(nil))
}

func action(c *gin.Context) {
	id := c.Param("id")
	found, err := dao.FindAction(id, map[string]interface{}{
		"preload": []string{"Category"},
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(found))
}

func actions(c *gin.Context) {
	var query dto.QueryAction
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

func grantAction(c *gin.Context) {
	var body dto.OPAction
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

func revokeAction(c *gin.Context) {
	var body dto.OPAction
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

func changeAction(c *gin.Context) {
	var body dto.OPAction
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

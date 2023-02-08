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

func createActionCategory(c *gin.Context) {
	var body dto.NewActionCategory
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

func updateActionCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	var body dto.UpdateActionCategory
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

func deleteActionCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	exists, found := dao.ActionCategoryExists(uint(id))
	if !exists {
		_ = c.Error(errors.New("权限分类不存在"))
		return
	}
	err = found.Delete()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(nil))
}

func actionCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	found, err := dao.FindActionCategory(uint(id), map[string]interface{}{
		"preload": []string{"Actions"},
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(found))
}

func actionCategories(c *gin.Context) {
	rows, err := dao.FindActionCategories(map[string]interface{}{
		"preload": []string{"Actions"},
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, lib.Reply(rows))
}

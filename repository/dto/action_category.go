package dto

import (
	"app/repository/dao"
	"errors"

	"gorm.io/gorm"
)

type NewActionCategory struct {
	Name        string `binding:"omitempty,lt=200" json:"name"`
	Description string `json:"description"`
}

func (body *NewActionCategory) Create() (dao.ActionCategory, error) {
	m := dao.ActionCategory{
		Name: body.Name, Description: body.Description,
	}
	return m.Create()
}

type UpdateActionCategory struct {
	Name        string `binding:"omitempty,lt=200" json:"name"`
	Description string `json:"description"`
}

func (body *UpdateActionCategory) Save(id uint) (dao.ActionCategory, error) {
	m, err := dao.FindActionCategory(id, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return m, errors.New("权限分类不存在")
		} else {
			return m, err
		}
	}
	values := map[string]interface{}{
		"name":        body.Name,
		"description": body.Description,
	}
	values = omitEmpty(values)
	return m.Update(values)
}

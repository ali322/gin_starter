package dto

import (
	"app/repository/dao"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type NewAction struct {
	Name        string `binding:"omitempty,lt=200" json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
	CategoryID  uint   `binding:"required,numeric,gt=0" json:"categoryID"`
}

func (body *NewAction) Create() (dao.Action, error) {
	m := dao.Action{
		Name: body.Name, Description: body.Description, Value: body.Value, CategoryID: body.CategoryID,
	}
	exists, _ := dao.ActionCategoryExists(body.CategoryID)
	if !exists {
		return m, errors.New("权限分类不存在")
	}
	return m.Create()
}

type UpdateAction struct {
	Name        string `binding:"omitempty,lt=200" json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
	CategoryID  *uint  `binding:"omitempty,numeric,gt=0" json:"categoryID"`
	IsActived   bool   `binding:"omitempty" json:"isActived"`
}

func (body *UpdateAction) Save(id string) (dao.Action, error) {
	m, err := dao.FindAction(id, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return m, errors.New("权限不存在")
		} else {
			return m, err
		}
	}
	if body.CategoryID != nil {
		exists, _ := dao.ActionCategoryExists(*body.CategoryID)
		if !exists {
			return m, errors.New("权限分类不存在")
		}
	}
	values := map[string]interface{}{
		"name":        body.Name,
		"description": body.Description,
		"value":       body.Value,
		"category_id": body.CategoryID,
		"is_actived":  body.IsActived,
	}
	values = omitEmpty(values)
	return m.Update(values)
}

type QueryAction struct {
	Key       string `form:"key" binding:"max=10"`
	Page      int    `form:"page,default=1" binding:"min=1" json:"page"`
	Limit     int    `form:"limit,default=10" binding:"min=1" json:"limit"`
	SortBy    string `form:"sortBy,default=created_at" binding:"oneof=created_at updated_at last_logined_at" json:"sortBy"`
	SortOrder string `form:"sortOrder,default=desc" binding:"oneof=asc desc" json:"sortOrder"`
}

func (query *QueryAction) Find() ([]dao.Action, int64, error) {
	where := make([][]interface{}, 0)
	if query.Key != "" {
		where = append(where, []interface{}{"name LIKE ?", fmt.Sprintf("%%%s%%", query.Key)})
	}
	return dao.FindAndCountActions(map[string]interface{}{
		"where":   where,
		"preload": []string{"Category"},
		"offset":  (query.Page - 1) * query.Limit,
		"limit":   query.Limit,
		"order":   fmt.Sprintf("%s %s", query.SortBy, query.SortOrder),
	})
}

func isActionExist(action dao.Action, actions []dao.Action) bool {
	for i := 0; i < len(actions); i++ {
		if action.ID == actions[i].ID {
			return true
		}
	}
	return false
}

type OPAction struct {
	RoleID   uint   `uri:"roleID" json:"roleID"`
	ActionID string `uri:"actionID" json:"actionID"`
}

func (body OPAction) Grant() (err error) {
	role, err := dao.FindRole(body.RoleID, map[string]interface{}{
		"preload": []string{"Actions"},
	})
	if err != nil {
		return err
	}
	actions, err := dao.FindActions(map[string]interface{}{
		"where": strings.Split(body.ActionID, ","),
	})
	if err != nil {
		return err
	}
	var next []dao.Action
	for _, action := range actions {
		if !isActionExist(action, role.Actions) {
			next = append(next, action)
		}
	}
	return role.Relations("Actions").Append(next)
}

func (body OPAction) Revoke() (err error) {
	role, err := dao.FindRole(body.RoleID, map[string]interface{}{
		"preload": []string{"Actions"},
	})
	if err != nil {
		return err
	}
	actions, err := dao.FindActions(map[string]interface{}{
		"where": strings.Split(body.ActionID, ","),
	})
	if err != nil {
		return err
	}
	var next []dao.Action
	for _, action := range actions {
		if isActionExist(action, role.Actions) {
			next = append(next, action)
		}
	}
	return role.Relations("Actions").Delete(next)
}

func (body OPAction) Change() (err error) {
	role, err := dao.FindRole(body.RoleID, map[string]interface{}{
		"preload": []string{"Actions"},
	})
	if err != nil {
		return err
	}
	actions, err := dao.FindActions(map[string]interface{}{
		"where": strings.Split(body.ActionID, ","),
	})
	if err != nil {
		return err
	}
	var next []dao.Action
	next = append(next, actions...)
	return role.Relations("Actions").Replace(next)
}

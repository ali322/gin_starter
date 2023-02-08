package dto

import (
	"app/repository/dao"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type NewRole struct {
	Name        string `binding:"omitempty,lt=200" json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
	IsDefault   bool   `binding:"omitempty" json:"isDefault"`
	ActionID    string `binding:"omitempty" json:"actionID"`
}

func (body *NewRole) Create() (dao.Role, error) {
	m := dao.Role{
		Name: body.Name, Description: body.Description, IsDefault: body.IsDefault, Code: body.Code,
	}
	if body.ActionID != "" {
		actions, err := dao.FindActions(map[string]interface{}{
			"where": strings.Split(body.ActionID, ","),
		})
		if err != nil {
			return m, err
		}
		return m.Create(actions)
	}
	return m.Create(nil)
}

type UpdateRole struct {
	Name        string  `binding:"omitempty,lt=200" json:"name"`
	Description string  `json:"description"`
	Code        string  `json:"code"`
	IsDefault   *bool   `binding:"omitempty" json:"isDefault"`
	IsActived   *bool   `binding:"omitempty" json:"isActived"`
	ActionID    *string `binding:"omitempty" json:"actionID"`
}

func (body *UpdateRole) Save(id uint) (dao.Role, error) {
	m, err := dao.FindRole(id, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return m, errors.New("角色不存在")
		} else {
			return m, err
		}
	}
	values := map[string]interface{}{
		"name":        body.Name,
		"description": body.Description,
		"code":        body.Code,
		// "is_default":  body.IsDefault,
		// "is_actived":  m.IsActived,
	}
	if body.IsDefault != nil {
		values["is_default"] = body.IsDefault
	}
	if body.IsActived != nil {
		values["is_actived"] = body.IsActived
	}
	values = omitEmpty(values)
	if body.ActionID != nil {
		actions, err := dao.FindActions(map[string]interface{}{
			"where": strings.Split(*body.ActionID, ","),
		})
		if err != nil {
			return m, err
		}
		return m.Update(values, actions)
	}
	return m.Update(values, []dao.Action{})
}

type QueryRole struct {
	Key       string `form:"key" binding:"max=10" json:"key"`
	IsDefault *uint  `form:"isDefault" binding:"omitempty,oneof=0 1" json:"isDefault"`
	IsActived *uint  `form:"isActived" binding:"omitempty,oneof=0 1" json:"isActived"`
	Page      int    `form:"page,default=1" binding:"min=1" json:"page"`
	Limit     int    `form:"limit,default=10" binding:"min=1" json:"limit"`
	SortBy    string `form:"sortBy,default=created_at" binding:"oneof=created_at updated_at last_logined_at" json:"sortBy"`
	SortOrder string `form:"sortOrder,default=desc" binding:"oneof=asc desc" json:"sortOrder"`
}

func (query *QueryRole) Find() ([]dao.Role, int64, error) {
	where := make([][]interface{}, 0)
	if query.Key != "" {
		where = append(where, []interface{}{"name LIKE ?", fmt.Sprintf("%%%s%%", query.Key)})
	}
	if query.IsDefault != nil {
		where = append(where, []interface{}{"is_default = ?", query.IsDefault})
	}
	if query.IsActived != nil {
		where = append(where, []interface{}{"is_actived = ?", query.IsActived})
	}
	return dao.FindAndCountRoles(map[string]interface{}{
		"where":   where,
		"preload": []string{"Actions"},
		"offset":  (query.Page - 1) * query.Limit,
		"limit":   query.Limit,
		"order":   fmt.Sprintf("%s %s", query.SortBy, query.SortOrder),
	})
}

type OPRole struct {
	UserID string `json:"userID"`
	RoleID uint   `json:"roleID"`
}

func isUserExist(row dao.User, rows []dao.User) bool {
	for i := 0; i < len(rows); i++ {
		if row.ID == rows[i].ID {
			return true
		}
	}
	return false
}

func (body OPRole) Grant() (err error) {
	role, err := dao.FindRole(body.RoleID, map[string]interface{}{
		"preload": []string{"Users"},
	})
	if err != nil {
		return err
	}
	users, err := dao.FindUsers(map[string]interface{}{
		"where": strings.Split(body.UserID, ","),
	})
	if err != nil {
		return err
	}
	var next []dao.User
	for _, user := range users {
		if !isUserExist(user, role.Users) {
			next = append(next, user)
		}
	}
	return role.Relations("Users").Append(next)
}

func (body OPRole) Revoke() (err error) {
	role, err := dao.FindRole(body.RoleID, map[string]interface{}{
		"preload": []string{"Users"},
	})
	if err != nil {
		return err
	}
	users, err := dao.FindUsers(map[string]interface{}{
		"where": strings.Split(body.UserID, ","),
	})
	if err != nil {
		return err
	}
	var next []dao.User
	for _, user := range users {
		if isUserExist(user, role.Users) {
			next = append(next, user)
		}
	}
	return role.Relations("Users").Delete(next)
}

func (body OPRole) Change() (err error) {
	var next []dao.User
	role, err := dao.FindRole(body.RoleID, map[string]interface{}{
		"preload": []string{"Users"},
	})
	if err != nil {
		return err
	}
	users, err := dao.FindUsers(map[string]interface{}{
		"where": strings.Split(body.UserID, ","),
	})
	if err != nil {
		return err
	}
	next = append(next, users...)
	return role.Relations("Users").Replace(next)
}

type ToggleRoleActive struct {
	RoleID string `binding:"required" json:"roleID"`
}

func (body ToggleRoleActive) Active() (err error) {
	values := map[string]interface{}{
		"is_actived": true,
	}
	return dao.UpdateRoles(values, strings.Split(body.RoleID, ","))
}

func (body ToggleRoleActive) Deactive() (err error) {
	values := map[string]interface{}{
		"is_actived": false,
	}
	return dao.UpdateRoles(values, strings.Split(body.RoleID, ","))
}

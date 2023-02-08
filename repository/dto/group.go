package dto

import (
	"app/repository/dao"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type NewGroup struct {
	Name        string `binding:"omitempty,lt=200" json:"name"`
	Description string `json:"description"`
	Size        string `json:"size"`
	Logo        string `json:"logo"`
	IndustryID  uint   `binding:"omitempty,numeric,gt=0" json:"industryID"`
}

func (body *NewGroup) Create(user *dao.User, role *dao.Role) (dao.Group, error) {
	m := dao.Group{
		Name: body.Name, Description: body.Description, Size: body.Size, Logo: body.Logo,
	}
	return m.Create(user, role)
}

type UpdateGroup struct {
	Name        string `binding:"omitempty,lt=200" json:"name"`
	Description string `json:"description"`
	Size        string `json:"size"`
	Logo        string `json:"logo"`
	IndustryID  uint   `binding:"omitempty,numeric,gt=0" json:"industryID"`
}

func (body *UpdateGroup) Save(id string) (dao.Group, error) {
	m, err := dao.FindGroup(id, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return m, errors.New("团队不存在")
		} else {
			return m, err
		}
	}
	values := map[string]interface{}{
		"name":        body.Name,
		"description": body.Description,
		"size":        body.Size,
		"logo":        body.Logo,
		"industry_id": body.IndustryID,
	}
	values = omitEmpty(values)
	return m.Update(values)
}

type QueryGroup struct {
	Key       string `form:"key" binding:"max=10"`
	Page      int    `form:"page,default=1" binding:"min=1" json:"page"`
	Limit     int    `form:"limit,default=10" binding:"min=1" json:"limit"`
	SortBy    string `form:"sortBy,default=created_at" binding:"oneof=created_at updated_at last_logined_at" json:"sortBy"`
	SortOrder string `form:"sortOrder,default=desc" binding:"oneof=asc desc" json:"sortOrder"`
}

func (query *QueryGroup) Find() ([]dao.Group, int64, error) {
	where := make([][]interface{}, 0)
	if query.Key != "" {
		where = append(where, []interface{}{"name LIKE ?", fmt.Sprintf("%%%s%%", query.Key)})
	}
	return dao.FindAndCountGroups(map[string]interface{}{
		"where": where,
		// "preload": []string{"Role", "AssetFolder"},
		"offset": (query.Page - 1) * query.Limit,
		"limit":  query.Limit,
		"order":  fmt.Sprintf("%s %s", query.SortBy, query.SortOrder),
	})
}

type DeleteGroup struct {
	ID string `binding:"omitempty" json:"id"`
}

func (body *DeleteGroup) Delete(defaultRole uint) (err error) {
	return dao.DeleteGroup(strings.Split(body.ID, ","), defaultRole)
}

type IOGroup struct {
	GroupID string `binding:"required" json:"groupID"`
	UserID  string `binding:"required" json:"userID"`
}

func (body *IOGroup) In() (dao.Group, error) {
	group, err := dao.FindGroup(body.GroupID, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return group, errors.New("团队不存在")
		} else {
			return group, err
		}
	}
	err = group.AddUsers(strings.Split(body.UserID, ","))
	if err != nil {
		return group, err
	}
	return group, nil
}

func (body *IOGroup) Out(defaultRoleID uint) (dao.Group, error) {
	group, err := dao.FindGroup(body.GroupID, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return group, errors.New("团队不存在")
		} else {
			return group, err
		}
	}
	userIDs := strings.Split(body.UserID, ",")
	err = group.RemoveUsers(userIDs, defaultRoleID)
	if err != nil {
		return group, err
	}
	return group, nil
}

package dao

import (
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Action struct {
	BaseModel
	ID          string          `gorm:"size:100;not null;primaryKey" json:"id"`
	Name        string          `gorm:"size:200;not null" json:"name"`
	Description string          `gorm:"type:text" json:"description"`
	Value       string          `gorm:"type:text" json:"value"`
	CategoryID  uint            `form:"categoryID" json:"categoryID"`
	Category    *ActionCategory `gorm:"foreignkey:CategoryID" binding:"-" json:"category,omitempty"`
	IsActived   bool            `gorm:"type:boolean;default:false" binding:"boolean" json:"isActived"`
	Roles       []Role          `gorm:"many2many:role_has_actions" binding:"-" json:"roles"`
}

func (m Action) Create() (Action, error) {
	id := uuid.NewV4().String()
	m.ID = id
	if err := db.Create(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

func (m Action) Update(values interface{}) (Action, error) {
	err := db.Model(&m).Updates(values).Error
	return m, err
}

func FindAction(id string, options map[string]interface{}) (Action, error) {
	var one Action
	if err := db.Scopes(applyQueryOptions(options)).First(&one, "id = ?", id).Error; err != nil {
		return one, err
	}
	return one, nil
}

func FindActions(options map[string]interface{}) ([]Action, error) {
	var rows []Action
	if err := db.Scopes(applyQueryOptions(options)).Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

func FindAndCountActions(options map[string]interface{}) ([]Action, int64, error) {
	var rows []Action
	var count int64
	if err := db.Scopes(applyQueryOptions(options)).Find(&rows).Error; err != nil {
		return rows, count, err
	}
	delete(options, "offset")
	delete(options, "limit")
	delete(options, "order")
	delete(options, "join")
	if err := db.Model(&Action{}).Scopes(applyQueryOptions(options)).Count(&count).Error; err != nil {
		return rows, count, err
	}
	return rows, count, nil
}

func ActionExists(id string) (bool, Action) {
	var one Action
	err := db.Where("id = ?", id).First(&one).Error
	notFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !notFound, one
}

func (m Action) Delete() error {
	// db.Model(&m).Association("Assets").Clear()
	return db.Delete(&m).Error
}

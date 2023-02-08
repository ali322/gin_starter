package dao

import (
	"errors"

	"gorm.io/gorm"
)

type ActionCategory struct {
	BaseModel
	Name        string   `gorm:"size:100" binding:"required,lt=100" json:"name"`
	Description string   `gorm:"type:text" binding:"required" json:"desc"`
	Actions     []Action `gorm:"foreignkey:CategoryID" binding:"-" json:"actions"`
}

func (m ActionCategory) Create() (ActionCategory, error) {
	if err := db.Create(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

func (m ActionCategory) Update(values interface{}) (ActionCategory, error) {
	err := db.Model(&m).Updates(values).Error
	return m, err
}

func FindActionCategory(id uint, options map[string]interface{}) (ActionCategory, error) {
	var one ActionCategory
	if err := db.Scopes(applyQueryOptions(options)).First(&one, "id = ?", id).Error; err != nil {
		return one, err
	}
	return one, nil
}

func FindActionCategories(options map[string]interface{}) ([]ActionCategory, error) {
	var rows []ActionCategory
	if err := db.Scopes(applyQueryOptions(options)).Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

func ActionCategoryExists(id uint) (bool, ActionCategory) {
	var one ActionCategory
	err := db.Where("id = ?", id).First(&one).Error
	notFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !notFound, one
}

func (m ActionCategory) Delete() error {
	db.Model(&m).Association("Actions").Clear()
	return db.Delete(&m).Error
}

package dao

import (
	"errors"

	"gorm.io/gorm"
)

type Role struct {
	BaseModel
	Name        string   `gorm:"size:200;uniqueIndex;not null" json:"name"`
	Description string   `gorm:"type:text" json:"description"`
	Code        string   `gorm:"type:text" json:"code"`
	IsDefault   bool     `gorm:"type:boolean;default:false" binding:"boolean" json:"isDefault"`
	IsActived   bool     `gorm:"type:boolean;default:true" binding:"boolean" json:"isActived"`
	Actions     []Action `gorm:"many2many:role_has_actions" binding:"-" json:"actions"`
	Users       []User   `gorm:"foreignkey:RoleID" binding:"-" json:"users"`
}

func (m *Role) AfterFind(tx *gorm.DB) (err error) {
	if !m.IsActived {
		m.Actions = []Action{}
	}
	return
}

func (m Role) Create(actions []Action) (Role, error) {
	tx := db.Begin()
	if err := tx.Create(&m).Error; err != nil {
		tx.Rollback()
		return m, err
	}
	err := tx.Model(&m).Association("Actions").Append(actions)
	if err != nil {
		tx.Rollback()
		return m, err
	}
	tx.Commit()
	return m, nil
}

func (m Role) Update(values interface{}, actions []Action) (Role, error) {
	tx := db.Begin()
	err := tx.Model(&m).Updates(values).Error
	if err != nil {
		tx.Rollback()
		return m, err
	}
	if len(actions) > 0 {
		err = tx.Model(&m).Association("Actions").Replace(actions)
	}
	if err != nil {
		tx.Rollback()
		return m, err
	}
	tx.Commit()
	return m, nil
}

func UpdateRoles(values interface{}, ids []string) error {
	return db.Model(&Role{}).Where("id IN (?)", ids).Updates(values).Error
}

func FindRole(id uint, options map[string]interface{}) (Role, error) {
	var one Role
	if err := db.Scopes(applyQueryOptions(options)).First(&one, "id = ?", id).Error; err != nil {
		return one, err
	}
	return one, nil
}

func RoleExistsByName(name string) (bool, Role) {
	var one Role
	err := db.First(&one, "name = ?", name).Error
	notFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !notFound, one
}

func FindRoles(options map[string]interface{}) ([]Role, error) {
	var rows []Role
	if err := db.Scopes(applyQueryOptions(options)).Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

func FindAndCountRoles(options map[string]interface{}) ([]Role, int64, error) {
	var rows []Role
	var count int64
	if err := db.Scopes(applyQueryOptions(options)).Find(&rows).Error; err != nil {
		return rows, count, err
	}
	delete(options, "offset")
	delete(options, "limit")
	delete(options, "order")
	delete(options, "join")
	if err := db.Model(&Role{}).Scopes(applyQueryOptions(options)).Count(&count).Error; err != nil {
		return rows, count, err
	}
	return rows, count, nil
}

func RoleExists(id uint) (bool, Role) {
	var one Role
	err := db.Where("id = ?", id).First(&one).Error
	notFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !notFound, one
}

func (m Role) Delete() error {
	// db.Model(&m).Association("Assets").Clear()
	return db.Delete(&m).Error
}

func (m Role) Relations(col string) *gorm.Association {
	return db.Model(&m).Association(col)
}

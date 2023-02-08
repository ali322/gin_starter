package dao

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Group struct {
	BaseModel
	ID          string `gorm:"size:100;not null;primaryKey" json:"id"`
	Name        string `gorm:"size:200;uniqueIndex;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Size        string `gorm:"type:text" json:"size"`
	Logo        string `gorm:"type:text" json:"logo"`
	Amount      uint   `gorm:"default:0" binding:"-" json:"amount"`
	Users       []User `binding:"-" json:"users"`
	OwnerID     string `json:"ownerID"`
	Owner       *User  `binding:"-" json:"owner"`
}

func (m *Group) AfterDelete(tx *gorm.DB) (err error) {
	return tx.Model(&Group{}).Where("id = ?", m.ID).Update("owner_id", gorm.Expr("NULL")).Error
}

func (m *Group) Create(user *User, role *Role) (Group, error) {
	id := uuid.NewV4().String()
	m.ID = id
	m.OwnerID = user.ID
	tx := db.Begin()
	if err := tx.Create(&m).Error; err != nil {
		tx.Rollback()
		return *m, err
	}
	err := tx.Model(&User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"role_id": role.ID, "group_id": m.ID,
	}).Error
	if err != nil {
		tx.Rollback()
		return *m, err
	}
	tx.Commit()
	return *m, nil
}

func (m Group) Update(values interface{}) (Group, error) {
	err := db.Model(&m).Updates(values).Error
	return m, err
}

func (m Group) Save() (Group, error) {
	if err := db.Save(m).Error; err != nil {
		return m, err
	}
	return m, nil
}

func FindGroup(id string, options map[string]interface{}) (Group, error) {
	var one Group
	if err := db.Scopes(applyQueryOptions(options)).First(&one, "id = ?", id).Error; err != nil {
		return one, err
	}
	return one, nil
}

func FindGroups(options map[string]interface{}) ([]Group, error) {
	var rows []Group
	if err := db.Scopes(applyQueryOptions(options)).Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

func FindAndCountGroups(options map[string]interface{}) ([]Group, int64, error) {
	var rows []Group
	var count int64
	if err := db.Scopes(applyQueryOptions(options)).Find(&rows).Error; err != nil {
		return rows, count, err
	}
	delete(options, "offset")
	delete(options, "limit")
	delete(options, "order")
	delete(options, "join")
	if err := db.Model(&Group{}).Scopes(applyQueryOptions(options)).Count(&count).Error; err != nil {
		return rows, count, err
	}
	if len(rows) > 0 {
		ownerIDs := make([]string, 0)
		for _, v := range rows {
			ownerIDs = append(ownerIDs, v.OwnerID)
		}
		var owners []User
		if err := db.Find(&owners, ownerIDs).Error; err != nil {
			return rows, count, err
		}
		for i, owner := range owners {
			rows[i].Owner = &owner
		}
	}
	return rows, count, nil
}

func GroupExists(id string) (bool, Group) {
	var one Group
	err := db.Where("id = ?", id).First(&one).Error
	notFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !notFound, one
}

func GroupExistsByName(name string) (bool, Group) {
	var one Group
	err := db.Where("name = ?", name).First(&one).Error
	notFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !notFound, one
}

func (m Group) Delete() error {
	return db.Delete(&m).Error
}

func (m Group) Relations(col string) *gorm.Association {
	return db.Model(&m).Association(col)
}

func GroupByDay(day uint) ([]map[string]interface{}, error) {
	all := make([]map[string]interface{}, 0)
	sql := "SELECT DATE_FORMAT(created_at,'%Y-%m-%d') AS createdDate,COUNT(*) AS count FROM groups WHERE deleted_at IS NULL GROUP BY createdDate"
	if err := db.Raw(sql).Scan(&all).Error; err != nil {
		return all, err
	}
	return all, nil
}

func GroupByMonth(month uint) ([]map[string]interface{}, error) {
	all := make([]map[string]interface{}, 0)
	sql := "SELECT DATE_FORMAT(created_at,'%Y-%m') AS createdDate,COUNT(*) AS count FROM groups WHERE deleted_at IS NULL GROUP BY createdDate"
	if err := db.Raw(sql).Scan(&all).Error; err != nil {
		return all, err
	}
	return all, nil
}

func GroupOfIndustry() ([]map[string]interface{}, error) {
	all := make([]map[string]interface{}, 0)
	if err := db.Model(&Group{}).Select("COUNT(*) AS count, groups.industry_id, industries.name as industry").Group("industry_id").Joins("LEFT JOIN industries ON industries.id = groups.industry_id").Scan(&all).Error; err != nil {
		return all, err
	}
	return all, nil
}

func DeleteGroup(id []string, defaultRole uint) (err error) {
	tx := db.Begin()
	var rows []Group
	err = tx.Unscoped().Find(&rows, id).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// user
	err = tx.Model(&User{}).Where("group_id IN (?)", id).Update("role_id", defaultRole).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// app of personal
	// err = tx.Where("group_id IN (?) AND type = ?", id, "personal").Delete(&App{}).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// // app of group
	// err = tx.Model(&App{}).Where("group_id IN (?) AND type = ?", id, "group").Updates(map[string]interface{}{
	// 	"owner_id": db.Model(&Group{}).Select("owner_id").Where("groups.id = apps.group_id"),
	// 	"user_id":  db.Model(&Group{}).Select("owner_id").Where("groups.id = apps.group_id"),
	// 	"type":     "personal",
	// }).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// // app folder of personal
	// err = tx.Where("group_id IN (?) AND type = ?", id, "personal").Delete(&AppFolder{}).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// // app folder of group
	// err = tx.Model(&AppFolder{}).Where("group_id IN (?) AND type = ?", id, "group").Updates(map[string]interface{}{
	// 	"owner_id": db.Model(&Group{}).Select("owner_id").Where("groups.id = app_folders.group_id"),
	// 	"user_id":  db.Model(&Group{}).Select("owner_id").Where("groups.id = app_folders.group_id"),
	// 	"type":     "personal",
	// }).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// // asset
	// err = tx.Model(&Asset{}).Where("group_id IN (?)", id).Updates(map[string]interface{}{
	// 	"owner_id": db.Model(&Group{}).Select("owner_id").Where("groups.id = assets.group_id"),
	// 	"user_id":  db.Model(&Group{}).Select("owner_id").Where("groups.id = assets.group_id"),
	// }).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// // asset folder
	// err = tx.Model(&AssetFolder{}).Where("group_id IN (?)", id).Updates(map[string]interface{}{
	// 	"owner_id": db.Model(&Group{}).Select("owner_id").Where("groups.id = asset_folders.group_id"),
	// 	"user_id":  db.Model(&Group{}).Select("owner_id").Where("groups.id = asset_folders.group_id"),
	// }).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	err = tx.Model(&rows).Association("Users").Clear()
	if err != nil {
		tx.Rollback()
		return err
	}
	// err = tx.Model(&rows).Association("Assets").Clear()
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// err = tx.Model(&rows).Association("AssetFolders").Clear()
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// err = tx.Model(&rows).Association("Apps").Clear()
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// err = tx.Model(&rows).Association("AppFolders").Clear()
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	err = tx.Delete(&Group{}, id).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return
}

func isIDExists(id string, ids []string) bool {
	for _, r := range ids {
		if r == id {
			return true
		}
	}
	return false
}

func (m *Group) AddUsers(id []string) (err error) {
	tx := db.Begin()
	var users []User
	err = tx.Where("id IN (?) AND group_id IS NULL", id).Find(&users).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(&m).Select("amount").Updates(map[string]interface{}{"amount": gorm.Expr("amount + ?", len(users))}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	m.Amount += uint(len(users))
	if len(users) > 0 {
		err = tx.Model(&m).Association("Users").Append(users)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return
}

func (m *Group) RemoveUsers(id []string, defaultRole uint) (err error) {
	tx := db.Begin()
	if isIDExists(m.OwnerID, id) {
		return fmt.Errorf("用户 %s 是团队管理员", m.OwnerID)
	}
	var users []User
	err = tx.Where("id IN (?) AND group_id = ?", id, m.ID).Find(&users).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// user
	err = tx.Model(&User{}).Where("id IN (?) AND group_id IN (?)", id, m.ID).Update("role_id", defaultRole).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// app of personal
	// err = tx.Where("group_id = ? AND user_id IN (?) AND type = ?", m.ID, id, "personal").Delete(&App{}).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// // app of group
	// err = tx.Model(&App{}).Where("group_id = ? AND user_id IN (?) AND type = ?", m.ID, id, "group").Updates(map[string]interface{}{
	// 	"owner_id": db.Model(&Group{}).Select("owner_id").Where("groups.id = apps.group_id"),
	// 	// "user_id":  db.Model(&Group{}).Select("owner_id").Where("groups.id = apps.group_id"),
	// 	// "type":     "personal",
	// }).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// // app folder of personal
	// err = tx.Where("group_id = ? AND user_id IN (?) AND type = ?", m.ID, id, "personal").Delete(&AppFolder{}).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// // app folder of group
	// err = tx.Model(&AppFolder{}).Where("group_id = ? AND user_id IN (?) AND type = ?", m.ID, id, "group").Updates(map[string]interface{}{
	// 	"owner_id": db.Model(&Group{}).Select("owner_id").Where("groups.id = app_folders.group_id"),
	// 	// "user_id":  db.Model(&Group{}).Select("owner_id").Where("groups.id = app_folders.group_id"),
	// 	// "type":     "personal",
	// }).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	if len(users) > 0 {
		err = tx.Model(&m).Select("amount").Updates(map[string]interface{}{"amount": gorm.Expr("amount - ?", len(users))}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		m.Amount = m.Amount - uint(len(users))
	}
	err = tx.Model(&m).Association("Users").Delete(&users)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return
}

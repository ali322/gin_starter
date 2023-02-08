package dao

import (
	"app/lib"
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	ID              string        `gorm:"size:100;not null;primaryKey" json:"id"`
	Username        string        `gorm:"size:100;uniqueIndex;not null;index:idx_username" json:"username"`
	Password        string        `gorm:"size:200,not null" json:"-"`
	Email           string        `gorm:"size:200" json:"email"`
	Nickname        string        `gorm:"size:200" json:"nickname"`
	Avatar          string        `gorm:"type:text" json:"avatar"`
	Gender          string        `gorm:"type:text" json:"gender"`
	Phone           string        `gorm:"type:text" json:"phone"`
	Industry        string        `gorm:"type:text" json:"industry"`
	Source          string        `gorm:"type:text" json:"source"`
	Memo            string        `gorm:"type:text" json:"memo"`
	FollowingAmount uint          `gorm:"default:0" binding:"-" json:"followingAmount"`
	FansAmount      uint          `gorm:"default:0" binding:"-" json:"fansAmount"`
	Fans            []User        `gorm:"many2many:user_has_fans;foreignKey:ID;references:ID;joinForeignKey:FanID;joinReferences:UserID" json:"fans"`
	Followings      []User        `gorm:"many2many:user_has_fans;foreignKey:ID;references:ID;joinForeignKey:UserID;joinReferences:FanID" json:"followings"`
	IsActived       bool          `gorm:"type:boolean;default:true" binding:"-" json:"isActived"`
	LastLoginedAt   lib.LocalTime `json:"lastLoginedAt"`
	RoleID          *uint         `json:"roleID"`
	Role            *Role         `gorm:"foreignkey:RoleID" binding:"-" json:"role,omitempty"`
	GroupID         *string       `gorm:"type:text" json:"groupID"`
	Group           *Group        `gorm:"foreignkey:GroupID" binding:"-" json:"group"`
}

func (m User) Create() (User, error) {
	id := uuid.NewV4().String()
	m.ID = id
	m.LastLoginedAt = lib.LocalTime{Time: time.Now()}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(m.Password), 4)
	if err != nil {
		return m, err
	}
	m.Password = string(hashedPassword)
	if err := db.Create(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

func (m User) Save(cols []string) (User, error) {
	tx := db.Session(&gorm.Session{})
	if len(cols) > 0 {
		for _, col := range cols {
			tx = tx.Select(col)
		}
	}
	if err := tx.Updates(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

func (m User) Update(values interface{}) (User, error) {
	err := db.Model(&m).Updates(values).Error
	return m, err
}

func UpdateUsers(values interface{}, ids []string) error {
	return db.Model(&User{}).Where("id IN (?)", ids).Updates(values).Error
}

func FindByUsername(username string) (bool, User) {
	var one User
	err := db.Where("username = ?", username).First(&one).Error
	notFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !notFound, one
}

func DeleteUser(id string) (User, error) {
	var one User
	if err := db.Find(&one, "id = ?", id).Error; err != nil {
		return one, err
	}
	err := db.Delete(&one).Error
	return one, err
}

func FindAndCountUsers(options map[string]interface{}) ([]User, int64, error) {
	var rows []User
	var count int64
	if err := db.Scopes(applyQueryOptions(options)).Find(&rows).Error; err != nil {
		return rows, count, err
	}
	delete(options, "offset")
	delete(options, "limit")
	delete(options, "order")
	delete(options, "join")
	if err := db.Model(&User{}).Scopes(applyQueryOptions(options)).Count(&count).Error; err != nil {
		return rows, count, err
	}
	return rows, count, nil
}

func FindUsers(options map[string]interface{}) ([]User, error) {
	var rows []User
	if err := db.Scopes(applyQueryOptions(options)).Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

func UserExists(id string) (bool, User) {
	var one User
	err := db.Where("id = ?", id).First(&one).Error
	notFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !notFound, one
}

func FindUser(id string, options map[string]interface{}) (User, error) {
	var one User
	if err := db.Scopes(applyQueryOptions(options)).First(&one, "id = ?", id).Error; err != nil {
		return one, err
	}
	return one, nil
}

func (m User) Relations(col string) *gorm.Association {
	return db.Model(&m).Association(col)
}

func (m *User) Follow(user User) (err error) {
	tx := db.Begin()
	err = tx.Model(m).Select("FollowingAmount").Updates(User{FollowingAmount: m.FollowingAmount + 1}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(&user).Select("FansAmount").Updates(User{FansAmount: user.FansAmount + 1}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// next := append(old, user)
	err = tx.Model(m).Association("Followings").Append(&user)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return
}

func (m *User) Unfollow(user User) (err error) {
	tx := db.Begin()
	if m.FollowingAmount > 0 {
		err = tx.Model(&m).Select("FollowingAmount").Updates(User{FollowingAmount: m.FollowingAmount - 1}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if user.FansAmount > 0 {
		err = tx.Model(&user).Select("FansAmount").Updates(User{FansAmount: user.FansAmount - 1}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if m.FollowingAmount == 0 {
		if err = tx.Model(&m).Association("Followings").Clear(); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if err = tx.Model(&m).Association("Followings").Delete(user); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return
}

func UserByDay(day uint) ([]map[string]interface{}, error) {
	all := make([]map[string]interface{}, 0)
	sql := "SELECT DATE_FORMAT(created_at,'%Y-%m-%d') AS createdDate,COUNT(*) AS count FROM users WHERE is_actived = ? GROUP BY createdDate"
	if err := db.Raw(sql, true).Scan(&all).Error; err != nil {
		return all, err
	}
	return all, nil
}

func UserByMonth(month uint) ([]map[string]interface{}, error) {
	all := make([]map[string]interface{}, 0)
	sql := "SELECT DATE_FORMAT(created_at,'%Y-%m') AS createdDate,COUNT(*) AS count FROM users WHERE is_actived = ? GROUP BY createdDate"
	if err := db.Raw(sql, true).Scan(&all).Error; err != nil {
		return all, err
	}
	return all, nil
}

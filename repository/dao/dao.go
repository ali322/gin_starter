package dao

import (
	"app/lib"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Init(dsn string) {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}
	// db.Debug().Logger
	db.AutoMigrate(&User{}, &Role{}, &Action{}, &ActionCategory{}, &Group{})
	// if err := initData(); err != nil {
	// 	log.Fatal(err)
	// }
}

func Close() error {
	d, err := db.DB()
	if err != nil {
		return err
	}
	return d.Close()
}

type BaseModel struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	CreatedAt lib.LocalTime `json:"createdAt"`
	UpdatedAt lib.LocalTime `json:"updatedAt"`
	DeletedAt lib.DeletedAt `gorm:"index" json:"deletedAt"`
}

func applyQueryOptions(options map[string]interface{}) func(db *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// tx := db.Session(&gorm.Session{})
		if options["preload"] != nil {
			if preload, ok := options["preload"].([]string); ok {
				for _, col := range preload {
					tx = tx.Preload(col)
				}
			}
			if preload, ok := options["preload"].(map[string]interface{}); ok {
				for key, val := range preload {
					if val == nil {
						tx = tx.Preload(key)
					} else {
						tx = tx.Preload(key, val)
					}
				}
			}
		}
		if options["select"] != nil {
			if selected, ok := options["select"].([]string); ok {
				tx = tx.Select(selected)
			}
		}
		if options["where"] != nil {
			switch options["where"].(type) {
			case []string, []uint:
				tx = tx.Where("id in (?)", options["where"])
			case [][]interface{}:
				for _, where := range options["where"].([][]interface{}) {
					tx = tx.Where(where[0], where[1:]...)
				}
			case map[string]interface{}:
				tx = tx.Where(options["where"])
			}
		}
		if options["join"] != nil {
			tx = tx.Joins(options["join"].(string))
		}
		if options["order"] != nil {
			if orders, ok := options["order"].([]string); ok {
				for _, order := range orders {
					tx = tx.Order(order)
				}
			} else {
				tx = tx.Order(options["order"])
			}
		}
		if options["offset"] != nil {
			tx = tx.Offset(options["offset"].(int))
		}
		if options["limit"] != nil {
			tx = tx.Limit(options["limit"].(int))
		}
		return tx
	}
}

func initData() error {
	newActionCategory := ActionCategory{
		Name: "基础权限",
	}
	actionCategory, err := newActionCategory.Create()
	if err != nil {
		return err
	}
	actions := []Action{
		{Name: "管理菜单可见", Value: "ADMIN_MENU_VISIBLE", IsActived: true, CategoryID: actionCategory.ID},
	}
	next := make([]Action, 0)
	for _, v := range actions {
		created, err := v.Create()
		if err != nil {
			return err
		}
		next = append(next, created)
	}
	role := Role{
		Name: "平台管理员", IsDefault: true, IsActived: true,
	}
	adminRole, err := role.Create(next)
	if err != nil {
		return err
	}
	roles := []Role{
		{Name: "团队管理员", IsDefault: true, IsActived: true},
		{Name: "普通成员", IsDefault: true, IsActived: true},
	}
	for _, v := range roles {
		_, err := v.Create(next)
		if err != nil {
			return err
		}
	}
	user := User{
		Username: "admin",
		Password: "321",
		Email:    "admin@live.com",
		RoleID:   &adminRole.ID,
	}
	_, err = user.Create()
	if err != nil {
		return err
	}
	return nil
}

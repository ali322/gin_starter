package lib

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type LocalTime struct {
	time.Time
}

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		return
	}
	timeStr := strings.Trim(string(data), "\"")
	now, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
	*t = LocalTime{now}
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	output := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(output), nil
}

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if zeroTime.UnixNano() == t.Time.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *LocalTime) Scan(v interface{}) (err error) {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalTime{Time: value}
		return
	}
	return fmt.Errorf("failed to convert %v to timestamp", v)
}

type DeletedAt sql.NullTime

// Scan implements the Scanner interface.
func (n *DeletedAt) Scan(value interface{}) error {
	return (*sql.NullTime)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n DeletedAt) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

func (n DeletedAt) MarshalJSON() ([]byte, error) {
	if n.Valid {
		output := fmt.Sprintf("\"%s\"", n.Time.Format("2006-01-02 15:04:05"))
		return []byte(output), nil
	}
	return json.Marshal(nil)
}

func (n *DeletedAt) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	timeStr := strings.Trim(string(b), "\"")
	now, err := time.Parse("2006-01-02 15:04:05", timeStr)
	*n = DeletedAt{Time: now}
	// err := json.Unmarshal(b, &n.Time)
	if err == nil {
		n.Valid = true
	}
	return err
}

func (DeletedAt) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{gorm.SoftDeleteQueryClause{Field: f}}
}

func (DeletedAt) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{gorm.SoftDeleteUpdateClause{Field: f}}
}

func (DeletedAt) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{gorm.SoftDeleteDeleteClause{Field: f}}
}

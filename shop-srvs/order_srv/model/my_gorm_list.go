package model

import (
	"database/sql/driver"
	"encoding/json"
)

// GormList 自定义gorm type类型
type GormList []string

func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

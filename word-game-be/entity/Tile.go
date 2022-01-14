package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Square struct {
	Tile  *Tile  `json:"tile,omitempty"`
	Bonus *Bonus `json:"bonus,omitempty"`
}

type Bonus struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type Tile struct {
	Letter    string `json:"letter"`
	Score     int    `json:"score"`
	Draggable bool   `json:"draggable,omitempty"`
}

// SquareSet exists to make GORM happy about using the JSON datatype
type SquareSet struct {
	Squares *[]Square `json:"squares"`
	Width   int       `json:"width,omitempty"`
	Height  int       `json:"height,omitempty"`
}

func (t *SquareSet) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Value is not jsonb", src))
	}

	return json.Unmarshal(bytes, &t)
}

func (t *SquareSet) GormDataType() string {
	return "json"
}

func (t *SquareSet) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *SquareSet) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

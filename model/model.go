package model

import "time"

const BatchSize = 500

type BaseModel struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id"`
}

type CommonTimestampsField struct {
	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at"`
}

package entity

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	CreatedAt time.Time `gorm:"column:created_at;CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"column:updated_at;CURRENT_TIMESTAMP"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	return nil
}

func (b *Base) BeforeUpdate(tx *gorm.DB) (err error) {
	b.UpdatedAt = time.Now()
	return nil
}

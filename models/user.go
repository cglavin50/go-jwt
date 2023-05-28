package models

import (
	"gorm.io/gorm"
)

// https://gorm.io/docs/models.html details field tags
// gorm.Model definition
// type Model struct {
// 	ID        uint           `gorm:"primaryKey"`
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt gorm.DeletedAt `gorm:"index"`
//   }

type User struct {
	gorm.Model
	Email    string `gorm:"unique"` // sets as a unique column
	Password string
}

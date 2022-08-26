package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title  string `json:"title" gorm:"type:VARCHAR(255);not null"`
	Author string `json:"author" gorm:"type:VARCHAR(255);not null"`
	Rating int    `json:"rating" gorm:"type:SMALLINT"`
}

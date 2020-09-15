package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title  string `josn:"title"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
}

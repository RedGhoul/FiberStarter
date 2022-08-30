package models

import "time"

type Book struct {
	ID        uint      `json:"Id"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	DeletedAt time.Time `json:"DeletedAt"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Rating    int       `json:"rating"`
}

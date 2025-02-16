package models

import (
	"time"
)

type Book struct {
	ID      uint      `json:"id" gorm:"primaryKey"`
	Title   string    `gorm:"title" json:"title"`
	Author  string    `gorm:"author" json:"author"`
	Price   float64   `gorm:"price" json:"price"`
	Created time.Time `gorm:"created" json:"created"`
	Updated time.Time `gorm:"updated" json:"updated"`
}

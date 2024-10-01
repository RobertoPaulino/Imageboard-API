package main

import "time"

type Post struct {
	ID        uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Username  string    `json:"username"`
	Body      string    `json:"body" gorm:"not null"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
}

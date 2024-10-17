package main

import "time"

type Post struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Body      string    `json:"body"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
}

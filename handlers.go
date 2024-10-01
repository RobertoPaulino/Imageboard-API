package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

func getPosts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, username, body, latitude, longitude, created_at FROM posts")

		if err != nil {
			http.Error(w, "Failed to get posts", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		var posts []Post
		for rows.Next() {
			var post Post
			err = rows.Scan(&post.ID, &post.Username, &post.Body, &post.Latitude, &post.Longitude, &post.CreatedAt)
			if err != nil {
				http.Error(w, "Failed to scan posts", http.StatusInternalServerError)
				return
			}
			posts = append(posts, post)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

// Handler for POST /posts to create a new post, coordinates are mandatory
func createPost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newPost Post

		// Parse JSON body into newPost struct
		if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil || newPost.Body == "" {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Ensure both latitude and longitude are provided
		if newPost.Latitude == 0 || newPost.Longitude == 0 {
			http.Error(w, "Latitude and Longitude are required", http.StatusBadRequest)
			return
		}

		// Set default username if not provided
		if newPost.Username == "" {
			newPost.Username = "Anonymous"
		}

		// Round coordinates to 3 decimal places
		newPost.Latitude = roundToThreeDecimals(newPost.Latitude)
		newPost.Longitude = roundToThreeDecimals(newPost.Longitude)

		// Insert the new post into the database with current time and coordinates
		result, err := db.Exec("INSERT INTO posts (username, body, latitude, longitude, created_at) VALUES (?, ?, ?, ?, ?)",
			newPost.Username, newPost.Body, newPost.Latitude, newPost.Longitude, time.Now())
		if err != nil {
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		// Retrieve the last inserted ID
		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Failed to retrieve post ID", http.StatusInternalServerError)
			return
		}

		// Create the post response to return, with the auto-incremented ID and current timestamp
		newPost.ID = uint(id)
		newPost.CreatedAt = time.Now()

		// Return the newly created post as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newPost)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Initialize the database
	db := initDB()
	defer db.Close()

	//create a limiter
	limiter := NewRateLimiter()

	// Define routes
	http.HandleFunc("/posts", rateLimitMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				getPosts(db)(w, r)
			case "POST":
				createPost(db)(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}, limiter))

	// Start the server
	fmt.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

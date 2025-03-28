package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	database "social-network/database"
	handlers "social-network/handlers"

	"github.com/rs/cors"
)

func main() {
	//  CORS  localhost:3000
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// Initialize the database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer database.DB.Close()

	fileServer := http.FileServer(http.Dir("./app"))
	http.Handle("/app/", http.StripPrefix("/app", fileServer))

	http.HandleFunc("/", handlers.GetPostsHandler)
	http.HandleFunc("/session", handlers.SessionHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/profile", handlers.ProfileHandler)
	http.HandleFunc("/new_post", handlers.CreatePostHandler)
	http.HandleFunc("/post", handlers.PostHandler)
	http.HandleFunc("/comment", handlers.CreateCommentHandler)
	http.HandleFunc("/like", handlers.LikeHandler)
	http.HandleFunc("/follow", handlers.FollowHandler)
	http.HandleFunc("/unfollow", handlers.UnfollowHandler)
	http.HandleFunc("/followers", handlers.FollowersHandler)
	http.HandleFunc("/following", handlers.FollowingHandler)
	// http.HandleFunc("/search", handlers.SearchHandler)
	// http.HandleFunc("/ws", handlers.WebSocketHandler)

	log.Println("Server started on :8404")
	fmt.Println("http://localhost:8404/")
	err := http.ListenAndServe(":8404", c.Handler(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	time.Sleep(10 * time.Second)
}

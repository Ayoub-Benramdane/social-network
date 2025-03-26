package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	database "social-network/backend/database"
	handlers "social-network/backend/handlers"

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

	// http.HandleFunc("/", handlers.HomePage)
	// http.HandleFunc("/show_posts", handlers.ShowPosts)
	// http.HandleFunc("/post_submit", handlers.PostSubmit)
	// http.HandleFunc("/comment_submit", handlers.CommentSubmit)
	// http.HandleFunc("/interact", handlers.HandleInteract)
	// http.HandleFunc("/get_categories", handlers.GetCategories)
	// http.HandleFunc("/Connections", handlers.Connections)

	http.HandleFunc("/login", handlers.LoginHandler)
	// http.HandleFunc("/check-session", auth.CheckSessionHandler)
	// http.HandleFunc("/logout", auth.LogoutHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)

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

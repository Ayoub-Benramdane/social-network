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
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// Initialize the database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer database.DB.Close()

	fileServer := http.FileServer(http.Dir("./app"))
	http.Handle("/app/", http.StripPrefix("/app", fileServer))

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/session", handlers.SessionHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/profile", handlers.ProfileHandler)
	http.HandleFunc("/new_post", handlers.CreatePostHandler)
	http.HandleFunc("/post", handlers.PostHandler)
	// http.HandleFunc("/group_post", handlers.GroupPostHandler)
	http.HandleFunc("/comment", handlers.CreateCommentHandler)
	http.HandleFunc("/group_comment", handlers.CreateGroupCommentHandler)
	http.HandleFunc("/like", handlers.LikeHandler)
	http.HandleFunc("/like_group", handlers.LikeGroupHandler)
	http.HandleFunc("/saves", handlers.SaveHandler)
	http.HandleFunc("/get_saved_posts", handlers.GetSavedPostsHandler)
	http.HandleFunc("/follow", handlers.FollowHandler)
	// http.HandleFunc("/unfollow", handlers.UnfollowHandler)
	http.HandleFunc("/followers", handlers.FollowersHandler)
	http.HandleFunc("/following", handlers.FollowingHandler)
	http.HandleFunc("/new_group", handlers.CreateGrpoupHandler)
	http.HandleFunc("/new_post_group", handlers.CreatePostGroupHandler)
	http.HandleFunc("/new_event", handlers.CreateEventHandler)
	http.HandleFunc("/groups", handlers.GetGroupsHandler)
	http.HandleFunc("/group", handlers.GroupHandler)
	http.HandleFunc("/notifications", handlers.NotificationsHandler)
	http.HandleFunc("/notifications/mark_as_read", handlers.MarkNotificationAsReadHandler)
	http.HandleFunc("/notifications/mark_all_as_read", handlers.MarkAllNotificationsAsReadHandler)
	http.HandleFunc("/chats", handlers.ChatHandler)
	http.HandleFunc("/chats_group", handlers.ChatGroupHandler)
	http.HandleFunc("/message", handlers.SendMessageHandler)
	http.HandleFunc("/ws", handlers.WebSocketHandler)
	// http.HandleFunc("/search", handlers.SearchHandler)

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

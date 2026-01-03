package main

import (
	"log"
	"net/http"

	database "social-network/database"
	handlers "social-network/handlers"

	"github.com/rs/cors"
)

func main() {
	//  CORS  localhost:3000
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Cookie"}, 
		AllowCredentials: true,
	})

	// Initialize the database
	if err := database.InitializeDatabase(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer database.Database.Close()

	fileServer := http.FileServer(http.Dir("./app"))
	http.Handle("/app/", http.StripPrefix("/app", fileServer))

	http.HandleFunc("/", handlers.SessionHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/home", handlers.Home)
	http.HandleFunc("/user", handlers.CheckTheUserHandler)
	http.HandleFunc("/connections", handlers.GetConnectionsHandler)
	http.HandleFunc("/profile", handlers.ProfileHandler)
	http.HandleFunc("/profile_posts", handlers.ProfilePostsHandler)
	http.HandleFunc("/get_saved_posts", handlers.GetSavedPostsHandler)
	http.HandleFunc("/followers", handlers.GetFollowersHandler)
	http.HandleFunc("/following", handlers.GetFollowingHandler)
	http.HandleFunc("/new_post", handlers.CreatePostHandler)
	http.HandleFunc("/post", handlers.PostHandler)
	http.HandleFunc("/story", handlers.CreateStoryHandler)
	http.HandleFunc("/seen_story", handlers.SeenStory)
	http.HandleFunc("/categories", handlers.GetTopCategoriesHandler)
	http.HandleFunc("/posts_category", handlers.GetPostsByCategory)
	http.HandleFunc("/comment", handlers.CreateCommentHandler)
	http.HandleFunc("/like", handlers.LikeHandler)
	http.HandleFunc("/save", handlers.SaveHandler)
	http.HandleFunc("/follow", handlers.InvitationsHandler)
	http.HandleFunc("/suggested_users", handlers.GetSuggestedUsersHandler)
	http.HandleFunc("/new_group", handlers.CreateGrpoupHandler)
	http.HandleFunc("/group", handlers.GroupHandler)
	http.HandleFunc("/group_details", handlers.GroupDetailsHandler)
	http.HandleFunc("/new_post_group", handlers.CreatePostGroupHandler)
	http.HandleFunc("/new_event", handlers.CreateEventHandler)
	http.HandleFunc("/add_members", handlers.AddMembers)
	http.HandleFunc("/join", handlers.InvitationsHandler)
	http.HandleFunc("/groups", handlers.GetGroupsHandler)
	http.HandleFunc("/events", handlers.GetEventsHandler)
	http.HandleFunc("/event", handlers.GetEventHandler)
	http.HandleFunc("/join_to_event", handlers.JoinToEventHandler)
	http.HandleFunc("/accept_invitation", handlers.AcceptInvitationHandler)
	http.HandleFunc("/reject_invitation", handlers.DeclineInvitationHandler)
	http.HandleFunc("/accept_invitation_other", handlers.AcceptOtherInvitationHandler)
	http.HandleFunc("/invitations_groups", handlers.GetGroupInvitations)
	http.HandleFunc("/notifications", handlers.NotificationsHandler)
	http.HandleFunc("/mark_notifications_as_read", handlers.MarkNotificationsAsReadHandler)
	http.HandleFunc("/read_notification", handlers.MarkNotificationAsReadHandler)
	http.HandleFunc("/chats", handlers.ChatHandler)
	http.HandleFunc("/chats_group", handlers.ChatGroupHandler)
	http.HandleFunc("/read_messages", handlers.ReadMessagesHandler)
	http.HandleFunc("/search", handlers.SearchHandler)
	http.HandleFunc("/ws", handlers.WebSocketHandler)

	log.Println("Server started on :8404")
	err := http.ListenAndServe(":8404", c.Handler(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}

package structs

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

type Db struct {
	db *sql.DB
}

type User struct {
	ID                 int64     `json:"id" sqlite:"id"`
	Username           string    `json:"username" sqlite:"username"`
	FirstName          string    `json:"first_name" sqlite:"first_name"`
	LastName           string    `json:"last_name" sqlite:"last_name"`
	Email              string    `json:"email" sqlite:"email"`
	DateOfBirth        time.Time `json:"date_of_birth" sqlite:"date_of_birth"`
	Password           []byte    `json:"password" sqlite:"password"`
	ConfirmPass        []byte    `json:"confirm_pass" sqlite:"confirm_pass"`
	CreatedAt          time.Time `json:"created_at" sqlite:"created_at"`
	Avatar             string    `json:"avatar" sqlite:"avatar"`
	Cover              string    `json:"cover" sqlite:"cover"`
	Bio                string    `json:"bio" sqlite:"bio"`
	Privacy            string    `json:"privacy" sqlite:"privacy"`
	TotalPosts         int64     `json:"total_posts" sqlite:"total_posts"`
	TotalLikes         int64     `json:"total_likes" sqlite:"total_likes"`
	TotalComments      int64     `json:"total_comments" sqlite:"total_comments"`
	TotalFollowers     int64     `json:"total_followers" sqlite:"total_followers"`
	TotalFollowing     int64     `json:"total_following" sqlite:"total_following"`
	TotalGroups        int64     `json:"total_groups" sqlite:"total_groups"`
	TotalEvents        int64     `json:"total_events" sqlite:"total_events"`
	TotalNotifications int64     `json:"total_notifications" sqlite:"total_notifications"`
	TotalMessages      int64     `json:"total_messages" sqlite:"total_messages"`
	TotalInvitations   int64     `json:"total_invitations" sqlite:"total_invitations"`
	Online             bool      `json:"online" sqlite:"online"`
	SessionToken       uuid.UUID `json:"session_token" sqlite:"session_token"`
}

type Post struct {
	ID                 int64     `json:"id" sqlite:"id"`
	Title              string    `json:"title" sqlite:"title"`
	UserID             int64     `json:"user_id" sqlite:"user_id"`
	Avatar             string    `json:"avatar" sqlite:"avatar"`
	Content            string    `json:"content" sqlite:"content"`
	CategoryID         int64     `json:"category_id" sqlite:"category_id"`
	GroupName          string    `json:"group_name" sqlite:"group_name"`
	GroupID            int64     `json:"group_id" sqlite:"group_id"`
	Category           string    `json:"category" sqlite:"category"`
	CategoryColor      string    `json:"category_color" sqlite:"category_color"`
	CategoryBackground string    `json:"category_background" sqlite:"category_background"`
	Image              string    `json:"image" sqlite:"image"`
	Author             string    `json:"author" sqlite:"author"`
	CreatedAt          string    `json:"created_at" sqlite:"created_at"`
	TotalLikes         int64     `json:"total_likes" sqlite:"total_likes"`
	IsLiked            bool      `json:"is_liked" sqlite:"is_liked"`
	TotalComments      int64     `json:"total_comments" sqlite:"total_comments"`
	Comments           []Comment `json:"comments" sqlite:"comments"`
	Privacy            string    `json:"privacy" sqlite:"privacy"`
	Saved              bool      `json:"saved" sqlite:"saved"`
}

type Comment struct {
	ID        int64  `json:"id" sqlite:"id"`
	PostID    int64  `json:"post_id" sqlite:"post_id"`
	GroupID   int64  `json:"group_id" sqlite:"group_id"`
	Content   string `json:"content" sqlite:"content"`
	Author    string `json:"author" sqlite:"author"`
	CreatedAt string `json:"created_at" sqlite:"created_at"`
}

type Category struct {
	ID    int64  `json:"id" sqlite:"id"`
	Name  string `json:"name" sqlite:"name"`
	Count int64  `json:"count" sqlite:"count"`
}

type Group struct {
	ID           int64     `json:"id" sqlite:"id"`
	Name         string    `json:"name" sqlite:"name"`
	Image        string    `json:"image" sqlite:"image"`
	Cover        string    `json:"cover" sqlite:"cover"`
	Description  string    `json:"description" sqlite:"description"`
	CreatedAt    time.Time `json:"created_at" sqlite:"created_at"`
	Admin        string    `json:"admin" sqlite:"admin"`
	Privacy      string    `json:"privacy" sqlite:"privacy"`
	Members      []User    `json:"members" sqlite:"members"`
	TotalMembers int64     `json:"total_members" sqlite:"total_members"`
}

type Message struct {
	ID         int64  `json:"id" sqlite:"id"`
	SenderID   int64  `json:"sender_id" sqlite:"sender_id"`
	GroupID    int64  `json:"group_id" sqlite:"group_id"`
	ReceiverID int64  `json:"receiver_id" sqlite:"receiver_id"`
	Username   string `json:"username" sqlite:"username"`
	Avatar     string `json:"avatar" sqlite:"avatar"`
	Content    string `json:"content" sqlite:"content"`
	Image      string `json:"image" sqlite:"image"`
	CreatedAt  string `json:"created_at" sqlite:"created_at"`
	Type       string `json:"type" sqlite:"type"`
}

type Notification struct {
	ID               int64     `json:"id" sqlite:"id"`
	Username         string    `json:"username" sqlite:"username"`
	Avatar           string    `json:"avatar" sqlite:"avatar"`
	TypeNotification string    `json:"type_notification" sqlite:"type_notification"`
	CreatedAt        time.Time `json:"created_at" sqlite:"created_at"`
	Read             bool      `json:"read" sqlite:"read"`
}

type Invitation struct {
	ID       int64  `json:"id" sqlite:"id"`
	SenderID int64  `json:"sender_id" sqlite:"sender_id"`
	GroupID  int64  `json:"group_id" sqlite:"group_id"`
	Sender   string `json:"sender" sqlite:"sender"`
	Avatar   string `json:"avatar" sqlite:"avatar"`
	Group    string `json:"group" sqlite:"group"`
}

type Event struct {
	ID          int64     `json:"id" sqlite:"id"`
	Group       string    `json:"group" sqlite:"group"`
	GroupID     int64     `json:"group_id" sqlite:"group_id"`
	GroupName   string    `json:"group_name" sqlite:"group_name"`
	Name        string    `json:"name" sqlite:"name"`
	Description string    `json:"description" sqlite:"description"`
	Image       string    `json:"image" sqlite:"image"`
	Location    string    `json:"location" sqlite:"location"`
	StartDate   time.Time `json:"start_date" sqlite:"start_date"`
	EndDate     time.Time `json:"end_date" sqlite:"end_date"`
	CreatedAt   string    `json:"created_at" sqlite:"created_at"`
	Creator     string    `json:"creator" sqlite:"creator"`
	Attendees   []User    `json:"attendees" sqlite:"attendees"`
}

package structs

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID             int64
	Username       string
	FirstName      string
	LastName       string
	Email          string
	DateOfBirth    time.Time
	Password       []byte
	CreatedAt      time.Time
	ProfileImage   string
	Bio            string
	TotalPosts     int64
	TotalLikes     int64
	TotalComments  int64
	TotalFollowers int64
	TotalFollowing int64
	SessionToken   uuid.UUID
}

type Post struct {
	ID            int64
	Title         string
	Content       string
	Category      string
	Image         string
	Author        string
	CreatedAt     string
	TotalLikes    int64
	TotalComments int64
	Comments      []Comment
	Privacy       string
}

type Comment struct {
	ID        int64
	Content   string
	Author    string
	CreatedAt string
}

type Category struct {
	ID   int64
	Name string
}

type Group struct {
	ID          int64
	Name        string
	Image       string
	Description string
	CreatedAt   time.Time
	Admin       string
	Members     []User
}

type Message struct {
	ID        int64
	Sender    int64
	Recipient int64
	Content   string
	CreatedAt time.Time
	Read      bool
}

type Notification struct {
	ID               int64
	UserID           int64
	NotifiedUser     int64
	Content          string
	TypeNotification string
	CreatedAt        time.Time
	Read             bool
}

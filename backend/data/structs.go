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
	Author        string
	Categories    []string
	CreatedAt     string
	TotalLikes    int
	TotalComments int
	Comments      []Comment
}

type Comment struct {
	ID        int
	Content   string
	Author    string
	CreatedAt string
}

type Category struct {
	ID   int
	Name string
}

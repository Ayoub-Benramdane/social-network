package structs

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID           int64
	Username     string
	FirstName    string
	LastName     string
	Email        string
	DateOfBirth  time.Time
	Password     []byte
	SessionToken uuid.UUID
}

type Post struct {
	ID        int
	Title     string
	Content   string
	Author    string
	Category  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Comments  []Comment
}

type Comment struct {
	ID        int
	Content   string
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Category struct {
	ID   int
	Name string
}

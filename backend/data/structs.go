package structs

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

var ConnectedClients = make(map[int64][]*websocket.Conn)

type User struct {
	UserID                int64     `json:"user_id" sqlite:"user_id"`
	Username              string    `json:"username" sqlite:"username"`
	FirstName             string    `json:"first_name" sqlite:"first_name"`
	LastName              string    `json:"last_name" sqlite:"last_name"`
	Email                 string    `json:"email" sqlite:"email"`
	BirthDate             time.Time `json:"date_of_birth" sqlite:"date_of_birth"`
	Password              string    `json:"password" sqlite:"password"`
	PasswordConfirmation  string    `json:"confirm_pass" sqlite:"confirm_pass"`
	CreatedAt             time.Time `json:"created_at" sqlite:"created_at"`
	AvatarURL             string    `json:"avatar" sqlite:"avatar"`
	CoverURL              string    `json:"cover" sqlite:"cover"`
	Bio                   string    `json:"bio" sqlite:"bio"`
	Role                  string    `json:"role" sqlite:"role"`
	PrivacyLevel          string    `json:"privacy" sqlite:"privacy"`
	IsTyping              bool      `json:"is_typing" sqlite:"is_typing"`
	MessageCount          int64     `json:"total_messages" sqlite:"total_messages"`
	GroupMessageCount     int64     `json:"total_groups_messages" sqlite:"total_groups_messages"`
	ChatMessageCount      int64     `json:"total_chats_messages" sqlite:"total_chats_messages"`
	FollowerCount         int64     `json:"total_followers" sqlite:"total_followers"`
	FollowingCount        int64     `json:"total_following" sqlite:"total_following"`
	GroupCount            int64     `json:"total_groups" sqlite:"total_groups"`
	EventCount            int64     `json:"total_events" sqlite:"total_events"`
	PostCount             int64     `json:"total_posts" sqlite:"total_posts"`
	LikeCount             int64     `json:"total_likes" sqlite:"total_likes"`
	CommentCount          int64     `json:"total_comments" sqlite:"total_comments"`
	SaveCount             int64     `json:"total_saves" sqlite:"total_saves"`
	NotificationCount     int64     `json:"total_notifications" sqlite:"total_notifications"`
	InvitationCount       int64     `json:"total_invitations" sqlite:"total_invitations"`
	IsFollowing           bool      `json:"is_following" sqlite:"is_following"`
	IsFollower            bool      `json:"is_follower" sqlite:"is_follower"`
	IsPending             bool      `json:"is_pending" sqlite:"is_pending"`
	IsOnline              bool      `json:"online" sqlite:"online"`
	SessionID             uuid.UUID `json:"session_token" sqlite:"session_token"`
	AccountType           string    `json:"type" sqlite:"type"`
	UserStories           []Stories `json:"stories" sqlite:"stories"`
}

type Post struct {
	PostID              int64     `json:"post_id" sqlite:"post_id"`
	Title               string    `json:"title" sqlite:"title"`
	AuthorID            int64     `json:"user_id" sqlite:"user_id"`
	AuthorAvatar        string    `json:"avatar" sqlite:"avatar"`
	Content             string    `json:"content" sqlite:"content"`
	CategoryID          int64     `json:"category_id" sqlite:"category_id"`
	GroupName           string    `json:"group_name" sqlite:"group_name"`
	GroupID             int64     `json:"group_id" sqlite:"group_id"`
	CategoryName        string    `json:"category" sqlite:"category"`
	CategoryColor       string    `json:"category_color" sqlite:"category_color"`
	CategoryBackground  string    `json:"category_background" sqlite:"category_background"`
	ImageURL            string    `json:"image" sqlite:"image"`
	AuthorName          string    `json:"author" sqlite:"author"`
	CreatedAt           string    `json:"created_at" sqlite:"created_at"`
	IsLiked             bool      `json:"is_liked" sqlite:"is_liked"`
	LikedBy             []User    `json:"who_liked" sqlite:"who_liked"`
	LikeCount           int64     `json:"total_likes" sqlite:"total_likes"`
	CommentCount        int64     `json:"total_comments" sqlite:"total_comments"`
	SaveCount           int64     `json:"total_saves" sqlite:"total_saves"`
	PostCount           int64     `json:"total_posts" sqlite:"total_posts"`
	Comments             []Comment `json:"comments" sqlite:"comments"`
	PrivacyLevel        string    `json:"privacy" sqlite:"privacy"`
	IsSaved             bool      `json:"saved" sqlite:"saved"`
}

type Comment struct {
	CommentID  int64  `json:"comment_id" sqlite:"comment_id"`
	AuthorID   int64  `json:"user_id" sqlite:"user_id"`
	AvatarURL  string `json:"avatar" sqlite:"avatar"`
	Username   string `json:"username" sqlite:"username"`
	PostID     int64  `json:"post_id" sqlite:"post_id"`
	GroupID    int64  `json:"group_id" sqlite:"group_id"`
	Content    string `json:"content" sqlite:"content"`
	ImageURL   string `json:"image" sqlite:"image"`
	CreatedAt  string `json:"created_at" sqlite:"created_at"`
}

type Stories struct {
	HasUnseen bool    `json:"unseen_story" sqlite:"unseen_story"`
	User      User    `json:"user" sqlite:"user"`
	Items     []Story `json:"stories" sqlite:"stories"`
}

type Story struct {
	StoryID   int64     `json:"story_id" sqlite:"story_id"`
	ImageURL  string    `json:"image" sqlite:"image"`
	IsRead    bool      `json:"status" sqlite:"status"`
	CreatedAt time.Time `json:"created_at" sqlite:"created_at"`
}

type Category struct {
	CategoryID int64  `json:"category_id" sqlite:"category_id"`
	Name       string `json:"name" sqlite:"name"`
	Color      string `json:"color" sqlite:"color"`
	Background string `json:"background" sqlite:"background"`
	ItemCount  int64  `json:"count" sqlite:"count"`
}

type Group struct {
	GroupID        int64  `json:"group_id" sqlite:"group_id"`
	Name           string `json:"name" sqlite:"name"`
	ImageURL       string `json:"image" sqlite:"image"`
	CoverURL       string `json:"cover" sqlite:"cover"`
	Description    string `json:"description" sqlite:"description"`
	CreatedAt      string `json:"created_at" sqlite:"created_at"`
	AdminName      string `json:"admin" sqlite:"admin"`
	AdminID        int64  `json:"admin_id" sqlite:"admin_id"`
	IsOwner        bool   `json:"owner" sqlite:"owner"`
	InvitedBy      int64  `json:"invited_by" sqlite:"invited_by"`
	PrivacyLevel   string `json:"privacy" sqlite:"privacy"`
	UserRole       string `json:"role" sqlite:"role"`
	ActionType     string `json:"type" sqlite:"type"`
	MemberCount    int64  `json:"total_members" sqlite:"total_members"`
	PostCount      int64  `json:"total_posts" sqlite:"total_posts"`
	MessageCount   int64  `json:"total_messages" sqlite:"total_messages"`
}

type Message struct {
	MessageID            int64  `json:"message_id" sqlite:"message_id"`
	CurrentUserID        int64  `json:"current_user" sqlite:"current_user"`
	GroupID              int64  `json:"group_id" sqlite:"group_id"`
	SenderID             int64  `json:"user_id" sqlite:"user_id"`
	AvatarURL             string `json:"avatar" sqlite:"avatar"`
	Username              string `json:"username" sqlite:"username"`
	FirstName             string `json:"first_name" sqlite:"first_name"`
	LastName              string `json:"last_name" sqlite:"last_name"`
	SenderUsername        string `json:"sender_username" sqlite:"sender_username"`
	SenderAvatarURL       string `json:"sender_avatar" sqlite:"sender_avatar"`
	Content               string `json:"content" sqlite:"content"`
	ImageURL              string `json:"image" sqlite:"image"`
	CreatedAt             string `json:"created_at" sqlite:"created_at"`
	MessageCount          int64  `json:"total_messages" sqlite:"total_messages"`
	ChatMessageCount      int64  `json:"total_chat_messages" sqlite:"total_chat_messages"`
	GroupMessageCount     int64  `json:"total_group_messages" sqlite:"total_group_messages"`
	MessageType           string `json:"type" sqlite:"type"`
}

type Notification struct {
	NotificationID   int64  `json:"notification_id" sqlite:"notification_id"`
	UserID           int64  `json:"user_id" sqlite:"user_id"`
	Username         string `json:"username" sqlite:"username"`
	AvatarURL        string `json:"avatar" sqlite:"avatar"`
	Content          string `json:"content" sqlite:"content"`
	PostID           int64  `json:"post_id" sqlite:"post_id"`
	GroupID          int64  `json:"group_id" sqlite:"group_id"`
	EventID          int64  `json:"event_id" sqlite:"event_id"`
	NotificationType string `json:"type_notification" sqlite:"type_notification"`
	Message          string `json:"notification_message" sqlite:"notification_message"`
	CreatedAt        string `json:"created_at" sqlite:"created_at"`
	IsRead           bool   `json:"read" sqlite:"read"`
}

type Invitation struct {
	InvitationID int64  `json:"invitation_id" sqlite:"invitation_id"`
	CreatedAt    string `json:"created_at" sqlite:"created_at"`
	IsOwner      bool   `json:"owner" sqlite:"owner"`
	User         User   `json:"user" sqlite:"user"`
	Group        Group  `json:"group" sqlite:"group"`
}

type Event struct {
	EventID       int64     `json:"event_id" sqlite:"event_id"`
	CreatorID     int64     `json:"user_id" sqlite:"user_id"`
	Username      string    `json:"username" sqlite:"username"`
	AvatarURL     string    `json:"avatar" sqlite:"avatar"`
	GroupID       int64     `json:"group_id" sqlite:"group_id"`
	GroupName     string    `json:"group_name" sqlite:"group_name"`
	Name          string    `json:"name" sqlite:"name"`
	Description   string    `json:"description" sqlite:"description"`
	ImageURL      string    `json:"image" sqlite:"image"`
	EventType     string    `json:"type" sqlite:"type"`
	Location      string    `json:"location" sqlite:"location"`
	StartDate     time.Time `json:"start_date" sqlite:"start_date"`
	EndDate       time.Time `json:"end_date" sqlite:"end_date"`
	CreatedAt     string    `json:"created_at" sqlite:"created_at"`
	CreatorName   string    `json:"creator" sqlite:"creator"`
	MemberCount   int64     `json:"total_members" sqlite:"total_members"`
}

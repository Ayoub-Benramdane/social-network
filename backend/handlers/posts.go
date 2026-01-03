package handlers

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"

	structs "social-network/data"
	"social-network/database"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		fmt.Println("Failed to retrieve user", err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !CheckLastActionTime(w, r, "posts") {
		return
	}

	switch r.Method {
	case http.MethodGet:
		NewPostGet(w, r, user)
	case http.MethodPost:
		NewPostPost(w, r, user)
	default:
		fmt.Println("Method not allowed", r.Method)
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}
}

func NewPostGet(w http.ResponseWriter, r *http.Request, user *structs.User) {
	categories, err := database.FetchAllCategories()
	if err != nil {
		fmt.Println("Error retrieving categories:", err)
		response := map[string]string{"error": "Failed to retrieve categories"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	users, err := database.GetUserFollowers(user.UserID)
	if err != nil {
		log.Printf("Error retrieving users: %v", err)
		response := map[string]string{"error": "Failed to retrieve users"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	data := struct {
		Categories []structs.Category
		Users      []structs.User
	}{
		Categories: categories,
		Users:      users,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func NewPostPost(w http.ResponseWriter, r *http.Request, user *structs.User) {
	var post structs.Post
	var err error
	post.Title = strings.TrimSpace(r.FormValue("title"))
	post.Content = strings.TrimSpace(r.FormValue("content"))
	post.PrivacyLevel = strings.TrimSpace(r.FormValue("privacy"))
	post.CategoryID, err = strconv.ParseInt(r.FormValue("category"), 10, 64)
	if err != nil {
		fmt.Println("Invalid category", err)
		response := map[string]string{"error": "Invalid category"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	errors, valid := ValidatePost(post.Title, post.Content, post.PrivacyLevel)
	if !valid {
		fmt.Println("Validation error", errors)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "Validation error",
			"fields": errors,
		})
		return
	}

	var imagePath string
	image, header, err := r.FormFile("postImage")
	if err != nil && err.Error() != "http: no such file" {
		fmt.Println("Failed to retrieve image", err)
		response := map[string]string{"error": "Failed to retrieve image"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if image != nil {
		imagePath, err = SaveImage(image, header, "../frontend/public/images/")
		if err != nil {
			fmt.Println("Failed to save image", err)
			response := map[string]string{"error": err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		newpath := strings.Split(imagePath, "/public")
		imagePath = newpath[1]
	}

	users := []string{}
	if post.PrivacyLevel == "private" {
		users = strings.Split(r.FormValue("users"), ",")
		if len(users) == 1 && users[0] == "" {
			fmt.Println("No users provided for private post")
			response := map[string]string{"error": "No users provided for private post"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	ClientsMutex.Lock()
	id, err := database.CreatePost(user.UserID, post.GroupID, post.CategoryID, post.Title, post.Content, imagePath, post.PrivacyLevel)
	if err != nil {
		fmt.Println("Failed to create post", err)
		response := map[string]string{"error": "Failed to create post"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	ClientsMutex.Unlock()

	if post.PrivacyLevel == "private" {
		for _, usr := range users {
			usr_id, err := strconv.ParseInt(usr, 10, 64)
			if err != nil {
				fmt.Println("Invalid user", err)
				response := map[string]string{"error": "Invalid user"}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
				return
			}

			if _, err = database.FindUserByID(usr_id); err != nil {
				fmt.Println("Invalid user", err)
				response := map[string]string{"error": "Invalid user"}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
				return
			}

			if err = database.AddAlmostPrivateUser(usr_id, id); err != nil {
				fmt.Println("Failed to give authorizen to this user", err)
				response := map[string]string{"error": "Failed to give authorizen to this user"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
		}
	}

	category, err := database.FetchCategoryByID(post.CategoryID)
	if err != nil {
		fmt.Println("Failed to retrieve category", err)
		response := map[string]string{"error": "Failed to retrieve category"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	newPost := structs.Post{
		PostID:             id,
		AuthorID:           user.UserID,
		AuthorName:         user.Username,
		CategoryID:         post.CategoryID,
		CategoryName:       category.Name,
		CategoryColor:      category.Color,
		CategoryBackground: category.Background,
		Title:              html.EscapeString(post.Title),
		Content:            html.EscapeString(post.Content),
		PrivacyLevel:       post.PrivacyLevel,
		ImageURL:           imagePath,
		CreatedAt:          "Just Now",
		LikeCount:          0,
		CommentCount:       0,
		Comments:           nil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newPost)
}

func CreatePostGroupHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		fmt.Println("Failed to retrieve user", err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	var group_id int64
	if r.Method == http.MethodGet {
		group_id, err = strconv.ParseInt(r.URL.Query().Get("group_id"), 10, 64)
		if err != nil {
			fmt.Println("Invalid group ID", err)
			response := map[string]string{"error": "Invalid group ID"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
	} else if r.Method == http.MethodPost {
		group_id, err = strconv.ParseInt(r.FormValue("group_id"), 10, 64)
		if err != nil {
			fmt.Println("Invalid group ID", err)
			response := map[string]string{"error": "Invalid group ID"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		fmt.Println("Method not allowed", r.Method)
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	if _, err = database.GetGroupByID(group_id); err != nil {
		fmt.Println("Failed to retrieve group", err)
		response := map[string]string{"error": "Failed to retrieve group"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	member, err := database.IsUserGroupMember(user.UserID, group_id)
	if err != nil {
		fmt.Println("Failed to check if user is a member", err)
		response := map[string]string{"error": "Failed to check if user is a member"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	} else if !member {
		fmt.Println("You are not a member of this group")
		response := map[string]string{"error": "You are not a member of this group"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	switch r.Method {
	case http.MethodGet:
		NewPostGroupGet(w, r)
	case http.MethodPost:
		NewPostGroupPost(w, r, user, group_id)
	}
}

func NewPostGroupGet(w http.ResponseWriter, r *http.Request) {
	categories, err := database.FetchAllCategories()
	if err != nil {
		log.Printf("Error retrieving categories: %v", err)
		response := map[string]string{"error": "Failed to retrieve categories"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	data := struct {
		Categories []structs.Category
	}{
		Categories: categories,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func NewPostGroupPost(w http.ResponseWriter, r *http.Request, user *structs.User, group_id int64) {
	var post structs.Post
	var err error
	post.Title = strings.TrimSpace(r.FormValue("title"))
	post.Content = strings.TrimSpace(r.FormValue("content"))
	post.GroupID = group_id
	post.CategoryID, err = strconv.ParseInt(r.FormValue("category"), 10, 64)
	if err != nil {
		fmt.Println("Invalid category", err)
		response := map[string]string{"error": "Invalid category"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	group, err := database.GetGroupByID(post.GroupID)
	if err != nil {
		fmt.Println("Failed to retrieve group", err)
		response := map[string]string{"error": "Failed to retrieve groups"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	errors, valid := ValidatePost(post.Title, post.Content, group.PrivacyLevel)
	if !valid {
		fmt.Println("Validation error", errors)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "Validation error",
			"fields": errors,
		})
		return
	}

	var imagePath string
	image, header, err := r.FormFile("postImage")
	if err != nil && err.Error() != "http: no such file" {
		fmt.Println("Failed to retrieve image", err)
		response := map[string]string{"error": "Failed to retrieve image"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if image != nil {
		imagePath, err = SaveImage(image, header, "../frontend/public/images/")
		if err != nil {
			fmt.Println("Failed to save image", err)
			response := map[string]string{"error": err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		newpath := strings.Split(imagePath, "/public")
		imagePath = newpath[1]
	}

	id, err := database.CreatePost(user.UserID, post.GroupID, post.CategoryID, post.Title, post.Content, imagePath, "public")
	if err != nil {
		fmt.Println("Failed to create post", err)
		response := map[string]string{"error": "Failed to create post"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	category, err := database.FetchCategoryByID(post.CategoryID)
	if err != nil {
		fmt.Println("Failed to retrieve category", err)
		response := map[string]string{"error": "Failed to retrieve category"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	newPost := structs.Post{
		PostID:             id,
		AuthorID:           user.UserID,
		AuthorName:         user.Username,
		GroupName:          group.Name,
		GroupID:            group.GroupID,
		CategoryID:         post.CategoryID,
		CategoryName:       category.Name,
		CategoryColor:      category.Color,
		CategoryBackground: category.Background,
		Title:              html.EscapeString(post.Title),
		Content:            html.EscapeString(post.Content),
		ImageURL:           post.ImageURL,
		CreatedAt:          "Just Now",
		PrivacyLevel:       group.PrivacyLevel,
		LikeCount:          0,
		CommentCount:       0,
		Comments:           nil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newPost)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Method not allowed", r.Method)
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		fmt.Println("Failed to retrieve user", err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	post_id, err := strconv.ParseInt(r.URL.Query().Get("post_id"), 10, 64)
	if err != nil {
		fmt.Println("Invalid post ID", err)
		response := map[string]string{"error": "Invalid post ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	post, err := database.GetPost(user.UserID, post_id)
	if err != nil {
		fmt.Println("Failed to retrieve post", err)
		response := map[string]string{"error": "Failed to retrieve post"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if post.GroupID != 0 {
		group, err := database.GetGroupByID(post.GroupID)
		if err != nil {
			fmt.Println("Failed to retrieve group", err)
			response := map[string]string{"error": "Failed to retrieve groups"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		if group.PrivacyLevel == "private" {
			if member, err := database.IsUserGroupMember(user.UserID, post.GroupID); err != nil || !member {
				fmt.Println("Failed to check if user is member of group", err)
				response := map[string]string{"error": "Failed to check if user is member of group"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
		}
	} else if (post.PrivacyLevel == "almost_private" || post.PrivacyLevel == "private") && post.AuthorName != user.Username {
		if followed, err := database.IsUserFollowing(user.UserID, post.AuthorID); err != nil || !followed {
			fmt.Println("Failed to check if user is following author", err)
			response := map[string]string{"error": "You are not authorized to view this post"}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}
		if post.PrivacyLevel == "private" {
			if authorized, err := database.IsAuthorized(user.UserID, post_id); err != nil || !authorized {
				fmt.Println("Failed to check if user is authorized", err)
				response := map[string]string{"error": "You are not authorized to view this post"}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response)
				return
			}
		}
	}

	post.SaveCount, err = database.CountSaves(post_id, post.GroupID)
	if err != nil {
		fmt.Println("Failed to count saves", err)
		response := map[string]string{"error": "Failed to count saves"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	post.Comments, err = database.FetchPostComments(post_id)
	if err != nil {
		fmt.Println("Failed to retrieve comments", err)
		response := map[string]string{"error": "Failed to retrieve comments"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func GetPostsByCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Method not allowed", r.Method)
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		fmt.Println("Failed to retrieve user", err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	category_id, err := strconv.ParseInt(r.URL.Query().Get("category_id"), 10, 64)
	if err != nil {
		fmt.Println("Invalid category ID", err)
		response := map[string]string{"error": "Invalid category ID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	posts, err := database.GetPostsByCategory(category_id, user.UserID)
	if err != nil {
		fmt.Println("Failed to retrieve posts", err)
		response := map[string]string{"error": "Failed to retrieve posts"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	var posts_category []structs.Post
	for i := range posts {
		if posts[i].GroupID != 0 {
			group, err := database.GetGroupByID(posts[i].GroupID)
			if err != nil {
				fmt.Println("Failed to retrieve group", err)
				response := map[string]string{"error": "Failed to retrieve groups"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}

			if group.PrivacyLevel == "private" {
				if member, err := database.IsUserGroupMember(user.UserID, posts[i].GroupID); err != nil || !member {
					continue
				}
			}

		} else if (posts[i].PrivacyLevel == "almost_private" || posts[i].PrivacyLevel == "private") && posts[i].AuthorName != user.Username {
			if followed, err := database.IsUserFollowing(user.UserID, posts[i].AuthorID); err != nil || !followed {
				continue
			}
			if posts[i].PrivacyLevel == "private" {
				if authorized, err := database.IsAuthorized(user.UserID, posts[i].PostID); err != nil || !authorized {
					continue
				}
			}
		}
		posts[i].SaveCount, err = database.CountSaves(posts[i].PostID, posts[i].GroupID)
		if err != nil {
			fmt.Println("Failed to count saves", err)
			response := map[string]string{"error": "Failed to count saves"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		posts_category = append(posts_category, posts[i])
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts_category)
}

func ValidatePost(title, content, privacy string) (map[string]string, bool) {
	errors := make(map[string]string)
	const maxTitle = 100
	const maxContent = 2000

	if title == "" {
		errors["title"] = "Title is required"
	} else if len(title) > maxTitle {
		errors["title"] = "Title must be less than " + strconv.Itoa(maxTitle) + " characters"
	}

	if content == "" {
		errors["content"] = "Content is required"
	} else if len(content) > maxContent {
		errors["content"] = "Content must be less than " + strconv.Itoa(maxContent) + " characters"
	}

	if privacy == "" {
		errors["privacy"] = "PrivacyLevel is required"
	} else if privacy != "public" && privacy != "almost_private" && privacy != "private" {
		errors["privacy"] = "PrivacyLevel must be public, almost_private, or private"
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}

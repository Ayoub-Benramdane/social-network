package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	structs "social-network/data"
	"social-network/database"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed", r.Method)
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !CheckLastActionTime(w, r, "users") {
		return
	}

	var loginPayload structs.User
	err := json.NewDecoder(r.Body).Decode(&loginPayload)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		response := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	errors, valid := ValidateInput("", "", "", loginPayload.Email, loginPayload.Password, "", "", "", time.Now())
	if !valid {
		fmt.Println("Validation error:", errors)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "Validation error",
			"fields": errors,
		})
		return
	}

	user, err := database.FindUserByEmail(loginPayload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Invalid email", err)
			response := map[string]string{"email": "Invalid email"}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
		} else {
			log.Printf("Database error: %v", err)
			response := map[string]string{"error": "Internal server error"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginPayload.Password)); err != nil {
		fmt.Println("Password is incorrect")
		response := map[string]string{"password": "Password is incorrect"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	sessionToken, err := uuid.NewV4()
	if err != nil {
		log.Printf("Error generating session token: %v", err)
		response := map[string]string{"error": "Failed to generate session token"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := database.UpdateUserSession(loginPayload.Email, sessionToken); err != nil {
		log.Printf("Error updating session token: %v", err)
		response := map[string]string{"error": "Failed to update session token"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken.String(),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"username":     user.Username,
		"sessionToken": sessionToken.String(),
	})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed", r.Method)
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !CheckLastActionTime(w, r, "users") {
		return
	}

	var userPayload structs.User
	var err error
	userPayload.AccountType = strings.TrimSpace(r.FormValue("type"))
	userPayload.Username = strings.TrimSpace(r.FormValue("username"))
	userPayload.FirstName = strings.TrimSpace(r.FormValue("firstName"))
	userPayload.LastName = strings.TrimSpace(r.FormValue("lastName"))
	userPayload.Bio = strings.TrimSpace(r.FormValue("aboutMe"))
	userPayload.PrivacyLevel = strings.TrimSpace(r.FormValue("privacy"))

	if userPayload.AccountType == "register" {
		userPayload.Email = strings.TrimSpace(r.FormValue("email"))
		userPayload.Password = strings.TrimSpace(r.FormValue("password"))
		userPayload.PasswordConfirmation = strings.TrimSpace(r.FormValue("confirmedPassword"))
		temp := r.FormValue("dateOfBirth")
		userPayload.BirthDate, err = time.Parse("2006-01-02", temp)
		if err != nil {
			fmt.Println("Error parsing date of birth:", err)
			response := map[string]string{"error": "Error Parsing Date"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		errors, valid := ValidateInput(userPayload.Username, userPayload.FirstName, userPayload.LastName, userPayload.Email, userPayload.Password, userPayload.PasswordConfirmation, userPayload.PrivacyLevel, userPayload.Bio, userPayload.BirthDate)
		if !valid {
			fmt.Println("Validation error:", errors)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":  "Validation error",
				"fields": errors,
			})
			return
		}

		var avatarPath string
		avatarFile, avatarHeader, err := r.FormFile("avatar")
		if err != nil && err.Error() != "http: no such file" {
			fmt.Println("Error retrieving image:", err)
			response := map[string]string{"error": "Failed to retrieve image"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		} else if avatarFile != nil {
			avatarPath, err = SaveImage(avatarFile, avatarHeader, "../frontend/public/avatars/")
			if err != nil {
				fmt.Println("Error saving image:", err)
				response := map[string]string{"error": err.Error()}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
			newpath := strings.Split(avatarPath, "/public")
			avatarPath = newpath[1]
		} else {
			avatarPath = "/inconnu/avatar.png"
		}

		var coverPath string
		coverFile, coverHeader, err := r.FormFile("cover")
		if err != nil && err.Error() != "http: no such file" {
			fmt.Println("Error retrieving cover:", err)
			response := map[string]string{"error": "Failed to retrieve cover"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		} else if coverFile != nil {
			coverPath, err = SaveImage(coverFile, coverHeader, "../frontend/public/covers/")
			if err != nil {
				fmt.Println("Error saving cover:", err)
				response := map[string]string{"error": err.Error()}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
			newpath := strings.Split(coverPath, "/public")
			coverPath = newpath[1]
		} else {
			coverPath = "/inconnu/cover.jpg"
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPayload.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("Error hashing password:", err)
			response := map[string]string{"password": "Error hashing password"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		sessionToken, err := uuid.NewV4()
		if err != nil {
			log.Printf("Error generating session token: %v", err)
			response := map[string]string{"error": "Failed to generate session token"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		if err := database.CreateUser(userPayload.Username, userPayload.FirstName, userPayload.LastName, userPayload.Email, userPayload.Bio, avatarPath, coverPath, userPayload.PrivacyLevel, hashedPassword, userPayload.BirthDate, sessionToken); err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
				fmt.Println("Email alreadU(userPayload.Uy exists")
				response := map[string]string{"email": "Email already exists"}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
			} else if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
				fmt.Println("Username already exists")
				response := map[string]string{"username": "Username already exists"}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
			} else {
				log.Printf("Error inserting user: %v", err)
				response := map[string]string{"error": "Registration failed"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
			}
			return
		}
	} else if userPayload.AccountType == "update" {
		userPayload.UserID, err = strconv.ParseInt(r.FormValue("id"), 10, 64)
		if err != nil {
			fmt.Println("Error parsing ID:", err)
			response := map[string]string{"error": "Error Parsing ID"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		if err := database.UpdateProfile(userPayload.UserID, userPayload.Username, userPayload.FirstName, userPayload.LastName, userPayload.Bio, userPayload.PrivacyLevel); err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
				fmt.Println("Email already exists")
				response := map[string]string{"email": "Email already exists"}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
			} else if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
				fmt.Println("Username already exists")
				response := map[string]string{"usename": "Username already exists"}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
			} else {
				log.Printf("Error updating user: %v", err)
				response := map[string]string{"error": "Profile update failed"}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
			}
			return
		}

		if err := database.AcceptAllFriendInvitations(userPayload.UserID); err != nil {
			log.Printf("Error accepting invitations: %v", err)
			response := map[string]string{"error": "Failed to accept invitations"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := map[string]interface{}{
			"user":    userPayload,
			"message": "Profile updated successfully!",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	} else {
		fmt.Println("Invalid request method")
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]string{"message": "Registration successful! Please log in."}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Invalid request method", r.Method)
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	currentUser, err := GetUserFromSession(r)
	if err != nil || currentUser == nil {
		fmt.Println("Error retrieving user:", err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	ClientsMutex.Lock()
	if err := database.ClearUserSession(currentUser.UserID); err != nil {
		log.Printf("Error deleting session: %v", err)
		response := map[string]string{"error": "Failed to delete session"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	ClientsMutex.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "guest",
		MaxAge: -1,
	})

	response := map[string]string{"message": "Logout successful!", "username": currentUser.Username}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ValidateInput(username, firstName, lastName, email, password, confirmPassword, privacyLevel, aboutMe string, birthDate time.Time) (map[string]string, bool) {

	errors := make(map[string]string)

	const maxUsernameLength = 10
	const maxEmailLength = 30
	const maxPasswordLength = 20
	const maxNameLength = 20
	const maxAboutMeLength = 100

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if len(email) == 0 {
		errors["email"] = "Email cannot be empty"
	} else if len(email) > maxEmailLength {
		errors["email"] = fmt.Sprintf("Email cannot be longer than %d characters.", maxEmailLength)
	} else if !emailRegex.MatchString(email) {
		errors["email"] = "Invalid email format"
	}

	if password != confirmPassword && username != "" {
		errors["password"] = "Passwords do not match"
	} else if len(password) < 8 {
		errors["password"] = "Password must be at least 8 characters long"
	} else if len(password) > maxPasswordLength {
		errors["password"] = fmt.Sprintf("Password cannot be longer than %d characters.", maxPasswordLength)
	} else {
		hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecialChar := regexp.MustCompile(`[\W_]`).MatchString(password)

		if !hasUppercase {
			errors["password"] = "Password must include at least one uppercase letter"
		} else if !hasLowercase {
			errors["password"] = "Password must include at least one lowercase letter"
		} else if !hasDigit {
			errors["password"] = "Password must include at least one digit"
		} else if !hasSpecialChar {
			errors["password"] = "Password must include at least one special character"
		}
	}

	if username != "" {

		if len(firstName) == 0 {
			errors["first_name"] = "First name cannot be empty"
		} else if len(firstName) > maxNameLength {
			errors["first_name"] = fmt.Sprintf("First name cannot be longer than %d characters.", maxNameLength)
		} else if !isAlphabetic(firstName) {
			errors["first_name"] = "First name must contain only letters"
		}

		if len(lastName) == 0 {
			errors["last_name"] = "Last name cannot be empty"
		} else if len(lastName) > maxNameLength {
			errors["last_name"] = fmt.Sprintf("Last name cannot be longer than %d characters.", maxNameLength)
		} else if !isAlphabetic(lastName) {
			errors["last_name"] = "Last name must contain only letters"
		}

		if len(username) == 0 {
			errors["username"] = "Username cannot be empty"
		} else if len(username) > maxUsernameLength {
			errors["username"] = fmt.Sprintf("Username cannot be longer than %d characters.", maxUsernameLength)
		}

		if birthDate.IsZero() {
			errors["dateOfBirth"] = "Date cannot be empty"
		} else {
			age := time.Since(birthDate).Hours() / 24 / 365.3
			if age < 18 {
				errors["dateOfBirth"] = "You must be at least 18 years old"
			}

			minAllowedDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
			if birthDate.Before(minAllowedDate) {
				errors["dateOfBirth"] = "Date cannot be before the year 1900"
			}
		}

		if len(aboutMe) > maxAboutMeLength {
			errors["about_me"] = fmt.Sprintf("About me cannot be longer than %d characters.", maxAboutMeLength)
		}

		if privacyLevel != "public" && privacyLevel != "private" {
			errors["privacy"] = "Privacy must be either 'public' or 'private'"
		}
	}

	if len(errors) > 0 {
		log.Println(errors)
		return errors, false
	}

	return nil, true
}

func isAlphabetic(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) && char != ' ' {
			return false
		}
	}
	return true
}

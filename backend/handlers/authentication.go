package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"

	"social-network/database"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	var login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		response := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	errors, valid := ValidateInput("", "", "", login.Email, login.Password, "", "", "", time.Now())
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "Validation error",
			"fields": errors,
		})
		return
	}

	user, err := database.GetUserByEmail(login.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("1")

			fmt.Println(err)
			response := map[string]string{"error": "Invalid email or password"}
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		fmt.Println("2")

		fmt.Println(err)
		response := map[string]string{"error": "Password is incorrect"}
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

	if err := database.UpdateSession(login.Email, sessionToken); err != nil {
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

	data := map[string]interface{}{
		"username":     user.Username,
		"sessionToken": sessionToken.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	var register struct {
		Username          string    `json:"username"`
		FirstName         string    `json:"firstName"`
		LastName          string    `json:"lastName"`
		Email             string    `json:"email"`
		DateOfBirth       time.Time `json:"dateOfBirth"`
		Password          string    `json:"password"`
		ConfirmedPassword string    `json:"confirmedPassword"`
		AboutMe           string    `json:"aboutMe"`
		Privacy           string    `json:"privacy"`
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Cannot parse form", http.StatusBadRequest)
		return
	}

	register.Username = r.FormValue("username")
	register.FirstName = r.FormValue("firstName")
	register.LastName = r.FormValue("lastName")
	register.Email = r.FormValue("email")
	register.Password = r.FormValue("password")
	register.ConfirmedPassword = r.FormValue("confirmedPassword")
	register.AboutMe = r.FormValue("aboutMe")
	register.Privacy = r.FormValue("privacy")

	temp := r.FormValue("dateOfBirth")
	register.DateOfBirth, err = time.Parse("2006-01-02", temp)
	if err != nil {
		fmt.Println("here")
		fmt.Println(err)
		response := map[string]string{"error": "Error Parsing Date"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	register.Privacy = "public"

	errors, valid := ValidateInput(register.Username, register.FirstName, register.LastName, register.Email, register.Password, register.ConfirmedPassword, register.Privacy, register.AboutMe, register.DateOfBirth)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "Validation error",
			"fields": errors,
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("1")
		fmt.Println(err)
		response := map[string]string{"error": "Error hashing password"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	sessionToken, err := uuid.NewV4()
	if err != nil {
		fmt.Println("2")
		fmt.Println(err)
		log.Printf("Error generating session token: %v", err)
		response := map[string]string{"error": "Failed to generate session token"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	var imagePath string
	image, header, err := r.FormFile("avatar")
	if err != nil && err.Error() != "http: no such file" {
		fmt.Println("3")
		fmt.Println(err)
		response := map[string]string{"error": "Failed to retrieve image"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if image != nil {
		imagePath, err = SaveImage(image, header, "../frontend/public/avatars/")
		if err != nil {
			fmt.Println("4")
			fmt.Println(err)
			response := map[string]string{"error": err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		newpath := strings.Split(imagePath, "/public")
		imagePath = newpath[1]
	} else {
		imagePath = "/inconnu/avatar.png"
	}

	if err := database.RegisterUser(register.Username, register.FirstName, register.LastName, register.Email, register.AboutMe, imagePath, register.Privacy, hashedPassword, register.DateOfBirth, sessionToken); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			response := map[string]string{"error": "Email already exists"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
		} else if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
			response := map[string]string{"error": "Username already exists"}
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

	response := map[string]string{"message": "Registration successful! Please log in."}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := GetUserFromSession(r)
	if err != nil || user == nil {
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := database.DeleteSession(user.ID); err != nil {
		log.Printf("Error deleting session: %v", err)
		response := map[string]string{"error": "Failed to delete session"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "guest",
		MaxAge: -1,
	})

	response := map[string]string{"message": "Logout successful!", "username": user.Username}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ValidateInput(username, firstName, lastName, email, password, confirm_pass, privacy, aboutMe string, date time.Time) (map[string]string, bool) {
	errors := make(map[string]string)
	const maxUsername = 10
	const maxEmail = 30
	const maxPassword = 20
	const maxNameLength = 20
	const maxAboutMe = 100

	// ✅ Validation de l'email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if len(email) == 0 {
		errors["email"] = "Email cannot be empty"
	} else if len(email) > maxEmail {
		errors["email"] = fmt.Sprintf("Email cannot be longer than %d characters.", maxEmail)
	} else if !emailRegex.MatchString(email) {
		errors["email"] = "Invalid email format"
	}

	// ✅ Validation du mot de passe
	if password != confirm_pass && username != "" {
		errors["password"] = "Passwords do not match"
	} else if len(password) < 8 {
		errors["password"] = "Password must be at least 8 characters long"
	} else if len(password) > maxPassword {
		errors["password"] = fmt.Sprintf("Password cannot be longer than %d characters.", maxPassword)
	} else {
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial := regexp.MustCompile(`[\W_]`).MatchString(password)

		if !hasUpper {
			errors["password"] = "Password must include at least one uppercase letter"
		} else if !hasLower {
			errors["password"] = "Password must include at least one lowercase letter"
		} else if !hasDigit {
			errors["password"] = "Password must include at least one digit"
		} else if !hasSpecial {
			errors["password"] = "Password must include at least one special character"
		}
	}

	if username != "" {
		// ✅ Validation du prénom
		if len(firstName) == 0 {
			errors["first_name"] = "First name cannot be empty"
		} else if len(firstName) > maxNameLength {
			errors["first_name"] = fmt.Sprintf("First name cannot be longer than %d characters.", maxNameLength)
		} else if !isAlphabetic(firstName) {
			errors["first_name"] = "First name must contain only letters"
		}

		// ✅ Validation du nom
		if len(lastName) == 0 {
			errors["last_name"] = "Last name cannot be empty"
		} else if len(lastName) > maxNameLength {
			errors["last_name"] = fmt.Sprintf("Last name cannot be longer than %d characters.", maxNameLength)
		} else if !isAlphabetic(lastName) {
			errors["last_name"] = "Last name must contain only letters"
		}

		// ✅ Validation du username
		if len(username) == 0 {
			errors["username"] = "Username cannot be empty"
		} else if len(username) > maxUsername {
			errors["username"] = fmt.Sprintf("Username cannot be longer than %d characters.", maxUsername)
		}

		// ✅ Validation de la date
		if date.IsZero() {
			errors["date"] = "Date cannot be empty"
		} else {
			now := time.Now()
			year1900 := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

			if date.After(now) {
				errors["date"] = "Date cannot be in the future"
			} else if date.Before(year1900) {
				errors["date"] = "Date cannot be before the year 1900"
			}
		}
		// ✅ Validation de la description
		if len(aboutMe) == 0 {
			errors["about_me"] = "About Me cannot be empty"
		} else if len(aboutMe) > maxAboutMe {
			errors["about_me"] = fmt.Sprintf("About me cannot be longer than %d characters.", maxAboutMe)
		}

		// ✅ Validation de la privacy
		if privacy != "public" && privacy != "private" {
			errors["privacy"] = "Privacy must be either 'public' or 'private'"
		}
	}

	// Retour des erreurs
	if len(errors) > 0 {
		log.Println(errors)
		return errors, false
	}
	return nil, true
}

// 🔍 Fonction pour vérifier si un string contient uniquement des lettres
func isAlphabetic(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) && char != ' ' {
			return false
		}
	}
	return true
}

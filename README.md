# Social Network

A Facebook-like social network built as a full‑stack web application. This project was developed according to the required specifications and focuses on authentication, social interactions, real‑time communication, and proper system architecture.

The goal of this project is to demonstrate a complete understanding of frontend, backend, database design, real‑time systems, and containerization.

---

## 🚀 Project Overview

This social network allows users to interact with each other through posts, follows, groups, chats, and notifications, similar to common social media platforms.

Main implemented features:

* Followers system
* User profiles
* Posts and comments
* Groups and events
* Notifications
* Real‑time chats

---

## ✨ Features

### 👤 Authentication

* User registration and login
* Session and cookie‑based authentication
* Persistent login until logout

Registration includes:

* Email
* Password (encrypted)
* First name & last name
* Date of birth
* Optional avatar
* Optional nickname
* Optional "About me"

---

### 👥 Followers

* Send follow requests
* Accept or decline follow requests
* Automatic follow for public profiles
* Unfollow users

---

### 🙍 Profile

* Public and private profiles
* Profile information display
* User posts history
* Followers & following lists
* Toggle profile privacy

---

### 📝 Posts

* Create posts with text and optional image/GIF
* Comment on posts
* Like and save posts
* Post privacy levels:

  * Public
  * Almost private (followers only)
  * Private (selected followers)

---

### 📸 Stories (Extra Feature)

* Create temporary stories with images
* Stories visible for a limited time
* Seen / unseen story status
* Stories visible based on follow & privacy rules
* Real‑time story updates

---

### 👥 Groups

* Create groups with title and description
* Public and private groups
* Invite users to groups
* Request to join groups
* Accept or refuse group requests
* Group posts and comments (visible only to members)

---

### 📅 Events

* Create events inside groups
* Event details:

  * Title
  * Description
  * Date & time
* Event participation:

  * Going
  * Not going

---

### 💬 Chat

* Private messages between users
* Real‑time messaging using WebSockets
* Typing indicators
* Online / offline status
* Group chat rooms
* Emoji support

---

### 🔔 Notifications

* Real‑time notifications visible on all pages
* Notifications for:

  * Follow requests
  * Group invitations
  * Group join requests
  * New group events
  * Other user interactions

---

## 🛠 Tech Stack

### Frontend

* JavaScript
* Next.js (App Router)
* CSS Modules
* Responsive design

### Backend

* Go (Golang)
* net/http
* Gorilla WebSocket
* Session & cookie authentication

### Database

* SQLite
* Relational schema
* Foreign keys & constraints
* Migration system

### Migrations

* SQL migration files
* Automatic table creation on startup
* Structured migration folders

### Docker

* Separate Docker image for backend
* Separate Docker image for frontend
* docker‑compose for orchestration

---

## 📂 Project Structure

```
social-network/
│
├── backend/
│   ├── handlers/
│   ├── database/
│   ├── migrations/
│   ├── structs/
│   └── main.go
│
├── frontend/
│   ├── app/
│   ├── public/
│   ├── styles/
│   └── next.config.js
│
├── docker-compose.yml
└── README.md
```

---

## 🧪 Database

* SQLite database
* Managed using migrations
* Pre‑filled with logical test data for demonstration

---

## 🧠 Learning Outcomes

This project demonstrates knowledge of:

* Session‑based authentication
* Cookies handling
* SQL & database migrations
* WebSocket real‑time communication
* Docker containerization
* Full‑stack architecture

---

## ▶️ Running the Project

```bash
docker-compose up --build
```

* Frontend: [http://localhost:3000](http://localhost:3000)
* Backend: [http://localhost:8404](http://localhost:8404)

---

## 👨‍💻 Author

Developed as a complete full‑stack project, focusing on clean architecture, real‑time features, and real‑world social network logic.

---

## 📌 Notes

This project was built for educational purposes and follows the provided specifications closely, while allowing room for additional improvements and features.

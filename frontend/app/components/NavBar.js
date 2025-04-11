"use client";
import { useState } from "react";
import "../styles/NavBar.css";

export default function Navbar({ user }) {
  const handleLogout = async () => {
    try {
      const response = await fetch("http://localhost:8404/logout", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
      });

      if (response.ok) {
        setIsLoggedIn(false);
        setUser(null);
      }
    } catch (error) {
      console.error("Error logging out:", error);
    }
  };
  return (
    <div className="navbar">
      <div className="navbar-left">
        <div className="logo">
          <img src="./icons/logo.svg" alt="Logo" width={24} height={24} />
        </div>

        <div className="nav-links">
          <div className="home-link">
            <img src="./icons/home.svg" alt="Home" width={20} height={20} />
            <span>Home</span>
          </div>
          <div className="icon-link">
            <img
              src="./icons/message.svg"
              alt="Messages"
              width={20}
              height={20}
            />
          </div>
          <div className="icon-link">
            <img
              src="./icons/notification.svg"
              alt="Notifications"
              width={20}
              height={20}
            />
          </div>
        </div>
      </div>

      <div className="search-bar">
        <div className="search-input-container">
          <input type="text" placeholder="Search users and groups" />
          <img
            src="./icons/search.svg"
            alt="Search"
            width={16}
            height={16}
            className="search-icon"
          />
        </div>
      </div>

      <div className="user-actions">
        <div className="notification-badge">
          <div className="badge">1</div>
          <img
            className="notification-icon"
            src="./icons/notification.svg"
            alt="Notifications"
            width={20}
            height={20}
          />
        </div>
        <div className="message-badge">
          <div className="badge">3</div>
          <img
            src="./icons/message.svg"
            alt="Messages"
            width={20}
            height={20}
          />
        </div>
        <div className="user-profile">
          {/* <div className="user-name-avatar"> */}
            <img
              src={user.avatar}
              alt={user.username || "User"}
              className="user-avatar"
            />
            <span>{`${user.first_name} ${user.last_name}`}</span>
          {/* </div> */}

          <button onClick={handleLogout} className="logoutBtn">
            {/* <img
              src={"./icons/logout.svg"}
              alt="logout"
              className="logout-icon"
            /> */}
            <span>Logout</span>
          </button>
        </div>
      </div>
    </div>
  );
}

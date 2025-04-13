"use client";
import { useState, useEffect, useRef, use } from "react";
import "../styles/ChatWidget.css";
import UserCard from "./UserCard";

export default function ChatWidget({ users, groups }) {
  const [activeTab, setActiveTab] = useState("friends");
  const [selectedUser, setSelectedUser] = useState(null);

  const listToRender = activeTab === "friends" ? users : groups;

  async function showUserTab(user) {
    console.log(user);
    // setSelectedUser(user);

    // setActiveTab("friends");
    const id = user.id;
    // console.log(id);

    try {
      const response = await fetch("http://localhost:8404/chats", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(user.id),
      });
      // console.log(user.id);

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data);
      }
      const data = await response.json();

      console.log(data);
    } catch (error) {
      console.log(error);
    }
  }

  return (
    <div className="overall-chat-container">
      <div className="chat-header">
        <h4 className="chat-title">Messages</h4>
      </div>

      <div className="chat-tabs">
        <div
          className={`users-chat-tab ${
            activeTab === "friends" ? "active-tab" : ""
          }`}
          onClick={() => setActiveTab("friends")}
        >
          <h4>Friends</h4>
        </div>
        <div
          className={`users-chat-tab ${
            activeTab === "groups" ? "active-tab" : ""
          }`}
          onClick={() => setActiveTab("groups")}
        >
          <h4>Groups</h4>
        </div>
      </div>

      <div className="chat-container">
        <ul className="chat-content">
          {listToRender.map((user) => (
            <UserCard
              key={user.id}
              user={user}
              onClick={() => showUserTab(user)}
            />
          ))}
        </ul>
      </div>

      {/* {selectedUser && (
        <div className="selected-user-chat">
          <h4>
            Chat with {selectedUser.first_name} {selectedUser.last_name}
          </h4>
        </div>
      )} */}
    </div>
  );
}

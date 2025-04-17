"use client";
import { useState, useEffect, useRef, use } from "react";
import "../styles/ChatWidget.css";
import UserCard from "./UserCard";
import Message from "./Message";

export default function ChatWidget({ users, groups, myData }) {
  const [activeTab, setActiveTab] = useState("friends");
  // const [userData, setUserData] = useState();
  // const [myMessages, setMyMessages] = useState([]);
  // const [receivedMessages, setReceivedMessages] = useState([]);
  const [selectedUser, setSelectedUser] = useState(null);
  const [messages, setMessages] = useState([]);
  const [openWidget, setOpenWidget] = useState(true);
  const [openChatWidget, setOpenChatWidget] = useState(true);
  const [messageSending, setMessageSending] = useState("");

  const listToRender = activeTab === "friends" ? users : groups;

  async function handleMessagesSend(id) {
    const formData = new FormData();
    formData.append("receiver_id", id);
    formData.append("content", messageSending);
    // console.log(messageSending);
    try {
      const response = await fetch("http://localhost:8404/message", {
        method: "POST",
        credentials: "include",
        body: formData,
      });
      // console.log(id);
      const data = await response.json();
      if (!response.ok) {
        console.log(data);
      }
      console.log(data);
      setMessageSending("");
      setMessages((prevMessages) => [
        ...prevMessages,
        {
          content: messageSending,
          username: myData.username,
        },
      ]);
    } catch (error) {
      console.log(error);
    }
    console.log("Message:", messageSending);

    const socket = new WebSocket("ws://localhost:8404/ws");
    socket.onopen = () => {
      console.log("WebSocket Connected");
    };
    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data) {
        console.log("Data: ", data);
      }
    };
  }

  async function handleMessagesSending(id) {
    const formData = new FormData();
    formData.append("receiver_id", id);
    formData.append("content", messageSending);

    // console.log(messageSending);
    try {
      const response = await fetch("http://localhost:8404/message", {
        method: "POST",
        credentials: "include",
        body: formData,
      });
      console.log(id);

      const data = await response.json();
      if (!response.ok) {
        console.log(data);
      }
      console.log(data);
    } catch (error) {
      console.log(error);
    }
  }

  async function showUserTab(user) {
    setSelectedUser(user);
    console.log(user.id);

    try {
      const response = await fetch(
        `http://localhost:8404/chats?id=${user.id}`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
          // body: JSON.stringify(user.id),
        }
      );

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Failed to fetch messages");
      }

      const data = await response.json();
      console.log("User Data: ", data);

      setMessages(data);
    } catch (error) {
      console.error("Error fetching messages:", error);
    }
  }
  const toggleWidget = () => {
    setOpenWidget(!openWidget);
  };
  const toggleChatWidget = () => {
    setOpenChatWidget(!openChatWidget);
  };

  return (
    <div className="chat-wrapper-fixed">
      {selectedUser && (
        <div className="chat-box">
          <div
            className={`chat-header ${openChatWidget ? "opened" : "closed"}`}
            onClick={toggleChatWidget}
          >
            <img
              src={
                selectedUser.avatar
                  ? `${selectedUser.avatar}`
                  : `${selectedUser.image}`
              }
              className="chat-header-avatar"
              alt={selectedUser.username}
            />
            <h4 className="chat-title">{selectedUser.username}</h4>
            {/* {openChatWidget && (
              <span className="close-tab" onClick={setSelectedUser(false)}>
                X
              </span>
            )} */}
          </div>

          {openChatWidget && messages && (
            <div className="chat-messages">
              <div className="messages-container">
                {messages.map((msg) => (
                  <Message
                    key={msg.id}
                    message={msg}
                    isSent={msg.username !== selectedUser.username}
                  />
                ))}
              </div>
            </div>
          )}
          {openChatWidget && !messages && (
            <div className="chat-messages">
              <div className="messages-container">
                <h4>No Messages yet.</h4>
                {/* {messages.map((msg) => (
                <Message
                  key={msg.id}
                  message={msg}
                  isSent={msg.username !== selectedUser.username}
                />
              ))} */}
              </div>
            </div>
          )}
          {openChatWidget && (
            <div className="message-input">
              <input
                onChange={(e) => {
                  setMessageSending(e.target.value);
                }}
                className="message-input-input"
                placeholder="your message..."
              ></input>
              <div
                className="send-message-container"
                onClick={(e) => {
                  e.preventDefault();
                  handleMessagesSend(selectedUser.id);

                  // handleMessagesSending(selectedUser.id);
                }}
              >
                <img className="send-message-icon" src="./icons/send.svg"></img>
                <p>Send</p>
              </div>
            </div>
          )}
        </div>
      )}

      {!openWidget && (
        <div className="overall-chat-container-closed">
          <div className="chats-header" onClick={toggleWidget}>
            <h4 className="chat-title">Messages</h4>
            <div className="unread-messages">
              <p className="unread-message-number">{myData.total_messages}</p>
            </div>
          </div>
        </div>
      )}
      {openWidget && (
        <div className="overall-chat-container">
          <div className="chats-header" onClick={toggleWidget}>
            <h4 className="chat-title">Messages</h4>
            <div className="unread-messages">
              <p className="unread-message-number">
                {myData.total_messages} New Messages
              </p>
            </div>
          </div>

          <div className="chat-tabs">
            <div
              className={`users-chat-tab ${
                activeTab === "friends" ? "active-tab" : ""
              }`}
              onClick={() => setActiveTab("friends")}
            >
              <h4 className="friends-message-labes">Friends (5)</h4>
            </div>
            <div
              className={`users-chat-tab ${
                activeTab === "groups" ? "active-tab" : ""
              }`}
              onClick={() => setActiveTab("groups")}
            >
              <h4 className="groups-message-labes">Groups (2)</h4>
            </div>
          </div>

          <div className="chat-container">
            <ul className="chat-content">
              {listToRender.map((user) => (
                <UserCard
                  key={user.id}
                  user={user}
                  onClick={() => {
                    showUserTab(user);
                    setSelectedUser(user);
                  }}
                />
              ))}
            </ul>
          </div>
        </div>
      )}
    </div>
  );
}

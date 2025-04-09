"use client";
import "../styles/LeftSideBar.css";

export default function LeftSidebar() {
  const onlineUsers = [
    {
      id: 1,
      name: "Abdelkhalek Laidi",
      avatar: "https://via.placeholder.com/40",
    },
    {
      id: 2,
      name: "Ayoub Benramdan",
      avatar: "https://via.placeholder.com/40",
    },
    { id: 3, name: "Badr Lakraid", avatar: "https://via.placeholder.com/40" },
    {
      id: 4,
      name: "Abdelkhalek Laidi",
      avatar: "https://via.placeholder.com/40",
    },
    {
      id: 5,
      name: "Abdelkhalek Laidi",
      avatar: "https://via.placeholder.com/40",
    },
    {
      id: 6,
      name: "Abdelkhalek Laidi",
      avatar: "https://via.placeholder.com/40",
    },
  ];

  const communities = [
    { id: 1, name: "Volkswagen", avatar: "https://via.placeholder.com/40" },
    { id: 2, name: "Coffee", avatar: "https://via.placeholder.com/40" },
    { id: 3, name: "Zone01Oujda", avatar: "https://via.placeholder.com/40" },
  ];

  return (
    <div className="left-sidebar">
      <div className="search-box">
        <div className="search-input-container">
          <input type="text" placeholder="Search users and groups" />
          <svg
            className="search-icon"
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="#9CA3AF"
          >
            <path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z" />
          </svg>
        </div>
      </div>

      <div className="online-users">
        <h3>Online users</h3>
        <ul className="user-list">
          {onlineUsers.map((user) => (
            <li key={user.id} className="user-item">
              <div className="avatar-container">
                <img src="avatar.jpg" className="avatar" />
                <div className="online-indicator"></div>
              </div>
              <h3 className="user-name">{user.name}</h3>
            </li>
          ))}
        </ul>
      </div>

      <div className="communities">
        <h3>Community</h3>
        <ul className="community-list">
          {communities.map((community) => (
            <li key={community.id} className="community-item">
              <img
                src={community.avatar}
                alt={community.name}
                className="avatar"
              />
              <span className="community-name">{community.name}</span>
            </li>
          ))}
        </ul>
      </div>

      <button className="view-more-btn">View more</button>
    </div>
  );
}

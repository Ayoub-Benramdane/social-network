"use client";
import { useState } from "react";
import "../styles/LeftSideBar.css";
import UserCard from "./UserCard";

export default function LeftSidebar({ users, bestcategories }) {
  const [searchQuery, setSearchQuery] = useState("");

  const handleSearch = (e) => {
    setSearchQuery(e.target.value);
  };

  const handleFollow = (userId) => {
    console.log(`Following user with ID: ${userId}`);
  };

  return (
    <div className="left-sidebar">
      <div className="search-box">
        <div className="search-input-container">
          <img src="/icons/search.svg" className="search-icon" alt="Search" />
          <input
            type="text"
            placeholder="Search users and groups"
            value={searchQuery}
            onChange={handleSearch}
          />
        </div>
      </div>

      <div className="sidebar-section">
        <div className="section-header">
          <h3>Suggested Users</h3>
          <button className="see-all-btn">
            See all <span className="arrow">→</span>
          </button>
        </div>

        <ul className="user-list">
          {users.map((user) => (
            <UserCard key={user.id} user={user} action={"follow"} />
            // <li key={user.id} className="user-item">
            //   <img
            //     src={
            //       user.avatar ||
            //       "./avatars/thorfinn-vinland-saga-episode-23-1.png"
            //     }
            //     className="user-avatar"
            //     alt={user.username}
            //   />
            //   <div className="user-details">
            //     <div className="user-info">
            //       <h4 className="user-name">{`${user.first_name} ${user.last_name}`}</h4>
            //       <p className="user-username">@{user.username}</p>
            //     </div>
            //     <button
            //       className="follow-btn"
            //       onClick={() => handleFollow(user.id)}
            //     >
            //       Follow
            //     </button>
            //   </div>
            // </li>
          ))}
        </ul>
      </div>

      <div className="sidebar-section">
        <div className="section-header">
          <h3>Popular Categories</h3>
          <button className="see-all-btn">
            See all <span className="arrow">→</span>
          </button>
        </div>

        <ul className="category-list">
          {bestcategories.map((category) => (
            <li key={category.id} className="category-item">
              <div className="category-icon">
                <img src={`/icons/${category.name}.png`} alt={category.name} />
              </div>
              <span className="category-name">{category.name}</span>
              <span className="category-count">
                {category.count || 0} posts
              </span>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}

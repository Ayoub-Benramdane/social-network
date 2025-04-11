"use client";
import "../styles/LeftSideBar.css";

export default function LeftSidebar({ users, bestcategories }) {
  return (
    <div className="left-sidebar">
      <div className="search-box">
        <div className="search-input-container">
          <input type="text" placeholder="Search users and groups" />
          <img
            src="/icons/search.svg"
            className="search-icon"
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="#9CA3AF"
          ></img>
        </div>
      </div>

      <div className="users">
        <div className="header">
          <h3>Not Following</h3>
          <p style={{ color: "#3555F9", fontSize: "14px", cursor: "pointer" }}>
            See all &rarr;
          </p>
        </div>
        <ul className="user-list">
          {users.map((user) => (
            <li key={user.id} className="user-item">
              {/* <div className="avatar-container"> */}
              <img
                src="./avatars/thorfinn-vinland-saga-episode-23-1.png"
                className="avatar"
              />
              {/* </div> */}
              <div className="name-follow">
                <div>
                  <h4 className="user-name">{`${user.first_name} ${user.last_name}`}</h4>
                  <p className="user-username">{user.username}</p>
                </div>
                <button className="follow-btn">Follow</button>
              </div>
            </li>
          ))}
        </ul>
      </div>

      <div className="categories">
        <div className="header">
          <h3>Best Categories</h3>
          <p style={{ color: "#3555F9", fontSize: "14px", cursor: "pointer" }}>
            See all &rarr;
          </p>
        </div>
        <ul className="category-list">
          {bestcategories.map((category) => (
            <li key={category.id} className="category-item">
              <img
                src={`/icons/${category.name}.png`}
                // alt={category.name}
                className="avatar"
              />
              <span className="category-name">{category.name}</span>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}

"use client";
import "../styles/ProfileCard.css";

export default function ProfileCard({ user }) {

  async function handleCreatePost() {
    try {
      const response = await fetch("http://localhost:8404/new_post", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
      });
      const data = await response.json();
      if (!response.ok) {
        console.log(data);
      }
      console.log(data);
    } catch (error) {
      console.log(error);
    }
  }

  return (
    <div className="profile-card">
      <div className="profile-image">
        <img src={user.avatar} alt={user.username} />
      </div>
      <div className="profile-info">
        <p className="profile-name">{`${user.first_name} ${user.last_name}`}</p>

        <div className="stats">
          <div className="stat">
            <div className="stat-value">{user.total_followers}</div>
            <div className="stat-label">Followers</div>
          </div>
          <div className="stat">
            <div className="stat-value">{user.total_following}</div>
            <div className="stat-label">Followings</div>
          </div>
          <div className="stat">
            <div className="stat-value">{user.total_posts}</div>
            <div className="stat-label">Posts</div>
          </div>
        </div>

        <button onClick={handleCreatePost} className="create-post-btn">
          <span>Create post</span>
        </button>
        <button onClick={handleCreatePost} className="create-group-btn">
          <span>Create new group</span>
        </button>
      </div>
    </div>
  );
}

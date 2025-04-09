"use client";
import "../styles/ProfileCard.css";

export default function ProfileCard({ user }) {
  const defaultUser = {
    name: "Amine Dinani",
    avatar: "https://via.placeholder.com/150",
    followers: 30,
    following: 30,
    posts: 30,
  };

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

  const userData = user || defaultUser;

  return (
    <div className="profile-card">
      <div className="profile-image">
        <img src="amine.jpeg" alt={userData.name} />
      </div>
      <div className="profile-info">
        <h3 className="profile-name">{userData.name}</h3>

        <div className="stats">
          <div className="stat">
            <div className="stat-value">{userData.followers}</div>
            <div className="stat-label">Followers</div>
          </div>
          <div className="stat">
            <div className="stat-value">{userData.following}</div>
            <div className="stat-label">Followings</div>
          </div>
          <div className="stat">
            <div className="stat-value">{userData.posts}</div>
            <div className="stat-label">Posts</div>
          </div>
        </div>

        <button onClick={handleCreatePost} className="create-post-btn">
          <span>+ Create post</span>
        </button>
      </div>
    </div>
  );
}

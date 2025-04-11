"use client";
import { useState } from "react";
import "../styles/PostComponent.css";

export default function PostComponent({ posts }) {
  // const [showComments, setShowComments] = useState(false);
  console.log("Helloo", posts);
  function handleLike(postID) {
    console.log("You liked post: ", postID);
  }

  return (
    <div>
      {posts.map((post) => (
        <div key={post.id} className="post-card">
          <div className="header">
            <div className="post-header">
              <img
                src="avatar.jpg"
                alt={post.author}
                className="author-avatar"
              />
              <div className="author-info">
                <h4 className="author-name">{post.author}</h4>
                <p
                  className="created-at"
                  style={{ color: "#B8C3E1", fontSize: "12px" }}
                >
                  {post.created_at}
                </p>
              </div>
            </div>
            <div className="post-privacy">
              {/* <p className="privacy-text">{post.privacy}</p> */}
              <img
                src={`./icons/${post.privacy}.svg`}
                width={"32px"}
                height={"32px"}
                className="privacy-icon"
              ></img>
            </div>
          </div>

          <div className="post-content">
            <h3 className="post-title">{post.title}</h3>
            <p className="post-text">{post.content}</p>
            <div className="post-category">{post.category}</div>
          </div>
          <div className="post-actions">
            <div className="action-like" onClick={handleLike(post.id)}>
              <img src="/icons/like.svg"></img>
              <p>{post.total_likes} Likes</p>
            </div>
            <div className="action-comment">
              <img src="/icons/comment.svg"></img>
              <p>{post.total_comments} Comments</p>
            </div>
          </div>
          {/* <hr style={{ width: "60%", color: "#eeeeee" }}></hr> */}
          <p
            style={{
              color: "#3555F9",
              fontSize: "14px",
              cursor: "pointer",
              // textAlign: "center",
              padding: "1rem",
            }}
          >
            See post &rarr;
          </p>
        </div>
      ))}
    </div>
  );
}

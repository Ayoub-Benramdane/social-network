"use client";
import { useState } from "react";
import "../styles/PostComponent.css";

export default function PostComponent({ posts }) {
  function handleLike(postID) {
    console.log("You liked post: ", postID);
  }

  return (
    <div className="posts-container">
      {posts.map((post) => (
        <div key={post.id} className="post-card">
          <div className="header">
            <div className="post-header">
              <img
                src={post.avatar || "avatar.jpg"}
                alt={post.author}
                className="author-avatar"
              />
              <div className="author-info">
                <h4 className="author-name">{post.author}</h4>
                <div className="timestamp">
                  <img src="./icons/created_at.svg" alt="Time" />
                  <p className="created-at">{post.created_at}</p>
                </div>
              </div>
            </div>
            <div className="post-privacy">
              <img
                src={`./icons/${post.privacy}.svg`}
                width="32"
                height="32"
                className="privacy-icon"
                alt={post.privacy}
              />
            </div>
          </div>

          <div className="post-content">
            <h3 className="post-title">{post.title}</h3>
            <p className="post-text">{post.content}</p>

            {/* Display post image if available */}
            {post.image && (
              <div className="post-image-container">
                <img
                  src={post.image}
                  alt="Post content"
                  className="post-image"
                />
              </div>
            )}

            <div className="post-category">{post.category}</div>
          </div>

          <div className="post-actions">
            <button
              className="action-button action-like"
              onClick={() => handleLike(post.id)}
            >
              <img src="/icons/like.svg" alt="Like" />
              <span>{post.total_likes} Likes</span>
            </button>

            <button className="action-button action-comment">
              <img src="/icons/comment.svg" alt="Comment" />
              <span>{post.total_comments} Comments</span>
            </button>
          </div>

          <button className="see-post-button">
            See post <span className="arrow">→</span>
          </button>
        </div>
      ))}
    </div>
  );
}

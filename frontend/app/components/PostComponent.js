"use client";
import { useState } from "react";
import "../styles/PostComponent.css";

export default function PostComponent({ post }) {
  const [showComments, setShowComments] = useState(false);

  const defaultPost = {
    id: 1,
    author: "Abdelkhaled Laidi",
    authorAvatar: "https://via.placeholder.com/40",
    title: "Neque porro quisquam est qui dolor, adipisci velit...",
    content:
      "Lorem ipsum dolor sit amet consectetur. Vel in molestie amet tellus dui. At arcu mauris purus magna . Volutpat dictum consectetur semper aliquam  volutpat. Felis tincidunt donec blandit mauris ornare orci. Tellus enim in eget in it urna sit habitant in et ac. Metus est elit risus leo enim viverra morbi eu.",
    category: "Technology",
    createdAt: "2 hours ago",
    likes: 23,
    comments: [
      {
        id: 1,
        author: "Abdelkhaled Laidi",
        authorAvatar: "https://via.placeholder.com/40",
        content:
          "Lorem ipsum dolor sit amet consectetur. Vel in molestie amet tellus dui. At arcu mauris purus magna . Volutpat dictum consectetur semper aliquam  volutpat. Felis tincidunt donec blandit mauris ornare orci. Tellus enim in eget in it urna sit habitant in et ac. Metus est elit risus leo enim viverra morbi eu.",
      },
    ],
    commentsCount: 18,
  };

  const postData = post || defaultPost;

  return (
    <div className="post-card">
      <div className="post-header">
        <img src="avatar.jpg" alt={postData.author} className="author-avatar" />
        <div className="author-info">
          <h4 className="author-name">{postData.author}</h4>
          <p className="post-time">{postData.createdAt}</p>
        </div>
      </div>

      <div className="post-content">
        <h3 className="post-title">{postData.title}</h3>
        <p className="post-text">{postData.content}</p>

        <div className="post-category">{postData.category}</div>
      </div>

      <div className="post-actions">
        <div className="like-action">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="#2563EB">
            <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z" />
          </svg>
          <span>{postData.likes} likes</span>
        </div>

        <div className="comments-action">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="#6B7280">
            <path d="M21.99 4c0-1.1-.89-2-1.99-2H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h14l4 4-.01-18zM18 14H6v-2h12v2zm0-3H6V9h12v2zm0-3H6V6h12v2z" />
          </svg>
          <span>{postData.commentsCount} comments</span>
        </div>
      </div>

      {/* Comments Section */}
      <div className="comments-section">
        <h4 className="comments-title">Comments</h4>

        {postData.comments.map((comment) => (
          <div key={comment.id} className="comment">
            <div className="comment-header">
              <img
                src={comment.authorAvatar}
                alt={comment.author}
                className="comment-avatar"
              />
              <h5 className="comment-author">{comment.author}</h5>
            </div>
            <p className="comment-text">{comment.content}</p>
          </div>
        ))}

        <button className="view-more-btn">View more</button>
      </div>
    </div>
  );
}

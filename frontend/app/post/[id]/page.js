"use client";
import { useState, useEffect } from "react";
// import { useRouter } from "next/navigation";
import { useParams } from "next/navigation";
import { useRouter } from "next/router";
import Navbar from "../../components/NavBar";
import "../../styles/PostPage.css";

export default function PostPage() {
  //   const router = useRouter();
  const { id } = useParams();

  const [post, setPost] = useState(null);
  const [comments, setComments] = useState([]);
  const [newComment, setNewComment] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);


  useEffect(() => {
    if (!id) return;

    const fetchPost = async () => {
      setIsLoading(true);
      try {
        const response = await fetch(
          `http://localhost:8404/post?post_id=${id}&group_id=0`,
          {
            method: "GET",
            credentials: "include",
          }
        );

        if (!response.ok) {
          throw new Error("Failed to fetch post");
        }

        const data = await response.json();
        setPost(data);
        console.log(data);
        if (post) {
          console.log(data);
        }

        // setComments(data.Comments || []);
      } catch (err) {
        console.log("Error fetching post:", err);
        setError("Failed to load post. Please try again later.");
      } finally {
        setIsLoading(false);
      }
    };

    fetchPost();
  }, [id]);

  const handleCommentSubmit = async (e) => {
    e.preventDefault();
    if (!newComment.trim()) return;

    try {
      const response = await fetch("http://localhost:8404/comments", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          postId: parseInt(id),
          content: newComment,
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to add comment");
      }

      const data = await response.json();

      // Add the new comment to the list
      setComments([...comments, data]);

      // Update the total comments count
      setPost({
        ...post,
        TotalComments: post.TotalComments + 1,
      });

      // Clear the comment input
      setNewComment("");
    } catch (err) {
      console.error("Error adding comment:", err);
      alert("Failed to add comment. Please try again.");
    }
  };

  const handleLike = async () => {
    try {
      const response = await fetch(`http://localhost:8404/like`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          postId: parseInt(id),
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to like post");
      }

      // Update the total likes count
      setPost({
        ...post,
        TotalLikes: post.TotalLikes + 1,
      });
    } catch (err) {
      console.error("Error liking post:", err);
      alert("Failed to like post. Please try again.");
    }
  };

  if (isLoading) {
    return (
      <div className="post-page-container">
        {/* <Navbar /> */}
        <div className="post-page-content">
          <div className="loading-spinner">Loading...</div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="post-page-container">
        {/* <Navbar /> */}
        <div className="post-page-content">
          <img
            src="/icons/no-post-state.svg"
            alt="Error"
            className="error-icon"
          />
          <div className="error-message">{error}</div>
          {/* <button className="back-button" onClick={() => router.push("/")}>
            Back to Home
          </button> */}
        </div>
      </div>
    );
  }

  if (!post) {
    return (
      <div className="post-page-container">
        {/* <Navbar /> */}
        <div className="post-page-content">
          <div className="error-message">Post not found</div>
          {/* <button className="back-button" onClick={() => router.push("/")}>
            Back to Home
          </button> */}
        </div>
      </div>
    );
  }

  return (
    <div className="post-page-container">
      {/* <Navbar /> */}
      <div className="post-page-content">
        <div className="post-page-card">
          <div className="post-page-header">
            <div className="post-author-info">
              <img
                src={post.avatar}
                alt={post.avatar}
                className="author-avatar"
              />
              <div className="author-details">
                <a href={`/profile/${post.user_id}`} className="author-link">
                  <h4 className="author-name">{post.author}</h4>
                </a>
                <div className="timestamp">
                  <img src="/icons/created_at.svg" alt="Created at" />
                  <p className="created-at">{post.created_at}</p>
                </div>
              </div>
            </div>
            <div className="post-privacy">
              <img
                src={`/icons/${post.privacy}.svg`}
                width={"32px"}
                height={"32px"}
                className="privacy-icon"
                alt={post.privacy}
              />
            </div>
          </div>

          <div className="post-page-body">
            <h2 className="post-title">{post.title}</h2>
            <p className="post-content">{post.content}</p>

            {post.image && (
              <div className="post-image-container">
                <img src={post.image} alt="Post image" className="post-image" />
              </div>
            )}

            <div className="post-category">{post.category}</div>
          </div>

          <div className="post-actions">
            <div className="action-like" onClick={handleLike}>
              <img src="/icons/like.svg" alt="Like" />
              <p>{post.total_likes} Likes</p>
            </div>
            <div className="action-comment">
              <img src="/icons/comment.svg" alt="Comment" />
              <p>{post.total_comments} Comments</p>
            </div>
          </div>

          <div className="post-comments-section">
            <h3 className="comments-title">Comments</h3>

            <div className="comments-list">
              {comments.length > 0 ? (
                comments.map((comment) => (
                  <div key={comment.ID} className="comment-item">
                    <div className="comment-header">
                      <img
                        src="avatar.jpg"
                        alt={comment.Author}
                        className="comment-avatar"
                      />
                      <div className="comment-author-info">
                        <h4 className="comment-author">{comment.Author}</h4>
                        <p className="comment-time">{comment.CreatedAt}</p>
                      </div>
                    </div>
                    <p className="comment-content">{comment.Content}</p>
                  </div>
                ))
              ) : (
                <p className="no-comments">
                  No comments yet. Be the first to comment!
                </p>
              )}
            </div>

            <form className="comment-form" onSubmit={handleCommentSubmit}>
              <textarea
                className="comment-input"
                placeholder="Write a comment..."
                value={newComment}
                onChange={(e) => setNewComment(e.target.value)}
                rows={3}
              />
              <button type="submit" className="comment-submit-btn">
                Post Comment
              </button>
            </form>
          </div>
        </div>

        {/* <button className="back-button" onClick={() => router.push("/")}>
          Back to Home
        </button> */}
      </div>
    </div>
  );
}

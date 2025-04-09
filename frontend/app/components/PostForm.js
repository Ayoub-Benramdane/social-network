import { useState } from "react";
import "../styles/PostFormStyle.css";

export default function PostForm() {
  const [postFormInput, setPostFormInput] = useState({
    title: "",
    content: "",
    privacy: "",
    category: "technology",
    postImage: null,
  });
  const [imageInputKey, setImageInputKey] = useState(Date.now());

  const handleImageChange = (event) => {
    const file = event.target.files[0];
    setPostFormInput({
      ...postFormInput,
      postImage: file,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    const formData = new FormData();

    // formData.append(
    //   "postData",
    //   JSON.stringify({
    //     title: postFormInput.title,
    //     content: postFormInput.content,
    //     privacy: postFormInput.privacy,
    //     category: postFormInput.category,
    //   })
    // );
    const fieldsToInclude = ["title", "content", "privacy", "category"];

    fieldsToInclude.forEach((field) => {
      formData.append(field, postFormInput[field]);
    });

    if (postFormInput.postImage) {
      formData.append("postImage", postFormInput.postImage);
    }
    console.log(postFormInput);

    try {
      const response = await fetch("http://localhost:8404/new_post", {
        method: "POST",

        body: formData,
        credentials: "include",
      });

      if (!response.ok) {
        const data = await response.json();
        console.log(data);

        throw new Error(data.error || "Failed to create the post");
      }

      console.log("Post created successfully");
      setPostFormInput({
        title: "",
        content: "",
        privacy: "",
        category: "",
        postImage: null,
      });
    } catch (error) {
      console.log(error);
    }
  };

  return (
    <div className="postDiv">
      <h3 style={{ color: "#18151b" }}>Create a new Post</h3>
      <form className="create-post-form" onSubmit={handleSubmit}>
        <div className="form-div">
          <div className="title-input">
            <input
              placeholder="title"
              required
              value={postFormInput.title}
              onChange={(e) => {
                setPostFormInput({ ...postFormInput, title: e.target.value });
              }}
            ></input>
          </div>

          <div className="content-input">
            <input
              placeholder="content"
              required
              value={postFormInput.content}
              onChange={(e) => {
                setPostFormInput({ ...postFormInput, content: e.target.value });
              }}
            ></input>
          </div>
          <div className="image-input">
            <p>Upload Image</p>
            {postFormInput.postImage && (
              <div>
                <img
                  style={{
                    width: "150px",
                    borderRadius: "8px",
                  }}
                  src={URL.createObjectURL(postFormInput.postImage)}
                  alt="Selected"
                />
                <button
                  onClick={(e) => {
                    e.preventDefault();
                    setPostFormInput({
                      ...postFormInput,
                      postImage: null,
                    });
                    setImageInputKey(Date.now());
                  }}
                >
                  Remove
                </button>
              </div>
            )}
            <input
              key={imageInputKey}
              type="file"
              name="postImage"
              onChange={handleImageChange}
            ></input>
          </div>

          <div className="privacy-input">
            <input
              required
              type="radio"
              value="private"
              name="privacy"
              checked={postFormInput.privacy === "private"}
              onChange={(e) => {
                setPostFormInput({ ...postFormInput, privacy: e.target.value });
              }}
            />{" "}
            Private
            <input
              type="radio"
              value="public"
              name="privacy"
              checked={postFormInput.privacy === "public"}
              onChange={(e) => {
                setPostFormInput({ ...postFormInput, privacy: e.target.value });
              }}
            />{" "}
            Public
            <input
              type="radio"
              value="almost-public"
              name="privacy"
              checked={postFormInput.privacy === "almost-private"}
              onChange={(e) => {
                setPostFormInput({ ...postFormInput, privacy: e.target.value });
              }}
            />{" "}
            Almost Private
          </div>
          <div className="category">
            <select
              onChange={(e) => {
                setPostFormInput({
                  ...postFormInput,
                  category: e.target.value,
                });
              }}
            >
              <option value="technology">Technology</option>
              <option value="sport">Sport</option>
              <option value="entertainment">Entertainment</option>
              <option value="other">Other</option>
            </select>
          </div>
        </div>
        <button
          style={{
            backgroundColor: "rgb(174, 0, 255)",
            width: "150px",
            height: "32px",
            border: "none",
            color: "white",
            borderRadius: "8px",
            fontSize: "14px",
            fontWeight: "500",
          }}
        >
          Publish
        </button>
      </form>
      {/* <div className="posts-list">
        <h3>Posts</h3>
        {posts.length > 0 ? (
          posts.map((post, index) => (
            <div key={index} className="post">
              <h4>{post.Title}</h4>
              <p>{post.Content}</p>
              <img
                src={post.Image}
                alt="Post"
                style={{ width: "100px", height: "100px" }}
              />
              <div>Categories: {post.Categories.join(", ")}</div>
              <div>Privacy: {post.Privacy}</div>
            </div>
          ))
        ) : (
          <p>No posts yet.</p>
        )}
      </div> */}
    </div>
  );
}

import { useState } from "react";
import "../styles/PostFormStyle.css";

export default function PostForm() {
  const [postFormInput, setPostFormInput] = useState({
    title: "",
    content: "",
    privacy: "",
  });
  return (
    <div className="postDiv">
      <h3>Create a new Post</h3>
      <form
        className="create-post-form"
        onSubmit={(e) => {
          e.preventDefault();
          handleSubmit(postFormInput);
          //   console.log(postFormInput);
        }}
      >
        <div className="form-div">
          <div className="title-input">
            <input
              placeholder="title"
              required
              onChange={(e) => {
                setPostFormInput({ ...postFormInput, title: e.target.value });
              }}
            ></input>
          </div>

          <div className="content-input">
            <input
              placeholder="content"
              required
              onChange={(e) => {
                setPostFormInput({ ...postFormInput, content: e.target.value });
              }}
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
              checked={postFormInput.privacy === "almost-public"}
              onChange={(e) => {
                setPostFormInput({ ...postFormInput, privacy: e.target.value });
              }}
            />{" "}
            Almost Public
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
    </div>
  );
}

async function handleSubmit(postFormInput) {
  try {
    const response = await fetch("http://localhost:8080/api/create_post", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(postFormInput),
    });
    if (!response.ok) {
      throw new Error("Failed to create the post");
    }
    console.log("Post created successfully");
  } catch (error) {
    console.log(error);
  }
}

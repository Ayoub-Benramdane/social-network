import { useState } from "react";
import "../styles/RegisterFormStyle.css";

export default function RegisterForm() {
  const [registerFormInputs, setRegisterFormInputs] = useState({
    email: "",
    username: "",
    firstName: "",
    lastName: "",
    password: "",
    confirmedPassword: "",
    dateOfBirth: "",
    aboutMe: "",
    avatar: null,
  });

  const [passwordError, setPasswordError] = useState("");
  const [imageInputKey, setImageInputKey] = useState(Date.now());

  const handleImageChange = (event) => {
    const file = event.target.files[0];
    setRegisterFormInputs({
      ...registerFormInputs,
      avatar: file,
    });
  };
  const handleSubmit = async (e) => {
    e.preventDefault();
    // console.log(registerFormInputs);
    const formData = new FormData();

    const fieldsToInclude = [
      "username",
      "firstName",
      "lastName",
      "email",
      "dateOfBirth",
      "password",
      "confirmedPassword",
      "aboutMe",
    ];

    fieldsToInclude.forEach((field) => {
      formData.append(field, registerFormInputs[field]);
    });

    if (registerFormInputs.avatar) {
      formData.append("avatar", registerFormInputs.avatar);
    }
    if (registerFormInputs.password !== registerFormInputs.confirmedPassword) {
      setPasswordError("Passwords do not match!");
    }

    setPasswordError("");

    const response = await fetch("http://localhost:8404/register", {
      method: "POST",
      body: formData,
    });

    // const data = await response.json();
    // console.log(data);

    const message = (document.querySelector(
      ".message"
    ).textContent = `Your registration has been done successfully.`);
  };

  return (
    <div className="registerDiv">
      <form className="register-form" onSubmit={handleSubmit}>
        <h3
          style={{ fontSize: "40px", textAlign: "center", fontWeight: "600" }}
        >
          Register
        </h3>
        <div className="inputs">
          {/* Email */}
          <div className="email">
            <label>Email:</label>
            <input
              required
              onChange={(e) => {
                setRegisterFormInputs({
                  ...registerFormInputs,
                  email: e.target.value,
                });
              }}
            />
          </div>

          {/* Username */}
          <div className="username">
            <label>Username:</label>
            <input
              required
              onChange={(e) => {
                setRegisterFormInputs({
                  ...registerFormInputs,
                  username: e.target.value,
                });
              }}
            />
          </div>

          {/* First Name */}
          <div className="first-name">
            <label>First Name:</label>
            <input
              required
              onChange={(e) => {
                setRegisterFormInputs({
                  ...registerFormInputs,
                  firstName: e.target.value,
                });
              }}
            />
          </div>

          {/* Last Name */}
          <div className="last-name">
            <label>Last Name:</label>
            <input
              required
              onChange={(e) => {
                setRegisterFormInputs({
                  ...registerFormInputs,
                  lastName: e.target.value,
                });
              }}
            />
          </div>

          {/* Password */}
          <div className="password">
            <label>Password:</label>
            <input
              required
              type="password"
              onChange={(e) => {
                setRegisterFormInputs({
                  ...registerFormInputs,
                  password: e.target.value,
                });
              }}
            />
          </div>

          {/* Confirm Password */}
          <div className="password">
            <label>Confirm password:</label>
            <input
              required
              type="password"
              onChange={(e) => {
                setRegisterFormInputs({
                  ...registerFormInputs,
                  confirmedPassword: e.target.value,
                });

                if (e.target.value === registerFormInputs.password) {
                  setPasswordError("");
                }
              }}
            />
            {passwordError && <p style={{ color: "red" }}>{passwordError}</p>}
          </div>

          {/* Avatar */}

          <div className="image-input">
            <p>Upload Image</p>
            {registerFormInputs.avatar && (
              <div>
                <img
                  style={{
                    width: "150px",
                    borderRadius: "8px",
                  }}
                  src={URL.createObjectURL(registerFormInputs.avatar)}
                  alt="Selected"
                />
                <button
                  onClick={(e) => {
                    setRegisterFormInputs({
                      ...registerFormInputs,
                      avatar: null,
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
              name="avatar"
              onChange={handleImageChange}
            ></input>
          </div>

          {/* About Me */}
          <div className="about-me">
            <label>About Me:</label>
            <textarea
              required
              onChange={(e) => {
                setRegisterFormInputs({
                  ...registerFormInputs,
                  aboutMe: e.target.value,
                });
              }}
            />
          </div>
        </div>

        <div className="date-of-birth">
          <label htmlFor="birthday">Birthday:</label>
          <input
            type="date"
            onChange={(e) => {
              setRegisterFormInputs({
                ...registerFormInputs,
                dateOfBirth: e.target.value,
              });
            }}
          ></input>
        </div>

        <p className="message"></p>
        <button className="register-btn">Register</button>
      </form>
      {/* <p>Already have an acoount?</p>
      <a href="api/login" style={{ color: "blue" }}>
        Login in
      </a> */}
    </div>
  );
}

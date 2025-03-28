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
    // aboutMe: "",
  });

  const [passwordError, setPasswordError] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    console.log(registerFormInputs);
    console.log(typeof registerFormInputs.dateOfBirth);

    if (registerFormInputs.password !== registerFormInputs.confirmedPassword) {
      setPasswordError("Passwords do not match!");
    }

    setPasswordError("");

    const response = await fetch("http://localhost:8404/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(registerFormInputs),
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

          {/* About Me */}
          {/* <div className="about-me">
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
          </div> */}
        </div>

        <div className="date-of-birth">
          <label htmlFor="birthday">Birthday:</label>
          <input
            type="date"
            onChange={(e) => {
              const formattedDate = new Date(e.target.value).toISOString();

              setRegisterFormInputs({
                ...registerFormInputs,
                dateOfBirth: formattedDate,
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

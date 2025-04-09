// "use client";
// import { useState } from 'react';

// export default function LoginForm() {
//     // console.log(props);

//     const [loginFormInputs, setLoginFormInputs] = useState({
//         username: "",
//         password: ""
//     });

//     const handleSubmit = async (e) => {
//         e.preventDefault();
//         if (loginFormInputs.username = "username" && loginFormInputs.password === "password") {
//             console.log("Logedin");

//         }
//         const response = await fetch("http://localhost:8080/api/login", {
//             method: "POST",
//             headers: {
//                 "Content-Type": "application/json",
//             },
//             body: JSON.stringify(loginFormInputs),
//         });
//         console.log(loginFormInputs);

//         const data = await response.json();
//         console.log(data);
//     };

//     return (
//         <div>
//             <form onSubmit={handleSubmit}>
//                 <label>
//                     Enter your username/email:
//                 </label>
//                 <input
//                     onChange={(e) => setLoginFormInputs({...loginFormInputs, username: e.target.value})}
//                 />
//                 <label>
//                     Enter your password:
//                 </label>
//                 <input
//                     type="password"
//                     onChange={(e) => setLoginFormInputs({...loginFormInputs, password: e.target.value})}
//                 />
//                 <button type="submit">Login</button>
//             </form>
//         </div>
//     );
// }

// "use client";
// import { useState } from "react";

// export default function LoginForm({ onLoginSuccess }) {
//   const [loginFormInputs, setLoginFormInputs] = useState({
//     username: "",
//     password: "",
//   });

//   const handleSubmit = async (e) => {
//     e.preventDefault();

//     if (loginFormInputs.username === "username" && loginFormInputs.password === "password") {
//       console.log("Logged in");

//       onLoginSuccess();

//       return;
//     }

//     const response = await fetch("http://localhost:8404/login", {
//       method: "POST",
//       headers: {
//         "Content-Type": "application/json",
//       },
//       body: JSON.stringify(loginFormInputs),
//     });
//     console.log(loginFormInputs);

//     const data = await response.json();
//     console.log(data);
//   };

//   return (
//     <div>
//       <form onSubmit={handleSubmit}>
//         <label>Enter your username/email:</label>
//         <input
//           onChange={(e) =>
//             setLoginFormInputs({ ...loginFormInputs, username: e.target.value })
//           }
//         />
//         <label>Enter your password:</label>
//         <input
//           type="password"
//           onChange={(e) =>
//             setLoginFormInputs({ ...loginFormInputs, password: e.target.value })
//           }
//         />
//         <button type="submit">Login</button>
//       </form>
//     </div>
//   );
// }

"use client";

import { useState } from "react";

export default function LoginForm({ onLoginSuccess }) {
  const [loginFormInputs, setLoginFormInputs] = useState({
    email: "",
    password: "",
  });

  const [errorMessage, setErrorMessage] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!loginFormInputs.email || !loginFormInputs.password) {
      setErrorMessage("Email and password are required.");
      return;
    }

    const response = await fetch("http://localhost:8404/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify(loginFormInputs),
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();

    if (response.ok) {
      // const data = await response.json();
      console.log(data);

      onLoginSuccess();
    } else {
      console.error("Login failed:", data.error || "Unknown error");
      setErrorMessage(data.error || "Login failed");
    }
  };

  return (
    <div>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Email:</label>
          <input
            type="email"
            value={loginFormInputs.email}
            onChange={(e) =>
              setLoginFormInputs({
                ...loginFormInputs,
                email: e.target.value,
              })
            }
          />
        </div>

        <div>
          <label>Password:</label>
          <input
            type="password"
            value={loginFormInputs.password}
            onChange={(e) =>
              setLoginFormInputs({
                ...loginFormInputs,
                password: e.target.value,
              })
            }
          />
        </div>

        <button type="submit">Login</button>
      </form>

      {errorMessage && <div className="error-message">{errorMessage}</div>}
    </div>
  );
}

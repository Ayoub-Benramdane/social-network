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

"use client";
import { useState } from "react";

export default function LoginForm({ onLoginSuccess }) {
  const [loginFormInputs, setLoginFormInputs] = useState({
    username: "",
    password: "",
  });

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (loginFormInputs.username === "username" && loginFormInputs.password === "password") {
      console.log("Logged in");

      onLoginSuccess();

      return; 
    }

    const response = await fetch("http://localhost:8404/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(loginFormInputs),
    });
    console.log(loginFormInputs);

    const data = await response.json();
    console.log(data);
  };

  return (
    <div>
      <form onSubmit={handleSubmit}>
        <label>Enter your username/email:</label>
        <input
          onChange={(e) =>
            setLoginFormInputs({ ...loginFormInputs, username: e.target.value })
          }
        />
        <label>Enter your password:</label>
        <input
          type="password"
          onChange={(e) =>
            setLoginFormInputs({ ...loginFormInputs, password: e.target.value })
          }
        />
        <button type="submit">Login</button>
      </form>
    </div>
  );
}

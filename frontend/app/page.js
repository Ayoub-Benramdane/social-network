// // "use client";

// // import LoginForm from "./components/LoginForm";
// // import RegisterForm from "./components/RegisterForm";
// // import PostForm from "./components/PostForm";

// // export default function Home() {
  
// //   return (

// //     // <PostForm />
// //     <RegisterForm />
// //     //  <LoginForm />
// //   );
// // }


"use client";

import { useState } from "react";
import LoginForm from "./components/LoginForm";
import RegisterForm from "./components/RegisterForm";
import PostForm from "./components/PostForm";

export default function Home() {
  const [isLogin, setIsLogin] = useState(true);
  const [isLoggedIn, setIsLoggedIn] = useState(false); 

  const toggleForm = () => {
    setIsLogin(!isLogin);
  };

  const handleLoginSuccess = () => {
    setIsLoggedIn(true); 
  };

  return (
    <div>
      {!isLoggedIn && (
        <button onClick={toggleForm}>
          {isLogin ? "Register a new account" : "Have an account"}
        </button>
      )}

      {isLoggedIn ? (
        <PostForm />
      ) : (
        isLogin ? (
          <LoginForm onLoginSuccess={handleLoginSuccess} />
        ) : (
          <RegisterForm />
        )
      )}
    </div>
  );
}

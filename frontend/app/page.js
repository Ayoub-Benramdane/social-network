"use client";

import { useState, useEffect } from "react";
import Navbar from "./components/NavBar";
import LeftSidebar from "./components/LeftSideBar";
import ProfileCard from "./components/ProfileCard";
import TopGroups from "./components/TopGroups";
import PostComponent from "./components/PostComponent";
import LoginForm from "./components/LoginForm";
import RegisterForm from "./components/RegisterForm";
import "./styles/page.css";

export default function Home() {
  const [isLogin, setIsLogin] = useState(true);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [homeData, setHomeData] = useState(null);


  useEffect(() => {
    const checkLoginStatus = async () => {
      try {
        const response = await fetch("http://localhost:8404/session", {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
        });

        if (response.ok) {
          const data = await response.json();
          if (data === true) {
            setIsLoggedIn(true);
          } else {
            setIsLoggedIn(false);
          }
        }
      } catch (error) {
        console.error("Error checking login status:", error);
      }
    };

    checkLoginStatus();
  }, []);

  useEffect(() => {
    if (isLoggedIn) {
      fetch("http://localhost:8404/", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
      })
        .then((response) => response.json())
        .then((data) => {
          setHomeData(data);
          console.log( "Data received : ",data);
          
        })
        .catch((error) => {
          console.error("Error fetching posts:", error);
        });
    }
  }, [isLoggedIn]);
  
  
  // console.log("All:", posts);

  const toggleForm = () => {
    setIsLogin(!isLogin);
  };

  const handleLoginSuccess = () => {
    setIsLoggedIn(true);
  };

  // const samplePosts = [
  //   {
  //     id: 1,
  //     author: "Abdelkhaled Laidi",
  //     authorAvatar: "https://via.placeholder.com/40",
  //     title: "Neque porro quisquam est qui dolor, adipisci velit...",
  //     content:
  //       "Lorem ipsum dolor sit amet consectetur. Vel in molestie amet tellus dui. At arcu mauris purus magna . Volutpat dictum consectetur semper aliquam volutpat. Felis tincidunt donec blandit mauris ornare orci. Tellus enim in eget in it urna sit habitant in et ac. Metus est elit risus leo enim viverra morbi eu.",
  //     category: "Technology",
  //     createdAt: "2 hours ago",
  //     likes: 23,
  //     comments: [
  //       {
  //         id: 1,
  //         author: "Abdelkhaled Laidi",
  //         authorAvatar: "https://via.placeholder.com/40",
  //         content:
  //           "Lorem ipsum dolor sit amet consectetur. Vel in molestie amet tellus dui. At arcu mauris purus magna . Volutpat dictum consectetur semper aliquam volutpat. Felis tincidunt donec blandit mauris ornare orci. Tellus enim in eget in it urna sit habitant in et ac. Metus est elit risus leo enim viverra morbi eu.",
  //       },
  //     ],
  //     commentsCount: 18,
  //   },
  //   {
  //     id: 2,
  //     author: "Ayoub Benramdan",
  //     authorAvatar: "https://via.placeholder.com/40",
  //     title: "Another interesting post title",
  //     content:
  //       "Lorem ipsum dolor sit amet consectetur. Vel in molestie amet tellus dui. At arcu mauris purus magna. Volutpat dictum consectetur semper aliquam volutpat.",
  //     category: "Design",
  //     createdAt: "5 hours ago",
  //     likes: 15,
  //     comments: [],
  //     commentsCount: 3,
  //   },
  // ];

  if (!isLoggedIn) {
    return (
      <div className="auth-container">
        <div className="auth-card">
          <h2 className="auth-title">
            {isLogin ? "Login to your account" : "Create a new account"}
          </h2>

          {isLogin ? (
            <LoginForm onLoginSuccess={handleLoginSuccess} />
          ) : (
            <RegisterForm />
          )}

          <div className="auth-toggle">
            <button onClick={toggleForm} className="toggle-btn">
              {isLogin ? "Register a new account" : "Have an account? Login"}
            </button>
          </div>
        </div>
      </div>
    );
  }


    if (isLoggedIn && homeData) {
      return (
        <div className="app-container">
          <Navbar user={homeData.user} />
    
          <div className="main-content">
            <div className="grid-layout">
              {/* Left Sidebar */}
              <div className="left-column">
                <LeftSidebar users={homeData.not_following} bestcategories={homeData.best_categories}/>
              </div>
    
              {/* Main Content */}
              <div className="center-column">
                <PostComponent posts={homeData.posts} />
              </div>
    
              {/* Right Sidebar */}
              <div className="right-column">
                <ProfileCard user={homeData.user} />
                <TopGroups groups={homeData.other_groups} />
              </div>
            </div>
          </div>
        </div>
      );
    }


}

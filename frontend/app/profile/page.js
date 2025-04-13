"use client";
import { useState, useEffect } from "react";
// console.log("Welcome:");
async function handleProfile() {
  try {
    const response = await fetch("http://localhost:8404/profile", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: 1,
    });

    if (response.ok) {
      const data = await response.json();
      console.log("profile data: ", data);
    }
  } catch (error) {
    console.error("Error logging out:", error);
  }
}
export default function Profile() {
  useEffect(() => {
    handleProfile();
  }, []);

  return (
    <div>
      <h1>hey</h1>
    </div>
  );
}

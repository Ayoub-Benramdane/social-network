"use client";




export default async function RegisterHandler() {
  const response = await fetch("http://localhost:8404/register", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(registerFormInputs),
  });
  const data = await response.json();
  console.log(data);
}

// import Link from "next/link";

// export default function Home() {
//   return (
//     <div className="container">
//       <h1>Welcome to the Social Network</h1>

//       <div>
//         <h2>Login or Register</h2>
//         <p>
//           <Link href="/login">Login</Link>
//         </p>
//         <p>
//           <Link href="/register">Register</Link>
//         </p>
//       </div>
//     </div>
//   );
// }
import { useState } from 'react';

export default function Home() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errorMessage, setErrorMessage] = useState('');
  const [successMessage, setSuccessMessage] = useState('');

  const handleLogin = async (e) => {
    e.preventDefault();

    const res = await fetch('/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });

    const data = await res.json();

    if (res.ok) {
      localStorage.setItem('auth_token', data.token);
      setSuccessMessage(data.message);
      setErrorMessage('');
    } else {
      setErrorMessage(data.message);
      setSuccessMessage('');
    }
  };

  return (
    <div>
      <h1>Login</h1>
      <form onSubmit={handleLogin}>
        <div>
          <label htmlFor="email">Email:</label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div>
          <label htmlFor="password">Password:</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <button type="submit">Login</button>
      </form>
      {errorMessage && <p style={{ color: 'red' }}>{errorMessage}</p>}
      {successMessage && <p style={{ color: 'green' }}>{successMessage}</p>}
      <p>Don't have an account? <a href="/register">Register here</a></p>
    </div>
  );
}

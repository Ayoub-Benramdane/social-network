import { useState } from 'react';
// import styles from './Register.module.css';

const Register = () => {
  const [username, setUsername] = useState('');
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [email, setEmail] = useState('');
  const [dob, setDob] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [errorMessage, setErrorMessage] = useState('');
  const [successMessage, setSuccessMessage] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (password !== confirmPassword) {
      setErrorMessage('Passwords do not match.');
      setSuccessMessage('');
      return;
    }

    const res = await fetch('/api/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, firstName, lastName, email, dob, password, confirmPassword }),
    });

    const data = await res.json();
    if (res.ok) {
      setSuccessMessage(data.message);
      setErrorMessage('');
    } else {
      setErrorMessage(data.message);
      setSuccessMessage('');
    }
  };

  return (
    <div>
      <h1>Register Form</h1>
      <form id="registerForm" className="auth-form" onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="registerUsername">Username</label>
          <input
            type="text"
            id="registerUsername"
            name="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
          <small style={{ color: 'red', display: 'none' }}></small>
        </div>

        <div className="form-group">
          <label htmlFor="firstName">First Name</label>
          <input
            type="text"
            id="firstName"
            name="firstName"
            value={firstName}
            onChange={(e) => setFirstName(e.target.value)}
            required
          />
          <small style={{ color: 'red', display: 'none' }}></small>
        </div>

        <div className="form-group">
          <label htmlFor="lastName">Last Name</label>
          <input
            type="text"
            id="lastName"
            name="lastName"
            value={lastName}
            onChange={(e) => setLastName(e.target.value)}
            required
          />
          <small style={{ color: 'red', display: 'none' }}></small>
        </div>

        <div className="form-group">
          <label htmlFor="registerEmail">Email</label>
          <input
            type="email"
            id="registerEmail"
            name="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <small style={{ color: 'red', display: 'none' }}></small>
        </div>

        <div className="form-group">
          <label htmlFor="dob">Date of Birth</label>
          <input
            type="date"
            id="dob"
            name="dob"
            value={dob}
            onChange={(e) => setDob(e.target.value)}
            required
          />
          <small
            id="dob-error"
            style={{ color: 'red', display: dob && (new Date().getFullYear() - new Date(dob).getFullYear() < 18) ? 'block' : 'none' }}
          >
            You must be at least 18 years old.
          </small>
        </div>

        <div className="form-group">
          <label htmlFor="registerPassword">Password</label>
          <input
            type="password"
            id="registerPassword"
            name="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          <small style={{ color: 'red', display: 'none' }}></small>
        </div>

        <div className="form-group">
          <label htmlFor="confirmPassword">Confirm Password</label>
          <input
            type="password"
            id="confirmPassword"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            required
          />
          <small
            id="passwordError"
            style={{ color: 'red', display: password !== confirmPassword ? 'block' : 'none' }}
          >
            Passwords do not match.
          </small>
        </div>

        <button type="submit">Register</button>
      </form>

      {errorMessage && <p style={{ color: 'red' }}>{errorMessage}</p>}
      {successMessage && <p style={{ color: 'green' }}>{successMessage}</p>}
      <p>have I account? <a href="/">login here</a></p>
    </div>
  );
};

export default Register;

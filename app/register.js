// import { useState } from 'react';
// import styles from './page.module.css';

// const Register = () => {
//     const [isRegistering, setIsRegistering] = useState(false);
//     const [username, setUsername] = useState('');
//     const [firstName, setFirstName] = useState('');
//     const [lastName, setLastName] = useState('');
//     const [email, setEmail] = useState('');
//     const [dateOfBirth, setDob] = useState('');
//     const [password, setPassword] = useState('');
//     const [confirmPassword, setConfirmPassword] = useState('');

//     const [errorMessage, setErrorMessage] = useState('');
//     const [successMessage, setSuccessMessage] = useState('');

//     const handleSubmit = async (e) => {
//         e.preventDefault();
//         const res = await fetch('/api/register', {
//             method: 'POST',
//             headers: {
//                 'Content-Type': 'application/json',
//             },
//             body: JSON.stringify({ username, firstName, lastName, email, dateOfBirth, password, confirmPassword }),
//             mode: 'no-cors',
//         });

//         if (isRegistering) {
//             if (password !== confirmPassword) {
//                 setErrorMessage("Passwords do not match");
//                 return;
//             }
//             setSuccessMessage("Registration successful!");
//         } else {
//             setSuccessMessage("Login successful!");
//         }
//     };

//     return (
//         <div className={styles.container}>
//             <h1 className={styles.title}>Social Network</h1>

//             <div>
//                 <h2 className={styles.subtitle}>{isRegistering ? "Create New Account" : "Log into your account"}</h2>

//                 {isRegistering ? (
//                     <form className={styles.formContainer} onSubmit={handleSubmit}>
//                         <div>
//                             <label htmlFor="username">Username</label>
//                             <input
//                                 className={styles.inputField}
//                                 type="text"
//                                 id="username"
//                                 name="username"
//                                 value={username}
//                                 onChange={(e) => setUsername(e.target.value)}
//                                 required
//                             />
//                         </div>
//                         <div>
//                             <label htmlFor="firstName">First Name</label>
//                             <input
//                                 className={styles.inputField}
//                                 type="text"
//                                 id="firstName"
//                                 name="firstName"
//                                 value={firstName}
//                                 onChange={(e) => setFirstName(e.target.value)}
//                                 required
//                             />
//                         </div>
//                         <div>
//                             <label htmlFor="lastName">Last Name</label>
//                             <input
//                                 className={styles.inputField}
//                                 type="text"
//                                 id="lastName"
//                                 name="lastName"
//                                 value={lastName}
//                                 onChange={(e) => setLastName(e.target.value)}
//                                 required
//                             />
//                         </div>
//                         <div>
//                             <label htmlFor="email">Email</label>
//                             <input
//                                 className={styles.inputField}
//                                 type="email"
//                                 id="email"
//                                 name="email"
//                                 value={email}
//                                 onChange={(e) => setEmail(e.target.value)}
//                                 required
//                             />
//                         </div>
//                         <div>
//                             <label htmlFor="dob">Date of Birth</label>
//                             <input
//                                 className={styles.inputField}
//                                 type="date"
//                                 id="dob"
//                                 name="dob"
//                                 value={dateOfBirth}
//                                 onChange={(e) => setDob(e.target.value)}
//                                 required
//                             />
//                         </div>
//                         <div>
//                             <label htmlFor="password">Password</label>
//                             <input
//                                 className={styles.inputField}
//                                 type="password"
//                                 id="password"
//                                 name="password"
//                                 value={password}
//                                 onChange={(e) => setPassword(e.target.value)}
//                                 required
//                             />
//                         </div>
//                         <div>
//                             <label htmlFor="confirmPassword">Confirm Password</label>
//                             <input
//                                 className={styles.inputField}
//                                 type="password"
//                                 id="confirmPassword"
//                                 name="confirmPassword"
//                                 value={confirmPassword}
//                                 onChange={(e) => setConfirmPassword(e.target.value)}
//                                 required
//                             />
//                         </div>
//                         <button className={styles.button} type="submit">Sign Up</button>
//                         {errorMessage && <p className={styles.error}>{errorMessage}</p>}
//                         {successMessage && <p className={styles.success}>{successMessage}</p>}
//                     </form>
//                 ) : (
//                     <form className={styles.formContainer} onSubmit={handleSubmit}>
//                         <div>
//                             <label htmlFor="email">Email</label>
//                             <input
//                                 className={styles.inputField}
//                                 type="email"
//                                 id="email"
//                                 name="email"
//                                 value={email}
//                                 onChange={(e) => setEmail(e.target.value)}
//                                 required
//                             />
//                         </div>
//                         <div>
//                             <label htmlFor="password">Password</label>
//                             <input
//                                 className={styles.inputField}
//                                 type="password"
//                                 id="password"
//                                 name="password"
//                                 value={password}
//                                 onChange={(e) => setPassword(e.target.value)}
//                                 required
//                             />
//                         </div>
//                         <button className={styles.button} type="submit">Log In</button>
//                         {errorMessage && <p className={styles.error}>{errorMessage}</p>}
//                         {successMessage && <p className={styles.success}>{successMessage}</p>}
//                     </form>
//                 )}

//                 <p className={styles.p}>
//                     {isRegistering ? (
//                         <>
//                             Already have an account?{" "}
//                             <button className={styles.button} onClick={() => setIsRegistering(false)}>Log In</button>
//                         </>
//                     ) : (
//                         <>
//                             Don't have an account?{" "}
//                             <button className={styles.button} onClick={() => setIsRegistering(true)}>Sign Up</button>
//                         </>
//                     )}
//                 </p>
//             </div>
//         </div>
//     );
// };

// export default Register;


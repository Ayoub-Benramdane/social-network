export default async function handler(req, res) {
  res.setHeader('Access-Control-Allow-Origin', 'http://localhost:3000');
  res.setHeader('Access-Control-Allow-Methods', 'POST');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization');
  
  if (req.method === 'OPTIONS') {
    return res.status(200).end();
  }

  if (req.method === 'POST') {
    const { email, password } = req.body;

    if (!email || !password) {
      return res.status(400).json({ message: 'Email et mot de passe sont requis!' });
    }

    try {
      const response = await fetch('http://localhost:8404/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
        mode: 'cors', 
      });

      const data = await response.json();
      console.log('API response:', data);

      if (response.ok) {
        return res.status(200).json({
          message: 'Login successful!',
          token: data.token,  
        });
      } else {
        return res.status(401).json({ message: 'Invalid email or password' });
      }

    } catch (error) {
      console.error('Error during login:', error);
      return res.status(500).json({ message: 'Error login' });
    }
  } else {
    return res.status(405).json({ message: 'Method Not Allowed' });
  }
}

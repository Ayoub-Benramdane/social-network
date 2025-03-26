export default async function handler(req, res) {  
  if (req.method === 'POST') {
    const { username, firstName, lastName, email, dob, password, confirmPassword } = req.body;
    // Check if all required fields are present
    if (!username || !firstName || !lastName || !email || !dob || !password || !confirmPassword) {
      return res.status(400).json({ message: 'Tous les champs sont obligatoires!' });
    }

    // Check if passwords match
    if (password !== confirmPassword) {
      return res.status(400).json({ message: 'Les mots de passe ne correspondent pas!' });
    }

    // Check if the user is at least 18 years old
    const age = new Date().getFullYear() - new Date(dob).getFullYear();
    if (age < 18) {
      return res.status(400).json({ message: 'Vous devez avoir au moins 18 ans.' });
    }

    try {
      const DateOfBirth = new Date(dob).toISOString();
      // Send the registration request to the external API
      const response = await fetch('http://localhost:8404/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, firstName, lastName, email, DateOfBirth, password, confirmPassword }),
        mode: 'no-cors',
      });

      const data = await response.json();

      if (response.ok) {
        return res.status(201).json({ message: 'Utilisateur créé avec succès' });
      } else {
        return res.status(400).json({ message: 'Échec de l\'inscription' });
      }
    } catch (error) {
      console.error('Error during registration:', error);
      return res.status(500).json({ message: 'Erreur lors de l\'inscription' });
    }
  } else {
    return res.status(405).json({ message: 'Method Not Allowed' });
  }
}

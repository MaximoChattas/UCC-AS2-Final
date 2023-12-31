import React, { useState, useEffect, useContext } from 'react';
import { useNavigate } from 'react-router-dom';
import jwt_decode from 'jwt-decode';
import Navbar from '../NavBar/NavBar';
import './Login.css'
import { UserProfileContext } from '../../App';

function Login() {
  const [loading, setLoading] = useState(false);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const { userProfile, setUserProfile } = useContext(UserProfileContext);

  useEffect(() => {
    const token = localStorage.getItem('token');

    if (token) {
      setUserProfile(jwt_decode(token));
    }
  }, []);
  
  const handleLogin = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
  
    try {
      const response = await fetch('http://localhost:8090/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });
  
      if (response.status === 202) {
        const { token } = await response.json();
        localStorage.setItem('token', token);
        console.log(token)
        setUserProfile(jwt_decode(token));
        navigate('/');
      } else {
        const data = await response.json();
        const errorMessage = data.error || 'Error';
        throw new Error(errorMessage);
      }
    } catch (error) {
      console.error(error);
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <UserProfileContext.Provider value={{ userProfile, setUserProfile }}>
      <>
        <Navbar />
        <div className="contenedorLogin">
          <div>
              <h2>Inicio de Sesion</h2>
              <form onSubmit={handleLogin}>
                <div>
                  <label>Email:</label>
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                </div>
                <div>
                  <label>Clave:</label>
                  <input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                  />
                </div>
                {error && <p className="error-message">{error}</p>}
                <button type="submit" disabled={loading}>
                  {loading ? 'Cargando...' : 'Iniciar Sesion'}
                </button>
              </form>
              </div>
              <div>
              <p>¿Aun no tienes una cuenta?</p>
              <button onClick={()=>navigate('/signup')}>Registrate</button>
            </div>
        </div>
      </>
    </UserProfileContext.Provider>
  );
}

export default Login;

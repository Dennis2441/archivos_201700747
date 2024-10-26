// frontend/src/Login.jsx  
import { useState } from 'react';
// import '../Login.css';

export function Login() {  
  const [idParticion, setIdParticion] = useState('');  
  const [usuario, setUsuario] = useState('');  
  const [password, setPassword] = useState('');  
  const [isAuthenticated, setIsAuthenticated] = useState(false);  
  const [message, setMessage] = useState('');  

  const handleLogin = async () => {  
    try {  
      const response = await fetch('http://localhost:8080/login', {  
        method: 'POST',  
        headers: {  
          'Content-Type': 'application/json',  
        },  
        body: JSON.stringify({ idParticion, usuario, password }),  
      });  

      if (!response.ok) {  
        throw new Error('Error en la solicitud de login');  
      }  

      const data = await response.json();  
      setIsAuthenticated(data.success);  

      if (data.success) {  
        setMessage('Login exitoso');  
      } else {  
        setMessage('Credenciales incorrectas');  
      }  
    } catch (error) {  
      console.error('Error:', error);  
      setMessage('Error al procesar la solicitud de login');  
    }  
  };  

  if (isAuthenticated) {  
    return <div className="login-message">Bienvenido, {usuario}!</div>;  
  }  

  return (  
    <div className="login-container">  
      <div className="login-box">  
        <h2 className="login-title">Formulario de Login</h2>  
        <input  
          type="text"  
          placeholder="ID Partición"  
          value={idParticion}  
          onChange={(e) => setIdParticion(e.target.value)}  
          className="login-input"  
        />  
        
        <input  
          type="text"  
          placeholder="Usuario"  
          value={usuario}  
          onChange={(e) => setUsuario(e.target.value)}  
          className="login-input"  
        />  
        <input  
          type="password"  
          placeholder="Password"  
          value={password}  
          onChange={(e) => setPassword(e.target.value)}  
          className="login-input"  
        />  
        <button className="login-button" onClick={handleLogin}>  
          Enviar  
        </button>  
        <p className="login-message">{message}</p>  
      </div>  
    </div>  
  ); 
} 
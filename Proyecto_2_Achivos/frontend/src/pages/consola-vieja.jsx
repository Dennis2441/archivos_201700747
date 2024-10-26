// frontend/src/App.jsx  
import { useState } from 'react';  
import '../App.css';  

function ConsolaVieja() {  
  const [input, setInput] = useState('');  
  const [output, setOutput] = useState('');  
  const [fileName, setFileName] = useState('');  

  const handleSubmit = async () => {  
    if (input) {  
      try {  
        const response = await fetch('http://localhost:8080/submit', {  
          method: 'POST',  
          headers: {  
            'Content-Type': 'application/json',  
          },  
          body: JSON.stringify({ input }),  
        });  

        if (!response.ok) {  
          throw new Error('Error en la solicitud');  
        }  

        const data = await response.json();  
        setOutput(data.output);  
      } catch (error) {  
        console.error('Error:', error);  
        setOutput('Error al procesar la solicitud');  
      }  
    } else {  
      setOutput('Salida: No hay entrada.');  
    }  
    setInput('');  
  };  

  const handleFileChange = (event) => {  
    const file = event.target.files[0];  
    if (file) {  
      setFileName(file.name);  
      const reader = new FileReader();  
      reader.onload = (e) => {  
        const fileContent = e.target.result;  
        setInput(fileContent);  
      };  
      reader.readAsText(file);  
    }  
  };  

  return (  
    <div className="container">  
      <h1>Aplicación de Entrada y Salida</h1>  
      <div className="input-area">  
        <textarea  
          className="input-textarea"  
          value={input}  
          onChange={(e) => setInput(e.target.value)}  
          placeholder="Escribe tu entrada aquí..."  
        />  
        <input  
          type="file"  
          onChange={handleFileChange}  
          className="file-input"  
          accept=".txt, .sh, .bash"  
        />  
        {fileName && <p className="file-name">Archivo seleccionado: {fileName}</p>}  
        <button className="submit-button" onClick={handleSubmit}>  
          Enviar  
        </button>  
      </div>  
      <div className="output-area">  
        <p>{output}</p>  
      </div>  
    </div>  
  );  
}  

export default ConsolaVieja;
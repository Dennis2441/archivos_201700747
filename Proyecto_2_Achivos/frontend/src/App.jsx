// Librerias
import { Routes, BrowserRouter as Router, Route} from "react-router-dom";
import { Explorador } from "./pages/explorador";
import { Login } from "./pages/login";
// import './App.css';

function App() {  
  return (  
    <Router>  
      <Routes>
        <Route path="/" element={<Explorador />} />  
        <Route path="/login" element={<Login />} />  
      </Routes>  
    </Router>  
  );  
}  

export default App;
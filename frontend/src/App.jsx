import { BrowserRouter, Routes, Route } from "react-router-dom";
import Dashboard from "./components/Dashboard";
import Navbar from "./components/Navbar";
import "./App.css";
import Login from "./components/Login";

function App() {
  return (
    <BrowserRouter>
      {/* We keep Navbar at the top because we want it to show on every page */}
      <Navbar/>
      {/** The routes allow us to say: based on this path, effect this prop (element) */}
      <Routes>
        <Route path="/" element={<Dashboard />}></Route>
        <Route path="/login" element={<Login />}></Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;

import { BrowserRouter, Routes, Route } from "react-router-dom";
import { useState, useEffect } from "react";
import Dashboard from "./components/Dashboard";
import Navbar from "./components/Navbar";
import "./App.css";
import Login from "./components/Login";

function App() {
  // check t o see if we're logged in when going on the website
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [username, setUsername] = useState("");

  useEffect(() => {
    const token = localStorage.getItem("jwt_token");
    const savedUsername = localStorage.getItem("username");
    if (token) {
      setIsLoggedIn(true);
      setUsername(savedUsername || "");
      console.log("User already logged in:", savedUsername);
    }
  }, []);

  return (
    <BrowserRouter>
      {/* We keep Navbar at the top because we want it to show on every page */}
      <Navbar isLoggedIn={isLoggedIn} username={username}/>
      {/** The routes allow us to say: based on this path, effect this prop (element) */}
      <Routes>
        <Route path="/" element={<Dashboard />}></Route>
        <Route path="/login" element={<Login />}></Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;

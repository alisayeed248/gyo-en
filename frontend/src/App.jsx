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

  const handleLogout = () => {
    localStorage.removeItem("jwt_token")
    localStorage.removeItem("username");
    setIsLoggedIn(false);
    setUsername("");
    window.location.href = "/login";
  }

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
      <Navbar isLoggedIn={isLoggedIn} username={username} onLogout={handleLogout}/>
      {/** The routes allow us to say: based on this path, effect this prop (element) */}
      <div className="pt-20">
        <Routes>
          <Route
            path="/"
            element={
              isLoggedIn ? (
                <Dashboard />
              ) : (
                <Login
                  setIsLoggedIn={setIsLoggedIn}
                  setUsername={setUsername}
                />
              )
            }
          ></Route>
          <Route
            path="/login"
            element={
              <Login setIsLoggedIn={setIsLoggedIn} setUsername={setUsername} />
            }
          />
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App;

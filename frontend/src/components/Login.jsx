import { useState } from "react";

function Login({ setIsLoggedIn, setUsername: setAppUsername }) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  // not just an async function, we need to send the event?
  // regular synchronous function, not arrow btw
  // e here is the form submission event object
  // handleSubmit won't be on the button level but on the form level
  const handleSubmit = async (e) => {
    e.preventDefault();
    console.log("Form submitted!");
    console.log("Username:", username);
    console.log("Password:", password);

    try {
      const response = await fetch("http://localhost:8080/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      if (response.ok) {
        const data = await response.json();
        console.log("Login successful! Token: ", data.token);

        localStorage.setItem("jwt_token", data.token);
        localStorage.setItem("username", data.username || username);

        setIsLoggedIn(true);
        setAppUsername(username);

        window.location.href = "/";

        alert("Login successful! Check localStorage")
      } else {
        const errorText = await response.text()
        console.log("Login failed:", errorText);
        alert("Login failed: " + errorText);
      }
     
    } catch (error) {
      console.log("Error:", error);
      alert("Network error occurred");
    }
  };

  return (
    <div>
      <div className="h-screen flex items-center justify-center">
        <div>
          <form onSubmit={handleSubmit}>
            <div>Username</div>
            <input
              type="text"
              className="mb-4 text-black"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            ></input>

            <div>Password</div>
            <input
              type="password"
              className="mb-4 text-black"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            ></input>

            <div>
              <button type="submit">Login</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}

export default Login;

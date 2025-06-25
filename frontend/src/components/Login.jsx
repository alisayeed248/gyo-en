import { useState } from "react";

function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  // not just an async function, we need to send the event?
  // regular synchronous function, not arrow btw
  // e here is the form submission event object
  // handleSubmit won't be on the button level but on the form level
  const handleSubmit = (e) => {
    e.preventDefault();
    console.log("Form submitted!");
    console.log("Username:", username);
    console.log("Password:", password);
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

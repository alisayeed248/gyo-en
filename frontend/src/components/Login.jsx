import { useState } from "react";

function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  return (
    <div>
      <div className="h-screen flex items-center justify-center">
        <div>
          <form>
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

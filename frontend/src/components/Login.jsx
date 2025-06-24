function Login() {
  return (
    <div>
      <div className="h-screen flex items-center justify-center">
        <div>
          <form>
            <div>Username</div>
            <input type="text" className="mb-4"></input>

            <div>Password</div>
            <input type="password" className="mb-4"></input>

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

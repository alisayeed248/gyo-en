function Navbar({ isLoggedIn, username }) {
  const handleLoginClick = () => {
    window.location.href = "/login";
  };

  const handleLogout = () => {
    localStorage.removeItem("jwt_token");
    localStorage.removeItem("username");
    window.location.href = "/login";
  };

  return (
    <div className="fixed top-0 left-0 w-full bg-red-500 flex justify-between items-center p-4 z-10">
      <div>gyo-en Monitor</div>
      <div>
        {isLoggedIn ? (
          <>
            <span className="mr-4">Hello, {username}</span>
            <button
              onClick={handleLoginClick}
              className="mr-2 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Logout
            </button>
          </>
        ) : (
          <>
            <button
              onClick={handleLoginClick}
              className="mr-2 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Login
            </button>
          </>
        )}
      </div>
    </div>
  );
}

export default Navbar;

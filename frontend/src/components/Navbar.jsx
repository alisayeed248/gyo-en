function Navbar() {
  const handleLoginClick = () => {
    console.log("Login clicked.")
  }

  return (
    <nav>
      <div className="fixed top-0 w-full bg-red-500 flex justify-between items-center p-4" style={{right: 0}}>
        <div>gyo-en Monitor</div>
        <div>
          <button onClick={handleLoginClick} className="mr-2 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">Login</button>
          <button>Register</button>
        </div>
      </div>
    </nav>
  );
}

export default Navbar;

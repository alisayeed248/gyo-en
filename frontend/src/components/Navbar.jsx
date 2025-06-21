function Navbar() {
  return (
    <nav>
      <div className="fixed top-0 w-full bg-red-500 flex justify-between items-center p-4" style={{right: 0}}>
        <div>gyo-en Monitor</div>
        <div>
          <button>Login</button>
          <button>Register</button>
        </div>
      </div>
    </nav>
  );
}

export default Navbar;

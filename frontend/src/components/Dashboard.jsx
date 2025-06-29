import { useState, useEffect } from "react";
import Navbar from "./Navbar";

function Dashboard() {
  const [urlStatus, setUrlStatus] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchStatus = async () => {
      try {
        const token = localStorage.getItem("jwt_token");
        console.log("Token from localStorage:", token);
        const response = await fetch("http://localhost:8080/api/status", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        console.log("Response status:", response.status);
        if (!response.ok) {
          const errorText = await response.text();
          console.log("Error response:", errorText); // Debug line
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        console.log("Response data:", data);
        setUrlStatus(data.urls);
        setLoading(false);
      } catch (err) {
        setError("Failed to fetch status");
        setLoading(false);
      }
    };

    fetchStatus();
    // Refresh every 10 seconds
    const interval = setInterval(fetchStatus, 10000);

    return () => clearInterval(interval);
  }, []);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>gyo-en Monitor</h1>
      <p>Last updated: {new Date().toLocaleTimeString()}</p>

      {urlStatus.map((site) => (
        <div
          key={site.url}
          style={{
            padding: "10px",
            margin: "10px 0",
            border: "1px solid #ccc",
            borderRadius: "4px",
            backgroundColor: site.status === "UP" ? "#d4edda" : "#f8d7da",
          }}
        >
          <strong>{site.url}</strong>
          <span
            style={{
              float: "right",
              color: site.status === "UP" ? "green" : "red",
            }}
          >
            {site.status}
          </span>
          <br />
          <small>Last check: {site.lastCheck || "Never"}</small>
        </div>
      ))}
    </div>
  );
}

export default Dashboard;

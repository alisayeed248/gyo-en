import { useState, useEffect } from 'react';

function Dashboard() {
  const [urlStatus, setUrlStatus] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    const fetchStatus = async () => {
      try {
        const response = await fetch('/api/status')
        const data = await response.json()
        setUrlStatus(data.urls)
        setLoading(false)
      } catch (err) {
        setError('Failed to fetch status')
        setLoading(false)
      }
    }

    fetchStatus()
    // Refresh every 10 seconds
    const interval = setInterval(fetchStatus, 10000)
     
    return () => clearInterval(interval)
  }, [])

  if (loading) return <div>Loading...</div>
  if (error) return <div>Error: {error}</div>

  return (
    <div>
      <h1>gyo-en Monitor</h1>
      <p>Last updated: {new Date().toLocaleTimeString()}</p>
    </div>
  )
}
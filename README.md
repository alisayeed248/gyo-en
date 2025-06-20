# gyo-en

Full-stack uptime monitoring platform built with Go and React. Monitor your websites and APIs with real-time dashboard, database persistence, and live status updates.

## Features

- 🔍 **Real-time URL monitoring** - Check website status every 30 seconds
- 📊 **Live React dashboard** - Auto-refreshing status with response times
- 💾 **SQLite database** - Persistent storage for users, URLs, and check history
- ☸️ **Kubernetes-ready** - Designed for cloud deployment on GCP
- 🏠 **Local development friendly** - Works without external dependencies
- 🚨 **Status change detection** - Track UP/DOWN transitions
- 📈 **Response time tracking** - Monitor site performance over time

## Quick Start

### Prerequisites
- Go 1.22+
- Node.js 18+

### Local Development (Fastest)

**Terminal 1: Start Backend**
```bash
go run cmd/gyo-en/main.go
```

**Terminal 2: Start Frontend**
```bash
cd frontend
npm install
npm run dev
```

**Visit:** http://localhost:5173 to see the live dashboard!

**API:** http://localhost:8080/api/status for JSON data

## Architecture

```
├── cmd/gyo-en/main.go          # Backend entry point
├── backend/
│   └── database/
│       ├── models.go           # Data models (User, MonitoredURL, CheckResult)
│       └── database.go         # SQLite connection & GORM setup
├── internal/monitor/checker.go # URL checking logic
├── frontend/                   # React dashboard
│   ├── src/components/Dashboard.jsx
│   └── package.json
├── test-urls.txt              # URLs for local development
├── gyo-en.db                  # SQLite database (auto-created)
└── k8s-*.yaml                 # Kubernetes deployment files
```

## Data Models

**User** - Authentication and user management
```go
type User struct {
    ID       uint
    Username string
    Email    string
    Password string
}
```

**MonitoredURL** - URLs each user wants to monitor
```go
type MonitoredURL struct {
    ID     uint
    UserID uint
    URL    string
    Name   string  // Optional friendly name
}
```

**CheckResult** - Historical monitoring data
```go
type CheckResult struct {
    URL          string
    IsUp         bool
    ResponseTime time.Duration
    CheckedAt    time.Time
}
```

## API Endpoints

- `GET /api/status` - Current status of all monitored URLs
- `GET /health` - Health check endpoint

## Configuration

### Environment Variables
- `ENVIRONMENT` - development/production (default: development)
- `PORT` - HTTP server port (default: 8080)
- `DB_PATH` - SQLite database file path (default: gyo-en.db)
- `URLS_FILE` - File containing URLs to monitor (default: test-urls.txt)

### Monitored URLs
**Local development:** Edit `test-urls.txt`
```
https://www.google.com
https://github.com
https://your-website.com
```

## Development Workflow

### Backend Changes
```bash
# Test locally
go run cmd/gyo-en/main.go

# Add dependencies
go get package-name
go mod tidy
```

### Frontend Changes
```bash
cd frontend
npm run dev  # Hot reload for development
npm run build  # Build for production
```

### Database
- **SQLite file:** `gyo-en.db` (auto-created)
- **Migrations:** Automatic via GORM AutoMigrate
- **Reset database:** Delete `gyo-en.db` file

## Kubernetes Deployment

### Prerequisites
- Docker Desktop
- minikube or GKE cluster

### Deploy to Kubernetes
```bash
# Build and load image
docker build -t gyo-en:latest .
minikube image load gyo-en:latest

# Deploy
kubectl apply -f k8s-configmap.yaml
kubectl apply -f k8s-deployment.yaml

# Monitor
kubectl logs -l app=gyo-en -f
```

## How It Works

1. **Backend** monitors URLs from `test-urls.txt` every 30 seconds
2. **Check results** are stored in SQLite database with timestamps
3. **API** serves current status and historical data
4. **Frontend** fetches data every 10 seconds and displays live dashboard
5. **Database** persists all data across restarts

## Technology Stack

**Backend:**
- Go 1.22 with net/http
- GORM ORM with SQLite
- JSON REST API

**Frontend:**
- React 19 with Vite
- JavaScript (no TypeScript yet)
- CSS for styling

**Database:**
- SQLite for development
- PostgreSQL planned for production (GCP Cloud SQL)

**Infrastructure:**
- Docker containers
- Kubernetes deployment
- GCP free tier ready

## Roadmap

### Phase 1: User Authentication (In Progress)
- [ ] User registration/login
- [ ] Session management
- [ ] Protected API endpoints

### Phase 2: Personal Monitoring
- [ ] Users can add/remove their own URLs
- [ ] Personal dashboards
- [ ] User-specific check history

### Phase 3: Notifications
- [ ] Email alerts when sites go down
- [ ] Notification preferences
- [ ] Alert history

### Phase 4: Production Deployment
- [ ] Deploy to GCP GKE
- [ ] PostgreSQL database
- [ ] CI/CD pipeline
- [ ] Custom domain with HTTPS

## Contributing

This is a learning project focused on Kubernetes and Go development. The goal is to build a real monitoring service while learning cloud-native technologies.

## License

MIT License - See LICENSE file for details
# gyo-en

Full-stack uptime monitoring platform built with Go and React. Monitor your websites and APIs with real-time dashboard, JWT authentication, database persistence, and live status updates.

## Features

- 🔍 **Real-time URL monitoring** - Check website status every 30 seconds
- 🔐 **JWT Authentication** - Secure login system with bcrypt password hashing
- 📊 **Live React dashboard** - Auto-refreshing status with response times
- 💾 **SQLite database** - Persistent storage for users, URLs, and check history
- 🛡️ **Protected routes** - Dashboard requires authentication
- 👤 **User management** - Login/logout with session persistence
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

**Test Login:**
- Username: `test`
- Password: `password123`

**API:** http://localhost:8080/api/status for JSON data

## Authentication System

### JWT-Based Authentication
- **Secure login** with bcrypt password hashing
- **JWT tokens** with 1-hour expiration
- **Protected routes** - dashboard requires authentication
- **Session persistence** - stays logged in across browser sessions
- **Automatic logout** when tokens expire

### Test User
A test user is automatically created on first run:
- **Username**: `test`
- **Password**: `password123`

### API Endpoints
- `POST /api/login` - Authenticate user and receive JWT token
- `GET /api/status` - Current status of all monitored URLs (protected)
- `GET /health` - Health check endpoint

## Architecture

```
├── cmd/gyo-en/main.go          # Backend entry point with auth setup
├── internal/
│   ├── auth/
│   │   ├── handlers.go         # Login endpoint and JWT generation
│   │   └── validation.go       # User validation and password checking
│   ├── database/
│   │   ├── models.go           # Data models (User, MonitoredURL, CheckResult)
│   │   └── database.go         # SQLite connection & GORM setup
│   └── monitor/checker.go      # URL checking logic
├── frontend/                   # React dashboard with authentication
│   ├── src/components/
│   │   ├── Dashboard.jsx       # Protected dashboard component
│   │   ├── Login.jsx          # JWT-enabled login form
│   │   └── Navbar.jsx         # Dynamic navbar with user status
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
    Username string `gorm:"unique;not null"`
    Email    string `gorm:"unique;not null"`
    Password string `gorm:"not null"`  // bcrypt hashed
    CreatedAt time.Time
}
```

**MonitoredURL** - URLs each user wants to monitor
```go
type MonitoredURL struct {
    ID     uint
    UserID uint     // Links to User.ID
    URL    string
    Name   string   // Optional friendly name
}
```

**CheckResult** - Historical monitoring data
```go
type CheckResult struct {
    URL          string
    IsUp         bool
    ResponseTime time.Duration
    StatusCode   int
    CheckedAt    time.Time
}
```

## Authentication Flow

1. **User logs in** → Frontend sends credentials to `/api/login`
2. **Backend validates** → Checks username/password against database
3. **JWT generation** → Creates signed token with user info
4. **Frontend storage** → Stores JWT in localStorage
5. **Protected routes** → Dashboard requires valid JWT
6. **Auto-login** → Checks localStorage on app startup
7. **Logout** → Clears localStorage and redirects

## Security Features

- **Password hashing** with bcrypt (cost factor 10)
- **JWT tokens** with HMAC-SHA256 signing
- **Token expiration** (1 hour default)
- **Protected API endpoints** require valid JWT
- **CORS handling** for cross-origin requests
- **Input validation** and error handling

## Configuration

### Environment Variables
- `ENVIRONMENT` - development/production (default: development)
- `PORT` - HTTP server port (default: 8080)
- `DB_PATH` - SQLite database file path (default: gyo-en.db)
- `URLS_FILE` - File containing URLs to monitor (default: test-urls.txt)
- `JWT_SECRET` - Secret key for JWT signing (default: hardcoded for dev)

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
- **Test user:** Automatically created on first run

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

1. **User authentication** via JWT tokens with secure password validation
2. **Backend** monitors URLs from `test-urls.txt` every 30 seconds
3. **Check results** stored in SQLite database with timestamps
4. **Protected API** serves current status and historical data (requires JWT)
5. **Frontend** fetches data every 10 seconds and displays live dashboard
6. **Session persistence** keeps users logged in across browser sessions
7. **Database** persists all data across restarts

## Technology Stack

**Backend:**
- Go 1.22 with net/http
- JWT authentication with golang-jwt/jwt/v5
- GORM ORM with SQLite
- bcrypt for password hashing
- JSON REST API with CORS support

**Frontend:**
- React 19 with Vite
- React Router for navigation
- localStorage for JWT persistence
- Tailwind CSS for styling
- JavaScript (no TypeScript yet)

**Database:**
- SQLite for development
- PostgreSQL planned for production (GCP Cloud SQL)

**Infrastructure:**
- Docker containers
- Kubernetes deployment
- GCP free tier ready

## Roadmap

### Phase 1: Enhanced Authentication ✅
- [x] User registration/login with JWT
- [x] Session management and persistence
- [x] Protected API endpoints
- [x] Secure password hashing

### Phase 2: Personal Monitoring (In Progress)
- [ ] Users can add/remove their own URLs
- [ ] Personal dashboards per user
- [ ] User-specific check history
- [ ] JWT middleware for API protection

### Phase 3: Advanced Features
- [ ] Email alerts when sites go down
- [ ] Notification preferences
- [ ] Alert history and escalation
- [ ] Multi-user organizations

### Phase 4: Production Deployment
- [ ] Deploy to GCP GKE
- [ ] PostgreSQL database
- [ ] CI/CD pipeline with GitHub Actions
- [ ] Custom domain with HTTPS
- [ ] Environment-based JWT secrets

## Contributing

This is a learning project focused on full-stack development, authentication systems, and Kubernetes deployment. The goal is to build a real monitoring service while learning modern web technologies.

## License

MIT License - See LICENSE file for details
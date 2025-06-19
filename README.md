# gyo-en

Kubernetes-based uptime monitoring platform with real-time dashboard. Monitor your websites and APIs with live status updates, built with Go backend and React frontend.

## Features

- 🔍 **Real-time URL monitoring** - Check website status every 30 seconds
- 📊 **Live dashboard** - React frontend with auto-refreshing status
- ☸️ **Kubernetes-native** - Designed for cloud deployment
- 🏠 **Local development friendly** - Works without Redis/K8s for development
- 🚨 **Status change detection** - Track UP/DOWN transitions
- 📈 **Response time tracking** - Monitor site performance

## Quick Start

### Local Development (Fastest)

```bash
# Terminal 1: Start backend
go run cmd/gyo-en/main.go

# Terminal 2: Start frontend (in another terminal)
cd frontend
npm install
npm run dev
```

Visit http://localhost:5173 to see the dashboard!

### Kubernetes Deployment

#### Prerequisites
- Docker Desktop running
- minikube installed

#### Setup
```bash
# Start Kubernetes cluster
minikube start

# Build and load image
docker build -t gyo-en:latest .
minikube image load gyo-en:latest

# Deploy to Kubernetes
kubectl apply -f k8s-configmap.yaml
kubectl apply -f k8s-deployment.yaml

# View logs
kubectl logs -l app=gyo-en -f
```

## Development Workflow

### Backend Changes
```bash
# Local testing
go run cmd/gyo-en/main.go

# For Kubernetes deployment after changes
docker build -t gyo-en:v6 .
minikube image load gyo-en:v6
# Update image tag in k8s-deployment.yaml
kubectl apply -f k8s-deployment.yaml
```

### Frontend Changes
```bash
cd frontend
npm run dev  # Hot reload for development
npm run build  # Build for production
```

### Configuration

#### Environment Variables
- `REDIS_ADDR` - Redis connection (default: localhost:6379)
- `PORT` - HTTP server port (default: 8080)
- `URLS_FILE` - File containing URLs to monitor (default: test-urls.txt)
- `ENVIRONMENT` - development/production (default: development)

#### Monitored URLs
**Local development**: Edit `test-urls.txt`
**Kubernetes**: Edit ConfigMap
```bash
kubectl edit configmap gyo-en-config
kubectl rollout restart deployment gyo-en-deployment
```

## Architecture

```
├── cmd/gyo-en/main.go          # Backend entry point
├── internal/monitor/checker.go  # URL checking logic
├── frontend/                   # React dashboard
│   ├── src/components/Dashboard.jsx
│   └── package.json
├── k8s-deployment.yaml         # Kubernetes deployment
├── k8s-configmap.yaml          # Configuration (URLs to monitor)
├── Dockerfile                  # Container build
└── test-urls.txt              # URLs for local development
```

## API Endpoints

- `GET /api/status` - Current status of all monitored URLs
- `GET /health` - Health check endpoint

## Monitoring Commands

```bash
# View running pods
kubectl get pods

# Stream logs
kubectl logs -l app=gyo-en -f

# Check pod details
kubectl describe pod -l app=gyo-en

# Force restart deployment
kubectl rollout restart deployment gyo-en-deployment
```

## How It Works

1. **Backend** reads URLs from file or ConfigMap
2. **Monitor loop** checks each URL every 30 seconds
3. **Redis** stores check history and detects status changes
4. **API** serves current status to frontend
5. **Dashboard** displays real-time status with auto-refresh
6. **Kubernetes** ensures service stays running in production

## Next Steps

- [ ] Add user authentication
- [ ] Allow users to add/remove URLs
- [ ] Email/Slack notifications
- [ ] Deploy to GKE
- [ ] Uptime percentage calculations
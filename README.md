# gyo-en

Kubernetes-based uptime monitoring service written in Go.

## Quick Start

### Prerequisites
- Docker Desktop running
- minikube installed

### Setup
```bash
# Start Kubernetes cluster
minikube start

# Build and load image
docker build -t gyo-en:v4 .
minikube image load gyo-en:v4

# Deploy to Kubernetes
kubectl apply -f k8s-configmap.yaml
kubectl apply -f k8s-deployment.yaml

# View logs
kubectl logs -l app=gyo-en -f
```

## Key Commands

### Build & Deploy
```bash
# Rebuild after code changes
docker build -t gyo-en:v5 .
minikube image load gyo-en:v5
# Update k8s-deployment.yaml image version
kubectl apply -f k8s-deployment.yaml

# Force restart
kubectl rollout restart deployment gyo-en-deployment
```

### Monitoring
```bash
kubectl get pods
kubectl logs -l app=gyo-en -f
kubectl describe pod -l app=gyo-en
```

### Configuration
```bash
# Edit monitored URLs
kubectl edit configmap gyo-en-config
# Then restart deployment to pick up changes
```

## Directory Structure

```
├── cmd/gyo-en/main.go          # Main application
├── internal/monitor/checker.go  # URL checking logic
├── k8s-deployment.yaml         # Kubernetes deployment
├── k8s-configmap.yaml          # Configuration (URLs to monitor)
├── Dockerfile                  # Container build
└── go.mod                      # Go module
```

## How It Works

1. **ConfigMap** stores URLs to monitor
2. **Deployment** runs the Go app with mounted config
3. **App** reads URLs from `/config/urls.txt` and checks them every 30s
4. **Kubernetes** ensures the service stays running
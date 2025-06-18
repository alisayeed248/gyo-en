Write-Host "Nuclear deploy - destroying everything."

# Delete everything
kubectl delete deployment gyo-en-deployment --force --grace-period=0
kubectl delete pods -l app=gyo-en --force --grace-period=0

# Wait for cleanup
Start-Sleep -Seconds 10

# Fresh build and deploy
docker build --no-cache -t gyo-en:latest .
minikube image load gyo-en:latest

# Fresh deployment
kubectl apply -f k8s-deployment.yaml

Write-Host "Deployed. Waiting for startup..."
Start-Sleep -Seconds 20
kubectl logs -l app=gyo-en --tail=50
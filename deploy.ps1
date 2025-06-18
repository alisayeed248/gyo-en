$TAG = [int][double]::Parse((Get-Date -UFormat %s))
Write-Host "Building gyo-en:$TAG"

docker build -t gyo-en:$TAG .
minikube image load gyo-en:$TAG

(Get-Content k8s-deployment.yaml) -replace 'image: gyo-en:.*', "image: gyo-en:$TAG" | kubectl apply -f -

Write-Host "Deployed gyo-en:$TAG"
kubectl logs -l app=gyo-en -f


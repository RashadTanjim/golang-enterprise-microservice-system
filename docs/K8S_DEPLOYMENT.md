# Kubernetes Deployment Guide (Local + AWS EC2)

This guide explains how to deploy the Enterprise Microservice System on:
- Local Kubernetes (kind or minikube)
- AWS EC2 (single-node k3s)

It also shows how to configure CI/CD with GitHub Actions.

## 1) Prerequisites

- Docker
- kubectl
- Git
- Access to GHCR (GitHub Container Registry)

## 2) Image Registry Setup (GHCR)

The CD workflow builds and pushes images to GHCR. Update the image names in `k8s/*.yaml` to your repo:

```
image: ghcr.io/<your-org>/<your-repo>/user-service:latest
```

For example:
```
image: ghcr.io/rashadtanjim/enterprise-microservice-system/user-service:latest
```

## 3) Local Kubernetes (kind)

### 3.1 Create a local cluster

```
cat > kind-config.yaml <<'KIND'
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    extraPortMappings:
      - containerPort: 80
        hostPort: 8080
      - containerPort: 443
        hostPort: 8443
KIND

kind create cluster --name enterprise-ms --config kind-config.yaml
```

### 3.2 Install ingress-nginx (optional but recommended)

```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
kubectl -n ingress-nginx wait --for=condition=available deployment/ingress-nginx-controller --timeout=180s
```

### 3.3 Deploy manifests

```
kubectl apply -k k8s
kubectl -n enterprise-ms apply -f k8s/migrations.yaml
```

### 3.4 Verify

```
kubectl -n enterprise-ms get pods
kubectl -n enterprise-ms get svc
```

If you enabled ingress-nginx, add a host entry:
```
127.0.0.1 enterprise.local
```

Then open:
```
http://enterprise.local
```

## 4) Local Kubernetes (minikube)

```
minikube start
minikube addons enable ingress
kubectl apply -k k8s
kubectl -n enterprise-ms apply -f k8s/migrations.yaml
```

If you use ingress, get the IP:
```
minikube ip
```

Add it to `/etc/hosts`:
```
<minikube-ip> enterprise.local
```

## 5) AWS EC2 (Single Node k3s)

### 5.1 Provision EC2

- Ubuntu 22.04 (t3.medium or bigger recommended)
- Open inbound ports in the security group:
  - 22 (SSH)
  - 80 (HTTP)
  - 443 (HTTPS)

### 5.2 Install k3s

```
curl -sfL https://get.k3s.io | sh -

# Add kubeconfig for your user
mkdir -p ~/.kube
sudo cat /etc/rancher/k3s/k3s.yaml > ~/.kube/config
sudo chown $(id -u):$(id -g) ~/.kube/config
```

Verify:
```
kubectl get nodes
```

### 5.3 Deploy the system

```
kubectl apply -k k8s
kubectl -n enterprise-ms apply -f k8s/migrations.yaml
```

### 5.4 Access the portal

Set the ingress hostname in your local `/etc/hosts`:
```
<EC2_PUBLIC_IP> enterprise.local
```

Then open:
```
http://enterprise.local
```

## 6) CI/CD (GitHub Actions)

The `CD` workflow builds and pushes images to GHCR and deploys to Kubernetes.

### 6.1 Generate KUBE_CONFIG_DATA

On a machine where `kubectl` can access the cluster:
```
base64 -w 0 ~/.kube/config
```

Add the output as the GitHub Actions secret `KUBE_CONFIG_DATA`.

### 6.2 Required Secrets

- `KUBE_CONFIG_DATA`
- (Optional) any other secrets you want to manage in `k8s/secret.yaml`

## 7) Managed Databases (Optional)

If you want to use managed Postgres/Redis (RDS/ElastiCache):
- Remove `k8s/postgres.yaml` and `k8s/redis.yaml` from `k8s/kustomization.yaml`.
- Update `k8s/configmap.yaml` and `k8s/secret.yaml` with managed host credentials.

## 8) Troubleshooting

- Check pods: `kubectl -n enterprise-ms get pods`
- Logs: `kubectl -n enterprise-ms logs deploy/user-service`
- Describe failures: `kubectl -n enterprise-ms describe pod <pod-name>`
- Re-run migrations: `kubectl -n enterprise-ms apply -f k8s/migrations.yaml`

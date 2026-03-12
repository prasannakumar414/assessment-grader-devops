# Step 5 - Deploy to Kubernetes and Complete K8s Check

## Goal

Run your app in a local Kubernetes cluster and let the health operator report your K8s result to the grader server.

## 1) Install and Start Minikube

Install minikube: [minikube.sigs.k8s.io/docs/start](https://minikube.sigs.k8s.io/docs/start/)

```bash
minikube start
kubectl get nodes
```

## 2) Install the Health Operator

From this repo's `operator/` directory:

```bash
cd operator
kubectl apply -f examples/crd.yaml
kubectl apply -f examples/rbac.yaml
kubectl apply -f examples/deployment.yaml
```

Note: `examples/deployment.yaml` uses `controller:latest` by default. Use the operator image provided by your instructor if needed.

## 3) Deploy Your App in Kubernetes

Create a Deployment + Service named `my-app` in `default` namespace, using your Docker Hub image:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
        - name: app
          image: <your-dockerhub-username>/docker-assessment-test:latest
          env:
            - name: EMAIL
              value: "your-email@example.com"
          ports:
            - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: my-app
  namespace: default
spec:
  selector:
    app: my-app
  ports:
    - port: 80
      targetPort: 8080
```

Apply it:

```bash
kubectl apply -f my-app.yaml
kubectl get deploy,svc -n default
```

## 4) Create AppHealth Custom Resource

Create `my-health-check.yaml`:

```yaml
apiVersion: example.com/v1
kind: AppHealth
metadata:
  name: my-health-check
spec:
  namespace: default
  appName: my-app
  reportURL: http://<instructor-server-ip>:8080
  infoPath: /api/info
```

Apply it:

```bash
kubectl apply -f my-health-check.yaml
kubectl get apphealth -o wide
kubectl describe apphealth my-health-check
```

## 5) Verify in Grader Dashboard

The operator checks your deployment and reports to:

`POST http://<instructor-server-ip>:8080/api/notify` with `stage: "k8s"`

Ask your instructor to verify K8s status is updated in the admin dashboard.

## Very Important: `/info` vs `/api/info`

- Operator default path is `/info`
- Starter app exposes `/api/info`

Use either approach:

- Preferred: set `infoPath: /api/info` in AppHealth spec (as shown above), or
- Add `/info` endpoint in your app code and rebuild/push image

## Common Failures

- Deployment and Service names do not match `spec.appName`
- Service/deployment not in `spec.namespace`
- `reportURL` is wrong or not reachable from cluster
- Operator cannot pull or run properly
- Endpoint path mismatch (`/info` vs `/api/info`)

## Useful Debug Commands

```bash
kubectl logs -n system deploy/health-operator
kubectl get apphealth my-health-check -o yaml
kubectl get pods -A
```

Then complete with [Troubleshooting and FAQ](./troubleshooting.md) if needed.

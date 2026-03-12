# Example Kubernetes manifests for health-operator

Apply in order (or use `kubectl apply -k config/default` for full kustomize deployment).

1. **crd.yaml** – Install the AppHealth CRD.
2. **rbac.yaml** – ServiceAccount, ClusterRole, ClusterRoleBinding for the operator.
3. **deployment.yaml** – Operator Deployment (update image to your registry).
4. **apphealth-sample.yaml** – Sample AppHealth to validate a student deployment.
5. **my-app.yaml** – Minimal test app (Deployment + Service) that serves `/info` with an email. Use to verify the operator without a real student app: `kubectl apply -f examples/my-app.yaml`.

## Quick apply (after building the operator image)

```bash
# Install CRD
kubectl apply -f examples/crd.yaml

# Install RBAC and operator
kubectl apply -f examples/rbac.yaml
kubectl apply -f examples/deployment.yaml

# Create an AppHealth to validate namespace "student-1", app name "my-app"
kubectl apply -f examples/apphealth-sample.yaml

# Check validation status
kubectl get apphealth -o wide
kubectl describe apphealth apphealth-sample
```

## Student app requirements

The operator expects the student to have deployed:

- A **Deployment** and **Service** with the same name (e.g. `my-app`) in the specified namespace.
- The app must expose an HTTP endpoint (default `/info`) that returns JSON: `{"email": "student@example.com"}`.

Optionally set `spec.expectedImage` to validate that the deployment uses a specific image.

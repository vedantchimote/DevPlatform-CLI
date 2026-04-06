# Helm Chart Customization

Complete guide to customizing the DevPlatform Base Helm chart for your applications.

## Table of Contents

- [Overview](#overview)
- [Chart Structure](#chart-structure)
- [values.yaml Reference](#valuesyaml-reference)
- [Environment-Specific Values](#environment-specific-values)
- [Custom Values Files](#custom-values-files)
- [Common Customizations](#common-customizations)
- [Advanced Configuration](#advanced-configuration)
- [Examples](#examples)

## Overview

The DevPlatform Base Helm chart provides a production-ready template for deploying applications to Kubernetes. It includes:

- Deployment with configurable replicas
- Service for internal/external access
- Ingress for HTTP/HTTPS routing
- Horizontal Pod Autoscaler (HPA)
- Pod Disruption Budget (PDB)
- Security contexts and best practices
- Resource requests and limits
- Health checks (liveness and readiness probes)

**Chart Location**: `charts/devplatform-base/`

## Chart Structure

```
charts/devplatform-base/
├── Chart.yaml                 # Chart metadata
├── values.yaml                # Default values
├── values-dev.yaml            # Development overrides
├── values-staging.yaml        # Staging overrides
├── values-prod.yaml           # Production overrides
├── templates/
│   ├── deployment.yaml        # Deployment template
│   ├── service.yaml           # Service template
│   ├── ingress.yaml           # Ingress template
│   ├── hpa.yaml               # HPA template
│   ├── pdb.yaml               # PDB template
│   ├── serviceaccount.yaml    # Service account template
│   ├── _helpers.tpl           # Template helpers
│   └── NOTES.txt              # Post-install notes
└── README.md                  # Chart documentation
```

## values.yaml Reference

### Application Configuration

```yaml
# Application metadata
app:
  name: myapp                  # Application name
  environment: dev             # Environment: dev, staging, prod
```

### Image Configuration

```yaml
image:
  repository: nginx            # Container image repository
  tag: latest                  # Image tag
  pullPolicy: IfNotPresent     # Pull policy: Always, IfNotPresent, Never
  
imagePullSecrets: []           # Image pull secrets for private registries
# - name: regcred
```

### Replica Configuration

```yaml
replicaCount: 1                # Number of pod replicas

# Autoscaling (overrides replicaCount when enabled)
autoscaling:
  enabled: false               # Enable HPA
  minReplicas: 1               # Minimum replicas
  maxReplicas: 10              # Maximum replicas
  targetCPUUtilizationPercentage: 80     # CPU target
  targetMemoryUtilizationPercentage: 80  # Memory target
```

### Service Configuration

```yaml
service:
  type: ClusterIP              # Service type: ClusterIP, NodePort, LoadBalancer
  port: 80                     # Service port
  targetPort: 8080             # Container port
  annotations: {}              # Service annotations
```

### Ingress Configuration

```yaml
ingress:
  enabled: false               # Enable ingress
  className: nginx             # Ingress class
  annotations: {}              # Ingress annotations
    # cert-manager.io/cluster-issuer: letsencrypt-prod
    # nginx.ingress.kubernetes.io/ssl-redirect: "true"
  
  hosts:
    - host: myapp.example.com  # Hostname
      paths:
        - path: /              # Path
          pathType: Prefix     # Path type: Prefix, Exact
  
  tls: []                      # TLS configuration
  # - secretName: myapp-tls
  #   hosts:
  #     - myapp.example.com
```

### Resource Configuration

```yaml
resources:
  requests:
    cpu: 100m                  # CPU request
    memory: 128Mi              # Memory request
  limits:
    cpu: 500m                  # CPU limit
    memory: 512Mi              # Memory limit
```

### Health Checks

```yaml
livenessProbe:
  httpGet:
    path: /health              # Health check path
    port: 8080                 # Health check port
  initialDelaySeconds: 30      # Initial delay
  periodSeconds: 10            # Check interval
  timeoutSeconds: 5            # Timeout
  failureThreshold: 3          # Failures before restart

readinessProbe:
  httpGet:
    path: /ready               # Readiness check path
    port: 8080                 # Readiness check port
  initialDelaySeconds: 10      # Initial delay
  periodSeconds: 5             # Check interval
  timeoutSeconds: 3            # Timeout
  failureThreshold: 3          # Failures before marking unready
```

### Environment Variables

```yaml
env: []                        # Environment variables
# - name: DATABASE_HOST
#   value: postgres.example.com
# - name: DATABASE_PASSWORD
#   valueFrom:
#     secretKeyRef:
#       name: db-secret
#       key: password
```

### Security Context

```yaml
securityContext:
  runAsNonRoot: true           # Run as non-root user
  runAsUser: 1000              # User ID
  fsGroup: 1000                # File system group
  capabilities:
    drop:
      - ALL                    # Drop all capabilities

containerSecurityContext:
  allowPrivilegeEscalation: false
  readOnlyRootFilesystem: true
  runAsNonRoot: true
  runAsUser: 1000
```

### Pod Disruption Budget

```yaml
podDisruptionBudget:
  enabled: false               # Enable PDB
  minAvailable: 1              # Minimum available pods
  # maxUnavailable: 1          # Maximum unavailable pods (alternative)
```

### Service Account

```yaml
serviceAccount:
  create: true                 # Create service account
  annotations: {}              # Service account annotations
    # eks.amazonaws.com/role-arn: arn:aws:iam::123456789012:role/my-role
    # azure.workload.identity/client-id: 12345678-1234-1234-1234-123456789012
  name: ""                     # Service account name (auto-generated if empty)
```

### Node Affinity and Tolerations

```yaml
nodeSelector: {}               # Node selector
  # disktype: ssd

tolerations: []                # Tolerations
  # - key: "key1"
  #   operator: "Equal"
  #   value: "value1"
  #   effect: "NoSchedule"

affinity: {}                   # Affinity rules
  # podAntiAffinity:
  #   preferredDuringSchedulingIgnoredDuringExecution:
  #     - weight: 100
  #       podAffinityTerm:
  #         labelSelector:
  #           matchExpressions:
  #             - key: app
  #               operator: In
  #               values:
  #                 - myapp
  #         topologyKey: kubernetes.io/hostname
```

## Environment-Specific Values

### Development (values-dev.yaml)

Optimized for development and testing:

```yaml
app:
  environment: dev

replicaCount: 1                # Single replica

resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 512Mi

autoscaling:
  enabled: false               # No autoscaling

ingress:
  enabled: false               # No ingress (use port-forward)

podDisruptionBudget:
  enabled: false               # No PDB needed
```

### Staging (values-staging.yaml)

Balanced configuration for pre-production testing:

```yaml
app:
  environment: staging

replicaCount: 2                # Two replicas

resources:
  requests:
    cpu: 250m
    memory: 256Mi
  limits:
    cpu: 1000m
    memory: 1Gi

autoscaling:
  enabled: true                # Enable autoscaling
  minReplicas: 2
  maxReplicas: 5
  targetCPUUtilizationPercentage: 80

ingress:
  enabled: true                # Enable ingress
  hosts:
    - host: myapp-staging.example.com
      paths:
        - path: /
          pathType: Prefix

podDisruptionBudget:
  enabled: true                # Enable PDB
  minAvailable: 1
```

### Production (values-prod.yaml)

High availability configuration for production:

```yaml
app:
  environment: prod

replicaCount: 3                # Three replicas

resources:
  requests:
    cpu: 500m
    memory: 512Mi
  limits:
    cpu: 2000m
    memory: 2Gi

autoscaling:
  enabled: true                # Enable autoscaling
  minReplicas: 3
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70

ingress:
  enabled: true                # Enable ingress
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
  hosts:
    - host: myapp.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: myapp-tls
      hosts:
        - myapp.example.com

podDisruptionBudget:
  enabled: true                # Enable PDB
  minAvailable: 2              # Keep at least 2 pods available

affinity:
  podAntiAffinity:             # Spread pods across nodes
    preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 100
        podAffinityTerm:
          labelSelector:
            matchExpressions:
              - key: app
                operator: In
                values:
                  - myapp
          topologyKey: kubernetes.io/hostname
```

## Custom Values Files

### Creating Custom Values

Create a custom values file for your application:

```yaml
# custom-values.yaml

# Override image
image:
  repository: myregistry/payment-service
  tag: v2.0.0

# Add environment variables
env:
  - name: DATABASE_HOST
    value: payment-db.example.com
  - name: DATABASE_PORT
    value: "5432"
  - name: DATABASE_NAME
    value: payment
  - name: DATABASE_USER
    value: payment_user
  - name: DATABASE_PASSWORD
    valueFrom:
      secretKeyRef:
        name: payment-db-secret
        key: password
  - name: LOG_LEVEL
    value: info
  - name: FEATURE_FLAG_NEW_UI
    value: "true"

# Increase resources
resources:
  requests:
    cpu: 500m
    memory: 1Gi
  limits:
    cpu: 2000m
    memory: 4Gi

# Custom health checks
livenessProbe:
  httpGet:
    path: /api/health
    port: 8080
  initialDelaySeconds: 60

readinessProbe:
  httpGet:
    path: /api/ready
    port: 8080
  initialDelaySeconds: 30

# Custom ingress
ingress:
  enabled: true
  hosts:
    - host: payment.example.com
      paths:
        - path: /
          pathType: Prefix
    - host: pay.example.com
      paths:
        - path: /
          pathType: Prefix
```

### Using Custom Values

```bash
# With DevPlatform CLI
devplatform create --app payment --env staging --provider aws \
  --values-file custom-values.yaml

# With Helm directly
helm install payment ./charts/devplatform-base \
  --namespace payment-staging \
  --values ./charts/devplatform-base/values-staging.yaml \
  --values custom-values.yaml
```

### Values Precedence

Values are merged with the following precedence (highest to lowest):

1. **CLI flags** (`--set` flags)
2. **Custom values file** (`--values-file`)
3. **Environment values** (`values-{env}.yaml`)
4. **Default values** (`values.yaml`)

Example:
```bash
helm install myapp ./charts/devplatform-base \
  --values values-prod.yaml \
  --values custom-values.yaml \
  --set image.tag=v3.0.0
```

## Common Customizations

### 1. Custom Container Image

```yaml
image:
  repository: myregistry.azurecr.io/myapp
  tag: v1.2.3
  pullPolicy: Always

imagePullSecrets:
  - name: acr-secret
```

### 2. Environment Variables from ConfigMap

```yaml
envFrom:
  - configMapRef:
      name: myapp-config
  - secretRef:
      name: myapp-secrets
```

### 3. Multiple Containers (Sidecar)

```yaml
# In custom values
extraContainers:
  - name: log-forwarder
    image: fluent/fluent-bit:latest
    volumeMounts:
      - name: logs
        mountPath: /var/log

extraVolumes:
  - name: logs
    emptyDir: {}
```

### 4. Init Containers

```yaml
initContainers:
  - name: wait-for-db
    image: busybox:latest
    command:
      - sh
      - -c
      - |
        until nc -z postgres.example.com 5432; do
          echo "Waiting for database..."
          sleep 2
        done
```

### 5. Persistent Storage

```yaml
persistence:
  enabled: true
  storageClass: gp3
  accessMode: ReadWriteOnce
  size: 10Gi
  mountPath: /data
```

### 6. Custom Service Ports

```yaml
service:
  type: ClusterIP
  ports:
    - name: http
      port: 80
      targetPort: 8080
    - name: metrics
      port: 9090
      targetPort: 9090
```

### 7. AWS IRSA Configuration

```yaml
serviceAccount:
  create: true
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::123456789012:role/payment-service-role
```

### 8. Azure Workload Identity

```yaml
serviceAccount:
  create: true
  annotations:
    azure.workload.identity/client-id: 12345678-1234-1234-1234-123456789012
    azure.workload.identity/tenant-id: 87654321-4321-4321-4321-210987654321

podLabels:
  azure.workload.identity/use: "true"
```

### 9. Custom Ingress Annotations

```yaml
ingress:
  enabled: true
  className: nginx
  annotations:
    # SSL/TLS
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    
    # Rate limiting
    nginx.ingress.kubernetes.io/limit-rps: "100"
    
    # CORS
    nginx.ingress.kubernetes.io/enable-cors: "true"
    nginx.ingress.kubernetes.io/cors-allow-origin: "https://example.com"
    
    # Timeouts
    nginx.ingress.kubernetes.io/proxy-connect-timeout: "60"
    nginx.ingress.kubernetes.io/proxy-send-timeout: "60"
    nginx.ingress.kubernetes.io/proxy-read-timeout: "60"
```

### 10. Resource Quotas

```yaml
resources:
  requests:
    cpu: 1000m
    memory: 2Gi
  limits:
    cpu: 4000m
    memory: 8Gi
```

## Advanced Configuration

### Blue-Green Deployment

```yaml
# blue-values.yaml
app:
  name: myapp-blue
  version: v1

image:
  tag: v1.0.0

service:
  selector:
    version: v1

# green-values.yaml
app:
  name: myapp-green
  version: v2

image:
  tag: v2.0.0

service:
  selector:
    version: v2
```

### Canary Deployment

```yaml
# Use with Flagger or Argo Rollouts
canary:
  enabled: true
  steps:
    - setWeight: 10
    - pause: {duration: 5m}
    - setWeight: 25
    - pause: {duration: 5m}
    - setWeight: 50
    - pause: {duration: 5m}
    - setWeight: 75
    - pause: {duration: 5m}
```

### Multi-Region Deployment

```yaml
# us-east-values.yaml
app:
  region: us-east

ingress:
  hosts:
    - host: us-east.myapp.example.com

# eu-west-values.yaml
app:
  region: eu-west

ingress:
  hosts:
    - host: eu-west.myapp.example.com
```

## Examples

### Example 1: Simple Web Application

```yaml
# web-app-values.yaml
image:
  repository: nginx
  tag: 1.25-alpine

replicaCount: 2

service:
  type: ClusterIP
  port: 80
  targetPort: 80

ingress:
  enabled: true
  hosts:
    - host: webapp.example.com
      paths:
        - path: /
          pathType: Prefix

resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 512Mi
```

### Example 2: API Service with Database

```yaml
# api-service-values.yaml
image:
  repository: myregistry/api-service
  tag: v1.5.0

replicaCount: 3

env:
  - name: DATABASE_HOST
    value: postgres.example.com
  - name: DATABASE_PORT
    value: "5432"
  - name: DATABASE_NAME
    value: api_db
  - name: DATABASE_USER
    valueFrom:
      secretKeyRef:
        name: db-credentials
        key: username
  - name: DATABASE_PASSWORD
    valueFrom:
      secretKeyRef:
        name: db-credentials
        key: password
  - name: REDIS_HOST
    value: redis.example.com
  - name: LOG_LEVEL
    value: info

service:
  type: ClusterIP
  port: 80
  targetPort: 8080

ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: api.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: api-tls
      hosts:
        - api.example.com

resources:
  requests:
    cpu: 500m
    memory: 1Gi
  limits:
    cpu: 2000m
    memory: 4Gi

autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70

livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 60
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 5
```

### Example 3: Background Worker

```yaml
# worker-values.yaml
image:
  repository: myregistry/worker
  tag: v2.0.0

replicaCount: 5

env:
  - name: QUEUE_URL
    value: sqs://queue.example.com
  - name: WORKER_THREADS
    value: "4"
  - name: LOG_LEVEL
    value: debug

service:
  enabled: false  # No service needed for worker

ingress:
  enabled: false  # No ingress needed for worker

resources:
  requests:
    cpu: 1000m
    memory: 2Gi
  limits:
    cpu: 4000m
    memory: 8Gi

autoscaling:
  enabled: true
  minReplicas: 5
  maxReplicas: 20
  targetCPUUtilizationPercentage: 80

# No HTTP probes for worker
livenessProbe:
  exec:
    command:
      - /app/healthcheck.sh
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  exec:
    command:
      - /app/readycheck.sh
  initialDelaySeconds: 10
  periodSeconds: 5
```

## Validation

### Validate Values

```bash
# Lint chart
helm lint ./charts/devplatform-base --values custom-values.yaml

# Dry-run install
helm install myapp ./charts/devplatform-base \
  --namespace myapp-dev \
  --values custom-values.yaml \
  --dry-run --debug

# Template rendering
helm template myapp ./charts/devplatform-base \
  --namespace myapp-dev \
  --values custom-values.yaml
```

### Test Deployment

```bash
# Install with test
helm install myapp ./charts/devplatform-base \
  --namespace myapp-dev \
  --values custom-values.yaml \
  --wait --timeout 5m

# Run tests
helm test myapp --namespace myapp-dev
```

## Troubleshooting

### Common Issues

1. **Image Pull Errors**: Check `imagePullSecrets` and registry credentials
2. **Pod Crashes**: Review `resources`, `livenessProbe`, and container logs
3. **Ingress Not Working**: Verify `ingress.className` and DNS configuration
4. **HPA Not Scaling**: Check metrics-server installation and resource requests
5. **PDB Blocking Updates**: Adjust `minAvailable` or `maxUnavailable`

### Debug Commands

```bash
# Check pod status
kubectl get pods -n myapp-dev

# View pod logs
kubectl logs -n myapp-dev -l app=myapp

# Describe pod
kubectl describe pod -n myapp-dev <pod-name>

# Check events
kubectl get events -n myapp-dev --sort-by='.lastTimestamp'

# Test service
kubectl port-forward -n myapp-dev svc/myapp 8080:80
```

## See Also

- [README.md](../README.md) - Main documentation
- [Command Reference](COMMAND_REFERENCE.md) - CLI commands
- [Terraform Modules](TERRAFORM_MODULES.md) - Infrastructure documentation
- [Helm Chart README](../charts/devplatform-base/README.md) - Chart-specific docs

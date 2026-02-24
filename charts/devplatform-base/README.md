# DevPlatform Base Helm Chart

A base Helm chart for deploying applications on Kubernetes using the DevPlatform CLI.

## Features

- **Environment-specific configurations**: Separate values files for dev, staging, and prod
- **Security-first**: Pod and container security contexts enabled by default
- **Autoscaling**: Horizontal Pod Autoscaler support
- **High availability**: Pod disruption budgets and anti-affinity rules
- **Observability**: Liveness and readiness probes
- **Ingress support**: TLS and multiple host configuration
- **Flexible**: Extensive customization through values

## Installation

### Using DevPlatform CLI (Recommended)

```bash
devplatform create --app myapp --env dev
```

### Using Helm directly

```bash
# Development
helm install myapp ./charts/devplatform-base \
  --namespace myapp-dev \
  --create-namespace \
  --values ./charts/devplatform-base/values-dev.yaml \
  --set app.name=myapp \
  --set image.repository=myregistry/myapp \
  --set image.tag=v1.0.0

# Staging
helm install myapp ./charts/devplatform-base \
  --namespace myapp-staging \
  --create-namespace \
  --values ./charts/devplatform-base/values-staging.yaml \
  --set app.name=myapp \
  --set image.repository=myregistry/myapp \
  --set image.tag=v1.0.0

# Production
helm install myapp ./charts/devplatform-base \
  --namespace myapp-prod \
  --create-namespace \
  --values ./charts/devplatform-base/values-prod.yaml \
  --set app.name=myapp \
  --set image.repository=myregistry/myapp \
  --set image.tag=v1.0.0
```

## Configuration

### Environment-Specific Values

The chart includes pre-configured values files for each environment:

- `values-dev.yaml`: Development environment (1 replica, no autoscaling, no ingress)
- `values-staging.yaml`: Staging environment (2 replicas, autoscaling 2-5, ingress enabled)
- `values-prod.yaml`: Production environment (3 replicas, autoscaling 3-10, ingress enabled, HA)

### Resource Limits

Resources are automatically configured based on the environment:

| Environment | CPU Request | CPU Limit | Memory Request | Memory Limit |
|-------------|-------------|-----------|----------------|--------------|
| dev         | 100m        | 500m      | 128Mi          | 512Mi        |
| staging     | 250m        | 1000m     | 256Mi          | 1Gi          |
| prod        | 500m        | 2000m     | 512Mi          | 2Gi          |

### Common Customizations

#### Custom Image

```bash
helm install myapp ./charts/devplatform-base \
  --set image.repository=myregistry/myapp \
  --set image.tag=v2.0.0
```

#### Environment Variables

```bash
helm install myapp ./charts/devplatform-base \
  --set env[0].name=DATABASE_HOST \
  --set env[0].value=postgres.example.com \
  --set env[1].name=LOG_LEVEL \
  --set env[1].value=debug
```

Or use a custom values file:

```yaml
env:
  - name: DATABASE_HOST
    value: postgres.example.com
  - name: DATABASE_PASSWORD
    valueFrom:
      secretKeyRef:
        name: db-secret
        key: password
```

#### Custom Ingress

```yaml
ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: myapp.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: myapp-tls
      hosts:
        - myapp.example.com
```

## Values Reference

See [values.yaml](values.yaml) for the complete list of configurable values.

### Key Values

| Parameter | Description | Default |
|-----------|-------------|---------|
| `app.name` | Application name | `myapp` |
| `app.environment` | Environment type (dev/staging/prod) | `dev` |
| `image.repository` | Container image repository | `nginx` |
| `image.tag` | Container image tag | `latest` |
| `replicaCount` | Number of replicas | `1` |
| `service.type` | Kubernetes service type | `ClusterIP` |
| `service.port` | Service port | `80` |
| `ingress.enabled` | Enable ingress | `false` |
| `autoscaling.enabled` | Enable HPA | `false` |

## Upgrading

```bash
helm upgrade myapp ./charts/devplatform-base \
  --namespace myapp-dev \
  --values ./charts/devplatform-base/values-dev.yaml \
  --set image.tag=v1.1.0
```

## Uninstalling

```bash
helm uninstall myapp --namespace myapp-dev
```

Or using DevPlatform CLI:

```bash
devplatform destroy --app myapp --env dev
```

## Labels

All resources are labeled with:

- `app`: Application name from `app.name`
- `environment`: Environment type from `app.environment`
- `managed-by`: Always set to `devplatform-cli`
- Standard Helm labels: `app.kubernetes.io/name`, `app.kubernetes.io/instance`, etc.

## Security

The chart follows security best practices:

- Non-root user (UID 1000)
- Read-only root filesystem
- Dropped all capabilities
- No privilege escalation
- Security contexts at pod and container level

## Support

For issues and questions, please refer to the DevPlatform CLI documentation.

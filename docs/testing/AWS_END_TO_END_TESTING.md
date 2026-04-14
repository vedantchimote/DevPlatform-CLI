# AWS End-to-End Testing Guide

This guide provides comprehensive instructions for testing the DevPlatform CLI with AWS from start to finish.

## Prerequisites

### Required Tools
- ✅ DevPlatform CLI binary
- ✅ AWS CLI configured with credentials
- ✅ Terraform 1.5+
- ✅ Helm 3.x
- ✅ kubectl 1.27+

### AWS Requirements
- ✅ AWS account with appropriate permissions
- ✅ IAM user or role with permissions for:
  - VPC, Subnets, NAT Gateways, Internet Gateways
  - RDS (PostgreSQL)
  - EKS cluster access
  - Secrets Manager
  - S3 (for Terraform state)
  - DynamoDB (for state locking)

### AWS Credentials Setup
```bash
# Configure AWS CLI
aws configure

# Verify credentials
aws sts get-caller-identity
```

---

## Test Scenario 1: Dev Environment Full Lifecycle

### Step 1: Create Dev Environment
```bash
# Create a dev environment
devplatform-cli create --app testapp --env dev --provider aws

# Expected output:
# ⏳ Validating inputs...
# ✓ Inputs validated successfully
# ⏳ Loading configuration...
# ✓ Configuration loaded
# ⏳ Initializing AWS provider...
# ✓ AWS provider initialized
# ⏳ Validating AWS credentials...
# ✓ AWS credentials validated
# 📊 Estimated monthly cost: $XXX.XX
# ⏳ Provisioning infrastructure with Terraform...
# ✓ Infrastructure provisioned successfully
# ⏳ Deploying application with Helm...
# ✓ Application deployed successfully
# ⏳ Configuring kubectl access...
# ✓ kubectl configured
# 🎉 Environment created successfully!
```

### Step 2: Verify Resources in AWS Console
1. **VPC**: Navigate to VPC console
   - Verify VPC created with name `testapp-dev-vpc`
   - Verify 2 public subnets across 2 AZs
   - Verify 2 private subnets across 2 AZs
   - Verify 1 NAT Gateway
   - Verify Internet Gateway
   - Verify route tables configured correctly

2. **RDS**: Navigate to RDS console
   - Verify PostgreSQL instance `testapp-dev-db`
   - Verify instance class: `db.t3.micro`
   - Verify storage: 20GB
   - Verify single-AZ deployment
   - Verify security group allows access from EKS

3. **Secrets Manager**: Navigate to Secrets Manager
   - Verify secret `testapp-dev-db-password` exists
   - Verify secret contains database credentials

4. **EKS**: Navigate to EKS console
   - Verify namespace `testapp-dev` exists
   - Verify resource quotas configured
   - Verify service account with IRSA

### Step 3: Check Environment Status
```bash
# Check status in table format
devplatform-cli status --app testapp --env dev --provider aws

# Expected output:
# Environment Status: testapp-dev (AWS)
# ┌────────────┬────────┬─────────────────────────────────────┐
# │ Component  │ Status │ Details                             │
# ├────────────┼────────┼─────────────────────────────────────┤
# │ VPC        │ ✓      │ vpc-xxxxx                           │
# │ Database   │ ✓      │ testapp-dev-db.xxxxx.rds.amazonaws  │
# │ Namespace  │ ✓      │ testapp-dev                         │
# │ Pods       │ ✓      │ 1/1 Running                         │
# │ Ingress    │ ✓      │ http://xxxxx.elb.amazonaws.com      │
# └────────────┴────────┴─────────────────────────────────────┘
```

```bash
# Check status in JSON format
devplatform-cli status --app testapp --env dev --provider aws --output json

# Expected output: JSON with all resource details
```

```bash
# Watch status with auto-refresh every 5 seconds
devplatform-cli status --app testapp --env dev --provider aws --watch 5

# Press Ctrl+C to stop watching
```

### Step 4: Verify kubectl Access
```bash
# Update kubeconfig
aws eks update-kubeconfig --name <cluster-name> --region <region>

# Switch to namespace
kubectl config set-context --current --namespace=testapp-dev

# Verify pods
kubectl get pods

# Expected output:
# NAME                       READY   STATUS    RESTARTS   AGE
# testapp-dev-xxxxx-xxxxx    1/1     Running   0          5m
```

### Step 5: Verify Application Deployment
```bash
# Check Helm release
helm list -n testapp-dev

# Expected output:
# NAME        NAMESPACE   REVISION  STATUS    CHART                 APP VERSION
# testapp-dev testapp-dev 1         deployed  devplatform-base-0.1.0 1.0.0

# Check deployment
kubectl get deployment -n testapp-dev

# Check service
kubectl get service -n testapp-dev

# Check ingress
kubectl get ingress -n testapp-dev
```

### Step 6: Destroy Environment
```bash
# Destroy with confirmation prompt
devplatform-cli destroy --app testapp --env dev --provider aws

# Expected prompt:
# ⚠️  WARNING: This will destroy the following resources:
# - VPC and all networking components
# - RDS database instance
# - Kubernetes namespace and all resources
# - Helm release
#
# 💰 Estimated monthly savings: $XXX.XX
# 📅 Estimated annual savings: $X,XXX.XX
#
# Type 'yes' to confirm destruction: 

# Type 'yes' and press Enter

# Expected output:
# ⏳ Uninstalling Helm release...
# ✓ Helm release uninstalled
# ⏳ Destroying infrastructure with Terraform...
# ✓ Infrastructure destroyed successfully
# 💰 Monthly savings: $XXX.XX
# 📅 Annual savings: $X,XXX.XX
# 🎉 Environment destroyed successfully!
```

### Step 7: Verify Cleanup in AWS Console
1. **VPC**: Verify VPC deleted
2. **RDS**: Verify database instance deleted or in "deleting" state
3. **Secrets Manager**: Verify secret deleted or scheduled for deletion
4. **EKS**: Verify namespace deleted

---

## Test Scenario 2: Staging Environment with Custom Values

### Step 1: Create Custom Helm Values File
```yaml
# custom-values.yaml
replicaCount: 2

image:
  repository: myregistry/myapp
  tag: "v1.2.3"

resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "512Mi"
    cpu: "500m"

ingress:
  enabled: true
  hosts:
    - host: testapp-staging.example.com
      paths:
        - path: /
          pathType: Prefix
```

### Step 2: Create Staging Environment with Custom Values
```bash
devplatform-cli create \
  --app testapp \
  --env staging \
  --provider aws \
  --values-file custom-values.yaml

# Verify custom values applied
kubectl get deployment testapp-staging -n testapp-staging -o yaml | grep -A 10 resources
```

### Step 3: Verify Staging Resources
- Verify 2 NAT Gateways (staging uses 2 AZs with 2 NATs)
- Verify RDS instance class: `db.t3.small`
- Verify storage: 50GB
- Verify 2 pod replicas running

### Step 4: Cleanup
```bash
devplatform-cli destroy --app testapp --env staging --provider aws --confirm
```

---

## Test Scenario 3: Production Environment

### Step 1: Create Production Environment
```bash
devplatform-cli create --app testapp --env prod --provider aws
```

### Step 2: Verify Production Resources
- Verify 3 NAT Gateways (prod uses 3 AZs with 3 NATs)
- Verify RDS instance class: `db.r6g.large`
- Verify storage: 100GB
- Verify Multi-AZ deployment
- Verify 3 pod replicas
- Verify HPA (Horizontal Pod Autoscaler) configured
- Verify PDB (Pod Disruption Budget) configured

### Step 3: Verify High Availability
```bash
# Check HPA
kubectl get hpa -n testapp-prod

# Check PDB
kubectl get pdb -n testapp-prod

# Verify pods across multiple nodes
kubectl get pods -n testapp-prod -o wide
```

### Step 4: Cleanup
```bash
devplatform-cli destroy --app testapp --env prod --provider aws --confirm
```

---

## Test Scenario 4: Dry-Run Mode

### Step 1: Dry-Run Create
```bash
devplatform-cli create \
  --app testapp \
  --env dev \
  --provider aws \
  --dry-run

# Expected output:
# 🔍 DRY-RUN MODE: No resources will be created
# ⏳ Validating inputs...
# ✓ Inputs validated successfully
# ⏳ Loading configuration...
# ✓ Configuration loaded
# ⏳ Initializing AWS provider...
# ✓ AWS provider initialized
# ⏳ Validating AWS credentials...
# ✓ AWS credentials validated
# 📊 Estimated monthly cost: $XXX.XX
# ⏳ Running Terraform plan...
# [Terraform plan output showing what would be created]
# ✓ Terraform plan completed
# 🔍 DRY-RUN COMPLETE: No resources were created
```

### Step 2: Verify No Resources Created
- Check AWS console to verify no resources were created
- Verify no Terraform state file exists

---

## Test Scenario 5: Error Handling and Rollback

### Step 1: Simulate Terraform Failure
```bash
# Create environment with invalid configuration
# (This would require modifying Terraform modules to simulate failure)

# Expected behavior:
# - Error detected during Terraform apply
# - Automatic rollback initiated
# - Terraform destroy executed
# - Clear error message displayed
# - Manual cleanup instructions provided if rollback fails
```

### Step 2: Simulate Helm Failure
```bash
# Create environment with invalid Helm values
# (This would require invalid values file)

# Expected behavior:
# - Infrastructure provisioned successfully
# - Error detected during Helm install
# - Automatic rollback initiated
# - Helm uninstall executed (if partially installed)
# - Terraform destroy executed
# - Clear error message displayed
```

---

## Test Scenario 6: Concurrent Execution Safety

### Step 1: Create Multiple Environments Concurrently
```bash
# Terminal 1
devplatform-cli create --app app1 --env dev --provider aws

# Terminal 2 (run simultaneously)
devplatform-cli create --app app2 --env dev --provider aws

# Expected behavior:
# - Both commands execute successfully
# - Separate Terraform state files used
# - No state locking conflicts
# - Resources created independently
```

### Step 2: Test State Locking
```bash
# Terminal 1
devplatform-cli create --app testapp --env dev --provider aws

# Terminal 2 (run while Terminal 1 is still running)
devplatform-cli create --app testapp --env dev --provider aws

# Expected behavior:
# - Terminal 2 detects state lock
# - Clear error message about lock holder
# - Retry instructions provided
# - No corruption of Terraform state
```

### Step 3: Cleanup
```bash
devplatform-cli destroy --app app1 --env dev --provider aws --confirm
devplatform-cli destroy --app app2 --env dev --provider aws --confirm
```

---

## Test Scenario 7: Cost Estimation Verification

### Step 1: Verify Dev Environment Cost
```bash
devplatform-cli create --app testapp --env dev --provider aws --dry-run

# Expected cost breakdown:
# - VPC: ~$32/month (1 NAT Gateway)
# - RDS: ~$15/month (db.t3.micro)
# - EKS: ~$0/month (namespace only, cluster assumed to exist)
# Total: ~$47/month
```

### Step 2: Verify Staging Environment Cost
```bash
devplatform-cli create --app testapp --env staging --provider aws --dry-run

# Expected cost breakdown:
# - VPC: ~$64/month (2 NAT Gateways)
# - RDS: ~$30/month (db.t3.small)
# - EKS: ~$0/month
# Total: ~$94/month
```

### Step 3: Verify Production Environment Cost
```bash
devplatform-cli create --app testapp --env prod --provider aws --dry-run

# Expected cost breakdown:
# - VPC: ~$96/month (3 NAT Gateways)
# - RDS: ~$300/month (db.r6g.large, Multi-AZ)
# - EKS: ~$0/month
# Total: ~$396/month
```

---

## Test Scenario 8: Logging and Debugging

### Step 1: Test Verbose Mode
```bash
devplatform-cli create \
  --app testapp \
  --env dev \
  --provider aws \
  --verbose

# Expected output:
# - All log levels displayed (Debug, Info, Warn, Error)
# - Detailed progress information
# - Terraform command output
# - Helm command output
```

### Step 2: Test Debug Mode
```bash
devplatform-cli create \
  --app testapp \
  --env dev \
  --provider aws \
  --debug

# Expected output:
# - All verbose output
# - AWS API calls logged
# - Terraform plan details
# - Helm values being used
# - Kubernetes API interactions
```

### Step 3: Verify Log Files
```bash
# Check log directory
ls ~/.devplatform/logs/

# View latest log file
cat ~/.devplatform/logs/devplatform-<timestamp>.log

# Expected: JSON formatted logs with timestamps, levels, messages
```

---

## Verification Checklist

### ✅ Create Command
- [ ] Dev environment creates successfully
- [ ] Staging environment creates successfully
- [ ] Production environment creates successfully
- [ ] Custom Helm values applied correctly
- [ ] Dry-run mode works without creating resources
- [ ] Cost estimation displayed accurately
- [ ] Progress indicators displayed
- [ ] Connection information displayed

### ✅ Status Command
- [ ] Table output format works
- [ ] JSON output format works
- [ ] YAML output format works
- [ ] Watch mode works with auto-refresh
- [ ] All components show correct status
- [ ] Connection information displayed

### ✅ Destroy Command
- [ ] Confirmation prompt works
- [ ] --confirm flag skips prompt
- [ ] Helm uninstall executes before Terraform destroy
- [ ] Cost savings calculated and displayed
- [ ] --force flag handles partial failures
- [ ] --keep-state flag preserves state file
- [ ] All resources cleaned up

### ✅ Error Handling
- [ ] Terraform failures trigger rollback
- [ ] Helm failures trigger rollback
- [ ] Clear error messages displayed
- [ ] Resolution guidance provided
- [ ] Manual cleanup instructions provided
- [ ] Log file path displayed

### ✅ Multi-Environment
- [ ] Multiple apps can coexist
- [ ] Multiple environments per app can coexist
- [ ] State isolation works correctly
- [ ] No resource naming conflicts

### ✅ Concurrent Execution
- [ ] Multiple creates run simultaneously
- [ ] State locking prevents conflicts
- [ ] Lock holder information displayed
- [ ] Retry instructions provided

### ✅ Logging
- [ ] Console logging works
- [ ] File logging works
- [ ] Log rotation works
- [ ] Verbose mode works
- [ ] Debug mode works
- [ ] --no-color flag works

---

## Troubleshooting

### Issue: AWS Credentials Not Found
**Solution**: Configure AWS CLI with `aws configure`

### Issue: Terraform State Lock
**Solution**: Wait for other operation to complete or manually release lock

### Issue: RDS Creation Timeout
**Solution**: Increase timeout with `--timeout 45m`

### Issue: Helm Pod Not Ready
**Solution**: Check pod logs with `kubectl logs -n <namespace> <pod-name>`

### Issue: VPC Deletion Fails
**Solution**: Manually delete dependent resources (ENIs, security groups) then retry

---

## Success Criteria

✅ **All test scenarios pass**
✅ **Resources created correctly in AWS**
✅ **Resources cleaned up completely**
✅ **Error handling works as expected**
✅ **Logging captures all operations**
✅ **Cost estimates are accurate**
✅ **Multi-environment isolation works**
✅ **Concurrent execution is safe**


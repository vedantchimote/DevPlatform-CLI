# Manual Testing Guide for DevPlatform CLI

This guide explains how to complete the three remaining testing tasks that require manual intervention.

---

## Task 1: Install Go and Run Test Suite

### Why This Is Needed
The Go test suite (100+ unit tests, 27 integration tests) cannot run without Go installed on your system.

### Installation Steps

**Option 1: Using Chocolatey (Recommended for Windows)**
```powershell
# Install Chocolatey if not already installed
Set-ExecutionPolicy Bypass -Scope Process -Force
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Install Go
choco install golang -y

# Refresh environment variables
refreshenv

# Verify installation
go version
```

**Option 2: Manual Installation**
1. Download Go from: https://go.dev/dl/
2. Download: `go1.21.x.windows-amd64.msi` (or latest version)
3. Run the installer
4. Restart PowerShell
5. Verify: `go version`

### Running the Tests

Once Go is installed:

```powershell
# Navigate to project directory
cd "C:\Users\Rigel\Desktop\CV Projects\DevPlatform-CLI"

# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out

# Run specific package tests
go test -v ./internal/config/...
go test -v ./internal/terraform/...
go test -v ./internal/helm/...
go test -v ./internal/aws/...
go test -v ./internal/azure/...

# Run integration tests only
go test -v ./test/integration/...
```

### Expected Results

✅ **Success Indicators:**
- All tests pass
- Coverage report shows ~75% overall coverage
- No compilation errors
- All mocks work correctly

❌ **If Tests Fail:**
- Check error messages carefully
- Ensure all dependencies are in go.mod
- Run `go mod tidy` to fix dependencies
- Check that test fixtures exist in `test/fixtures/`

---

## Task 2: Fix Azure Authentication and Test Dry-Run

### The Problem
The DevPlatform CLI uses Azure SDK's DefaultAzureCredential, which requires a specific authentication scope that standard `az login` doesn't provide.

### Solution Options

**Option 1: Use Environment Variables (Recommended)**

Create a service principal for testing:

```powershell
# Create service principal
az ad sp create-for-rbac --name "devplatform-cli-test" --role Contributor --scopes /subscriptions/df05ba72-d7ee-4934-9d1d-13e1b59d971c

# Output will look like:
# {
#   "appId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
#   "displayName": "devplatform-cli-test",
#   "password": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
#   "tenant": "b9d92904-6a83-4a01-a865-4f4bf9a992ad"
# }

# Set environment variables (replace with your values)
$env:AZURE_CLIENT_ID = "your-appId"
$env:AZURE_CLIENT_SECRET = "your-password"
$env:AZURE_TENANT_ID = "b9d92904-6a83-4a01-a865-4f4bf9a992ad"
$env:AZURE_SUBSCRIPTION_ID = "df05ba72-d7ee-4934-9d1d-13e1b59d971c"

# Test dry-run
.\devplatform-cli.exe create --app testapp --env dev --provider azure --dry-run --verbose
```

**Option 2: Fix the Code (Alternative)**

If the authentication continues to fail, you may need to modify the Azure authentication code in `internal/azure/auth.go` to use Azure CLI credentials directly instead of DefaultAzureCredential.

### Testing Dry-Run

Once authentication is fixed:

```powershell
# Test dry-run (no resources created)
.\devplatform-cli.exe create --app testapp --env dev --provider azure --dry-run --verbose

# Expected output:
# ✓ Validating credentials...
# ✓ DRY-RUN MODE: No resources will be created
# ✓ Would create VNet (10.0.0.0/16)
# ✓ Would create Azure Database (B_Gen5_1)
# ✓ Would create AKS namespace
# ✓ Estimated cost: $12-20/month
```

---

## Task 3: Full End-to-End Testing with Real Azure Resources

### ⚠️ WARNING
This will create REAL Azure resources and consume your student credits!

### Prerequisites

Before proceeding, ensure you have:

1. **Go installed** (from Task 1)
2. **Azure authentication working** (from Task 2)
3. **Required tools installed:**

```powershell
# Check if tools are installed
terraform version
helm version
kubectl version --client

# If missing, install them:
choco install terraform helm kubernetes-cli -y
```

### Cost Estimate

**Development Environment (--env dev):**
- VNet: ~$5/month
- Azure Database (B_Gen5_1): ~$12/month
- AKS Namespace: ~$5/month (shared cluster)
- **Total: ~$22/month**

**Your student credit: $100** (enough for ~4 months of dev environment)

### Step-by-Step Testing

**Step 1: Create a Test Environment**

```powershell
# Create dev environment (smallest/cheapest)
.\devplatform-cli.exe create --app testcli --env dev --provider azure --verbose

# This will take 10-15 minutes
# Watch for:
# ✓ Validating credentials...
# ✓ Creating VNet...
# ✓ Creating Azure Database...
# ✓ Creating AKS namespace...
# ✓ Deploying application...
# ✓ Environment created successfully!
```

**Step 2: Verify the Environment**

```powershell
# Check status
.\devplatform-cli.exe status --app testcli --env dev --provider azure

# Expected output:
# Environment: dev-testcli-azure
# Status: Healthy
# Database: Running
# Pods: 2/2 Ready
# Ingress: Active
```

**Step 3: Test Status Command with Different Formats**

```powershell
# JSON output
.\devplatform-cli.exe status --app testcli --env dev --provider azure --output json

# YAML output
.\devplatform-cli.exe status --app testcli --env dev --provider azure --output yaml

# Watch mode (refresh every 5 seconds)
.\devplatform-cli.exe status --app testcli --env dev --provider azure --watch 5
# Press Ctrl+C to exit
```

**Step 4: Verify Azure Resources in Portal**

1. Go to: https://portal.azure.com
2. Navigate to Resource Groups
3. Find: `devplatform-testcli-dev-rg`
4. Verify resources:
   - Virtual Network
   - Azure Database for PostgreSQL
   - Network Security Groups
   - NAT Gateway

**Step 5: Test Kubernetes Access**

```powershell
# Get AKS credentials
az aks get-credentials --name shared-devplatform-cluster --resource-group devplatform-rg

# Set namespace context
kubectl config set-context --current --namespace=dev-testcli-azure

# Check pods
kubectl get pods

# Check services
kubectl get svc

# Check ingress
kubectl get ingress
```

**Step 6: Destroy the Environment (IMPORTANT!)**

```powershell
# Destroy to avoid ongoing charges
.\devplatform-cli.exe destroy --app testcli --env dev --provider azure --confirm

# This will take 5-10 minutes
# Watch for:
# ✓ Uninstalling Helm release...
# ✓ Destroying Terraform infrastructure...
# ✓ Environment destroyed successfully!
# ✓ Estimated monthly savings: $22
```

**Step 7: Verify Cleanup**

```powershell
# Check status (should fail)
.\devplatform-cli.exe status --app testcli --env dev --provider azure

# Expected: Error: environment not found

# Verify in Azure Portal
# Resource group should be deleted
```

### Monitoring Costs

```powershell
# Check your Azure spending
az consumption usage list --output table

# View student credit balance
# Go to: https://www.microsoftazuresponsoredstudents.com/
```

---

## Alternative: Use GitHub Actions CI/CD

Instead of running tests locally, you can push your code to GitHub and let the CI/CD pipeline run all tests automatically:

```powershell
# Push your commits to GitHub
git push origin main

# Go to GitHub Actions tab
# Watch the test workflow run
# All tests will execute in the cloud (free for public repos)
```

The `.github/workflows/test.yml` file is already configured to:
- Run all unit tests
- Run all integration tests
- Generate coverage reports
- Fail if coverage < 70%

---

## Summary

| Task | Can Be Automated | Requires Manual Steps | Risk Level |
|------|------------------|----------------------|------------|
| Install Go & Run Tests | ❌ No | ✅ Yes | 🟢 Low (no costs) |
| Fix Azure Auth & Dry-Run | ⚠️ Partial | ✅ Yes | 🟢 Low (no costs) |
| Full E2E Testing | ❌ No | ✅ Yes | 🟡 Medium (uses credits) |

**Recommendation:**
1. Install Go and run the test suite first (safest, no costs)
2. Fix Azure authentication and test dry-run (safe, no costs)
3. Only do full E2E testing if you need to verify real resource creation
4. Always destroy resources immediately after testing to avoid charges

---

## Getting Help

If you encounter issues:

1. **Go Installation Issues:**
   - Check PATH: `echo $env:PATH`
   - Restart PowerShell after installation
   - Try: `go env` to see Go configuration

2. **Azure Authentication Issues:**
   - Verify: `az account show`
   - Try: `az login` again
   - Check service principal credentials

3. **Test Failures:**
   - Read error messages carefully
   - Check `go.mod` for missing dependencies
   - Run `go mod tidy`
   - Verify test fixtures exist

4. **Resource Creation Issues:**
   - Check Azure quotas
   - Verify subscription is active
   - Check Terraform/Helm/kubectl versions
   - Review logs in `~/.devplatform/logs/`

---

**Last Updated**: April 18, 2026  
**Your Azure Subscription**: Azure for Students (df05ba72-d7ee-4934-9d1d-13e1b59d971c)  
**Student Credit Balance**: $100

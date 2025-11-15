# DevPlatform CLI - Architecture Documentation

## System Overview

The DevPlatform CLI is an Internal Developer Platform (IDP) that enables self-service infrastructure provisioning on AWS or Azure. It orchestrates Terraform for infrastructure management and Helm for Kubernetes deployments, providing a consistent developer experience across both cloud providers.

## High-Level Architecture

```mermaid
graph TB
    subgraph "Developer Workstation"
        CLI[DevPlatform CLI Binary]
        Config[.devplatform.yaml]
    end
    
    subgraph "CLI Core Components"
        Root[Root Command Handler]
        Create[Create Command]
        Status[Status Command]
        Destroy[Destroy Command]
        Validator[Config Validator]
    end
    
    subgraph "Internal Packages"
        TFWrapper[Terraform Wrapper]
        HelmWrapper[Helm Wrapper]
        Provider[Cloud Provider Abstraction]
        AWSUtil[AWS Provider]
        AzureUtil[Azure Provider]
        Logger[Logging System]
    end
    
    subgraph "External Tools"
        TF[Terraform Binary]
        Helm[Helm Binary]
        Kubectl[Kubectl Binary]
        AWSCLI[AWS CLI]
        AzureCLI[Azure CLI]
    end
    
    subgraph "AWS Cloud"
        S3[S3 State Backend]
        DynamoDB[DynamoDB Lock Table]
        VPC[VPC Resources]
        RDS[RDS Database]
        EKS[EKS Cluster]
        SecretsManager[Secrets Manager]
    end
    
    subgraph "Azure Cloud"
        AzureStorage[Azure Storage Backend]
        BlobLease[Blob Lease Lock]
        VNet[VNet Resources]
        AzureDB[Azure Database]
        AKS[AKS Cluster]
        KeyVault[Key Vault]
    end
    
    CLI --> Root
    Config --> Validator
    Root --> Create
    Root --> Status
    Root --> Destroy
    
    Create --> Validator
    Create --> TFWrapper
    Create --> HelmWrapper
    
    Status --> TFWrapper
    Status --> Provider
    
    Destroy --> HelmWrapper
    Destroy --> TFWrapper
    
    Provider --> AWSUtil
    Provider --> AzureUtil
    
    TFWrapper --> TF
    HelmWrapper --> Helm
    AWSUtil --> AWSCLI
    AWSUtil --> Kubectl
    AzureUtil --> AzureCLI
    AzureUtil --> Kubectl
    
    TF --> S3
    TF --> DynamoDB
    TF --> VPC
    TF --> RDS
    TF --> EKS
    TF --> SecretsManager
    
    TF --> AzureStorage
    TF --> BlobLease
    TF --> VNet
    TF --> AzureDB
    TF --> AKS
    TF --> KeyVault
    
    Helm --> EKS
    Helm --> AKS
    
    Create --> Logger
    Status --> Logger
    Destroy --> Logger
```

## Component Architecture

```mermaid
graph LR
    subgraph "cmd Package"
        RootCmd[root.go]
        CreateCmd[create.go]
        StatusCmd[status.go]
        DestroyCmd[destroy.go]
        VersionCmd[version.go]
    end
    
    subgraph "internal/config"
        Parser[parser.go]
        Validator[validator.go]
        Schema[schema.go]
    end
    
    subgraph "internal/provider"
        ProviderInterface[provider.go]
        Factory[factory.go]
    end
    
    subgraph "internal/terraform"
        Executor[executor.go]
        StateManager[state.go]
        OutputParser[output.go]
    end
    
    subgraph "internal/helm"
        Client[client.go]
        ChartManager[chart.go]
        ValuesMerger[values.go]
    end
    
    subgraph "internal/aws"
        AWSAuth[auth.go]
        AWSKubeConfig[kubeconfig.go]
        AWSPricing[pricing.go]
    end
    
    subgraph "internal/azure"
        AzureAuth[auth.go]
        AzureKubeConfig[kubeconfig.go]
        AzurePricing[pricing.go]
    end
    
    RootCmd --> CreateCmd
    RootCmd --> StatusCmd
    RootCmd --> DestroyCmd
    RootCmd --> VersionCmd
    
    CreateCmd --> Parser
    CreateCmd --> Validator
    CreateCmd --> Factory
    CreateCmd --> Executor
    CreateCmd --> Client
    
    StatusCmd --> StateManager
    StatusCmd --> Factory
    
    DestroyCmd --> Client
    DestroyCmd --> Executor
    
    Factory --> AWSAuth
    Factory --> AzureAuth
    
    Executor --> StateManager
    Executor --> OutputParser
    
    Client --> ChartManager
    Client --> ValuesMerger
    
    Parser --> Schema
    Validator --> Schema
```

## Data Flow Architecture

```mermaid
flowchart TD
    Start([Developer Input]) --> Parse[Parse CLI Arguments]
    Parse --> ValidateInput[Validate Input]
    ValidateInput --> SelectProvider{Cloud Provider?}
    
    SelectProvider -->|AWS| CheckAWS[Check AWS Credentials]
    SelectProvider -->|Azure| CheckAzure[Check Azure Credentials]
    
    CheckAWS -->|Valid| LoadConfig[Load Configuration]
    CheckAWS -->|Invalid| ErrorAuth[Display Auth Error]
    
    CheckAzure -->|Valid| LoadConfig
    CheckAzure -->|Invalid| ErrorAuth
    
    LoadConfig --> DryRun{Dry Run Mode?}
    
    DryRun -->|Yes| TFPlan[Terraform Plan]
    DryRun -->|No| TFInit[Terraform Init]
    
    TFPlan --> DisplayPlan[Display Plan Output]
    DisplayPlan --> End([Exit])
    
    TFInit --> TFApply[Terraform Apply]
    TFApply -->|Success| ExtractOutputs[Extract TF Outputs]
    TFApply -->|Failure| Rollback[Rollback Infrastructure]
    
    ExtractOutputs --> HelmInstall[Helm Install]
    HelmInstall -->|Success| VerifyPods[Verify Pods Running]
    HelmInstall -->|Failure| RollbackHelm[Rollback Helm + Infrastructure]
    
    VerifyPods --> UpdateKubeconfig[Update Kubeconfig]
    UpdateKubeconfig --> DisplaySuccess[Display Success Message]
    DisplaySuccess --> End
    
    Rollback --> DisplayError[Display Error]
    RollbackHelm --> DisplayError
    ErrorAuth --> End
    DisplayError --> End
```

## Deployment Architecture

```mermaid
graph TB
    subgraph "Developer Machine"
        CLI[DevPlatform CLI]
    end
    
    subgraph "AWS Account"
        subgraph "Shared Infrastructure (AWS)"
            S3State[S3 Bucket<br/>terraform-state]
            DynamoLock[DynamoDB Table<br/>terraform-locks]
            EKSCluster[EKS Cluster<br/>shared-cluster]
        end
        
        subgraph "App Environment: payment-dev (AWS)"
            VPC1[VPC<br/>10.0.0.0/16]
            RDS1[RDS Instance<br/>db.t3.micro]
            NS1[Namespace<br/>dev-payment]
            Pods1[Application Pods]
            Ingress1[ALB Ingress]
        end
    end
    
    subgraph "Azure Subscription"
        subgraph "Shared Infrastructure (Azure)"
            AzureStorageState[Azure Storage<br/>terraform-state]
            BlobLease[Blob Lease Lock]
            AKSCluster[AKS Cluster<br/>shared-cluster]
        end
        
        subgraph "App Environment: payment-dev (Azure)"
            VNet1[VNet<br/>10.0.0.0/16]
            AzureDB1[Azure Database<br/>B_Gen5_1]
            NS2[Namespace<br/>dev-payment]
            Pods2[Application Pods]
            Ingress2[Azure LB Ingress]
        end
    end
    
    CLI -->|--provider aws| VPC1
    CLI -->|--provider aws| RDS1
    CLI -->|--provider aws| NS1
    
    CLI -->|--provider azure| VNet1
    CLI -->|--provider azure| AzureDB1
    CLI -->|--provider azure| NS2
    
    VPC1 --> RDS1
    NS1 --> Pods1
    Pods1 --> Ingress1
    Pods1 -.->|Connect| RDS1
    
    VNet1 --> AzureDB1
    NS2 --> Pods2
    Pods2 --> Ingress2
    Pods2 -.->|Connect| AzureDB1
    
    CLI -.->|State (AWS)| S3State
    CLI -.->|Lock (AWS)| DynamoLock
    NS1 -.->|Runs in| EKSCluster
    
    CLI -.->|State (Azure)| AzureStorageState
    CLI -.->|Lock (Azure)| BlobLease
    NS2 -.->|Runs in| AKSCluster
```

## Security Architecture

```mermaid
graph TB
    subgraph "Authentication Layer"
        DevCreds[Developer Credentials]
        AWSAuth[AWS IAM Role/User]
        AzureAuth[Azure Service Principal/Managed Identity]
        AssumeRole[Assume Role / Get Token]
    end
    
    subgraph "Authorization Layer"
        IAMPolicies[IAM Policies (AWS)]
        AzureRBAC[Azure RBAC]
        K8sRBAC[Kubernetes RBAC]
        IRSA[IRSA (AWS)]
        WorkloadIdentity[Workload Identity (Azure)]
    end
    
    subgraph "Network Security"
        AWSSG[Security Groups (AWS)]
        AzureNSG[Network Security Groups (Azure)]
        NACL[Network ACLs (AWS)]
        PrivateSubnet[Private Subnets]
        NAT[NAT Gateway]
    end
    
    subgraph "Data Security"
        Encryption[Database Encryption at Rest]
        SecretsManager[Secrets Manager (AWS)]
        KeyVault[Key Vault (Azure)]
        TLS[TLS in Transit]
    end
    
    subgraph "Audit & Compliance"
        CloudTrail[CloudTrail (AWS)]
        ActivityLog[Activity Log (Azure)]
        VPCFlowLogs[VPC Flow Logs (AWS)]
        NSGFlowLogs[NSG Flow Logs (Azure)]
        Tags[Resource Tags]
    end
    
    DevCreds --> AWSAuth
    DevCreds --> AzureAuth
    AWSAuth --> AssumeRole
    AzureAuth --> AssumeRole
    AssumeRole --> IAMPolicies
    AssumeRole --> AzureRBAC
    
    IAMPolicies --> AWSSG
    IAMPolicies --> SecretsManager
    IAMPolicies --> CloudTrail
    
    AzureRBAC --> AzureNSG
    AzureRBAC --> KeyVault
    AzureRBAC --> ActivityLog
    
    K8sRBAC --> IRSA
    K8sRBAC --> WorkloadIdentity
    IRSA --> IAMPolicies
    WorkloadIdentity --> AzureRBAC
    
    AWSSG --> PrivateSubnet
    AzureNSG --> PrivateSubnet
    PrivateSubnet --> NAT
    AWSSG --> RDS[Database Instance]
    AzureNSG --> RDS
    
    RDS --> Encryption
    RDS --> TLS
    
    SecretsManager --> DBPassword[Database Password]
    KeyVault --> DBPassword
    
    CloudTrail --> AuditLog[Audit Logs]
    ActivityLog --> AuditLog
    VPCFlowLogs --> NetworkLog[Network Logs]
    NSGFlowLogs --> NetworkLog
    Tags --> CostTracking[Cost Tracking]
```

## State Management Architecture

```mermaid
stateDiagram-v2
    [*] --> Unprovisioned
    
    Unprovisioned --> Validating: devplatform create --provider aws|azure
    Validating --> Provisioning: Validation Success
    Validating --> Unprovisioned: Validation Failed
    
    Provisioning --> InfrastructureCreating: Terraform Init
    InfrastructureCreating --> ApplicationDeploying: Infrastructure Ready
    InfrastructureCreating --> RollingBack: Infrastructure Failed
    
    ApplicationDeploying --> Provisioned: Deployment Success
    ApplicationDeploying --> RollingBack: Deployment Failed
    
    RollingBack --> Unprovisioned: Rollback Complete
    
    Provisioned --> Checking: devplatform status --provider aws|azure
    Checking --> Provisioned: Status Retrieved
    
    Provisioned --> Destroying: devplatform destroy --provider aws|azure
    Destroying --> ApplicationRemoving: Confirmation Received
    ApplicationRemoving --> InfrastructureDestroying: App Removed
    InfrastructureDestroying --> Unprovisioned: Destruction Complete
    InfrastructureDestroying --> PartiallyDestroyed: Destruction Failed
    
    PartiallyDestroyed --> ManualCleanup: Manual Intervention
    ManualCleanup --> Unprovisioned: Cleanup Complete
```

## Concurrency Model

```mermaid
sequenceDiagram
    participant Dev1 as Developer 1
    participant Dev2 as Developer 2
    participant CLI1 as CLI Instance 1
    participant CLI2 as CLI Instance 2
    participant Lock as State Lock (S3+DynamoDB or Azure Storage)
    participant State as State Backend
    
    Dev1->>CLI1: create --app payment --env dev --provider aws
    Dev2->>CLI2: create --app payment --env dev --provider aws
    
    CLI1->>Lock: Acquire Lock (payment-dev-aws)
    Lock-->>CLI1: Lock Acquired
    
    CLI2->>Lock: Acquire Lock (payment-dev-aws)
    Lock-->>CLI2: Lock Held by CLI1
    
    CLI1->>State: Read State (payment-dev-aws)
    CLI1->>State: Write State (payment-dev-aws)
    CLI1->>Lock: Release Lock (payment-dev-aws)
    
    CLI2->>Lock: Retry Acquire Lock
    Lock-->>CLI2: Lock Acquired
    CLI2->>State: Read State (payment-dev-aws)
    CLI2->>State: Write State (payment-dev-aws)
    CLI2->>Lock: Release Lock (payment-dev-aws)
```

## Error Handling Architecture

```mermaid
graph TD
    Operation[CLI Operation] --> Try{Execute}
    
    Try -->|Success| Success[Return Success]
    Try -->|Error| Classify[Classify Error Type]
    
    Classify --> AuthError{Auth Error?}
    Classify --> ValidationError{Validation Error?}
    Classify --> TerraformError{Terraform Error?}
    Classify --> HelmError{Helm Error?}
    Classify --> NetworkError{Network Error?}
    
    AuthError -->|Yes| DisplayAuthHelp[Display Auth Instructions]
    ValidationError -->|Yes| DisplayValidationMsg[Display Validation Message]
    TerraformError -->|Yes| TFRollback[Trigger Terraform Rollback]
    HelmError -->|Yes| HelmRollback[Trigger Helm Rollback]
    NetworkError -->|Yes| RetryLogic[Retry with Backoff]
    
    TFRollback --> LogError[Log Detailed Error]
    HelmRollback --> LogError
    RetryLogic -->|Max Retries| LogError
    RetryLogic -->|Retry| Try
    
    DisplayAuthHelp --> Exit[Exit with Error Code]
    DisplayValidationMsg --> Exit
    LogError --> DisplayError[Display User-Friendly Error]
    DisplayError --> Exit
    
    Success --> ExitSuccess[Exit with Success Code]
```

## Monitoring and Observability

```mermaid
graph LR
    subgraph "CLI Instrumentation"
        Operations[CLI Operations]
        Metrics[Metrics Collector]
        Logs[Log Writer]
    end
    
    subgraph "Local Storage"
        LogFiles[~/.devplatform/logs/]
        MetricsCache[Metrics Cache]
    end
    
    subgraph "AWS Resources"
        CloudWatch[CloudWatch Logs]
        CloudTrail[CloudTrail Events]
        CostExplorer[Cost Explorer]
    end
    
    subgraph "Observability Data"
        ProvisionTime[Provision Duration]
        ErrorRate[Error Rate]
        CostData[Cost Data]
        ResourceCount[Resource Count]
    end
    
    Operations --> Metrics
    Operations --> Logs
    
    Metrics --> MetricsCache
    Logs --> LogFiles
    
    Logs --> CloudWatch
    Operations --> CloudTrail
    
    MetricsCache --> ProvisionTime
    MetricsCache --> ErrorRate
    CloudTrail --> ResourceCount
    CostExplorer --> CostData
```

## Technology Stack

```mermaid
mindmap
  root((DevPlatform CLI))
    Core Language
      Go 1.21+
      Cobra Framework
      Viper Config
    Infrastructure
      Terraform 1.5+
      AWS Provider
      Azure Provider
      Kubernetes Provider
    Container Orchestration
      Helm 3.x
      Kubernetes 1.27+
      AWS EKS
      Azure AKS
    Cloud Provider
      AWS
        VPC
        RDS
        EKS
        S3
        DynamoDB
        Secrets Manager
      Azure
        VNet
        Azure Database
        AKS
        Azure Storage
        Key Vault
    Development Tools
      GitHub Actions
      golangci-lint
      go test
    Dependencies
      aws-sdk-go-v2
      azure-sdk-for-go
      client-go
      terraform-exec
      helm-go
```

## Build and Release Pipeline

```mermaid
graph LR
    subgraph "Development"
        Code[Source Code]
        Test[Unit Tests]
        Lint[Linting]
    end
    
    subgraph "CI Pipeline"
        Build[Build Binary]
        IntegrationTest[Integration Tests]
        SecurityScan[Security Scan]
    end
    
    subgraph "Release Pipeline"
        Tag[Git Tag]
        BuildMatrix[Build Matrix<br/>Linux/Mac/Windows]
        Sign[Sign Binaries]
        Upload[Upload to GitHub]
    end
    
    subgraph "Distribution"
        GHRelease[GitHub Releases]
        Checksums[SHA256 Checksums]
        Docs[Release Notes]
    end
    
    Code --> Test
    Test --> Lint
    Lint --> Build
    Build --> IntegrationTest
    IntegrationTest --> SecurityScan
    
    SecurityScan --> Tag
    Tag --> BuildMatrix
    BuildMatrix --> Sign
    Sign --> Upload
    
    Upload --> GHRelease
    Upload --> Checksums
    Upload --> Docs
```

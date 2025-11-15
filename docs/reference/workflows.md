# DevPlatform CLI - Workflow Documentation

## Create Command Workflow

### Complete Create Flow

```mermaid
flowchart TD
    Start([Developer runs create command]) --> ParseArgs[Parse Command Arguments]
    ParseArgs --> ValidateFlags{All Required<br/>Flags Present?}
    
    ValidateFlags -->|No| ShowUsage[Display Usage Help]
    ValidateFlags -->|Yes| ValidateAppName[Validate App Name]
    
    ShowUsage --> End([Exit with Error])
    
    ValidateAppName --> AppNameValid{Valid Format?}
    AppNameValid -->|No| ShowAppError[Display App Name Error]
    AppNameValid -->|Yes| ValidateEnv[Validate Environment Type]
    
    ShowAppError --> End
    
    ValidateEnv --> EnvValid{Valid Env Type?}
    EnvValid -->|No| ShowEnvError[Display Env Type Error]
    EnvValid -->|Yes| SelectProvider{Cloud Provider?}
    
    ShowEnvError --> End
    
    SelectProvider -->|AWS| CheckAWS[Check AWS Credentials]
    SelectProvider -->|Azure| CheckAzure[Check Azure Credentials]
    
    CheckAWS --> AWSValid{Credentials<br/>Valid?}
    AWSValid -->|No| ShowAWSError[Display AWS Auth Error]
    AWSValid -->|Yes| LoadConfig[Load Configuration File]
    
    CheckAzure --> AzureValid{Credentials<br/>Valid?}
    AzureValid -->|No| ShowAzureError[Display Azure Auth Error]
    AzureValid -->|Yes| LoadConfig
    
    ShowAWSError --> End
    ShowAzureError --> End
    
    LoadConfig --> MergeConfig[Merge CLI Flags with Config]
    MergeConfig --> DryRunCheck{Dry Run<br/>Mode?}
    
    DryRunCheck -->|Yes| TFPlan[Execute Terraform Plan]
    DryRunCheck -->|No| TFInit[Execute Terraform Init]
    
    TFPlan --> DisplayPlan[Display Plan Output]
    TFPlan --> EstimateCost[Display Cost Estimate]
    EstimateCost --> EndSuccess([Exit Success])
    
    TFInit --> TFInitSuccess{Init<br/>Success?}
    TFInitSuccess -->|No| TFInitError[Display Terraform Init Error]
    TFInitSuccess -->|Yes| TFApply[Execute Terraform Apply]
    
    TFInitError --> End
    
    TFApply --> TFApplySuccess{Apply<br/>Success?}
    TFApplySuccess -->|No| TFApplyError[Display Terraform Error]
    TFApplySuccess -->|Yes| ExtractOutputs[Extract Terraform Outputs]
    
    TFApplyError --> TFRollback[Execute Terraform Destroy]
    TFRollback --> LogRollback[Log Rollback Details]
    LogRollback --> End
    
    ExtractOutputs --> GetDBEndpoint[Get Database Endpoint]
    GetDBEndpoint --> GetNetworkID[Get Network ID]
    GetNetworkID --> PrepareHelm[Prepare Helm Values]
    
    PrepareHelm --> HelmInstall[Execute Helm Install]
    HelmInstall --> HelmSuccess{Install<br/>Success?}
    
    HelmSuccess -->|No| HelmError[Display Helm Error]
    HelmSuccess -->|Yes| VerifyPods[Verify Pods Running]
    
    HelmError --> HelmRollback[Execute Helm Uninstall]
    HelmRollback --> TFRollback
    
    VerifyPods --> PodsReady{Pods<br/>Ready?}
    PodsReady -->|No| PodError[Display Pod Status Error]
    PodsReady -->|Yes| GetIngress[Get Ingress URL]
    
    PodError --> HelmRollback
    
    GetIngress --> UpdateKubeconfig[Display Kubeconfig Instructions]
    UpdateKubeconfig --> DisplayOutputs[Display Success Message]
    DisplayOutputs --> ShowDB[Show Database Endpoint]
    ShowDB --> ShowURL[Show Ingress URL]
    ShowURL --> ShowNamespace[Show Namespace Name]
    ShowNamespace --> LogSuccess[Log Success Event]
    LogSuccess --> EndSuccess
    
    DisplayPlan --> EndSuccess
```

### Terraform Execution Detail

```mermaid
sequenceDiagram
    participant CLI as CLI Core
    participant TFWrapper as Terraform Wrapper
    participant TF as Terraform Binary
    participant Backend as State Backend (S3 or Azure Storage)
    participant Cloud as Cloud API (AWS or Azure)
    
    CLI->>TFWrapper: Execute(app, env, provider, action)
    TFWrapper->>TFWrapper: Build Backend Config
    TFWrapper->>TFWrapper: Build Variable Overrides
    
    TFWrapper->>TF: terraform init -backend-config=...
    TF->>Backend: Initialize State Backend
    Backend-->>TF: Backend Ready
    TF-->>TFWrapper: Init Success
    
    TFWrapper->>TF: terraform apply -var app_name=... -var env=... -var provider=...
    TF->>Backend: Acquire State Lock
    Backend-->>TF: Lock Acquired
    
    alt AWS Provider
        TF->>Cloud: Create VPC
        Cloud-->>TF: VPC Created
        TF->>Cloud: Create RDS Instance
        Cloud-->>TF: RDS Instance Creating
    else Azure Provider
        TF->>Cloud: Create VNet
        Cloud-->>TF: VNet Created
        TF->>Cloud: Create Azure Database
        Cloud-->>TF: Azure Database Creating
    end
    
    TF->>Cloud: Create Subnets
    Cloud-->>TF: Subnets Created
    
    TF->>Cloud: Create Security Groups/NSGs
    Cloud-->>TF: Security Groups Created
    
    TF->>Cloud: Wait for Database Ready
    Cloud-->>TF: Database Ready
    
    TF->>Cloud: Create Kubernetes Namespace
    Cloud-->>TF: Namespace Created
    
    TF->>Backend: Write State
    Backend-->>TF: State Written
    
    TF->>Backend: Release Lock
    Backend-->>TF: Lock Released
    
    TF-->>TFWrapper: Apply Complete + Outputs
    TFWrapper->>TFWrapper: Parse Outputs
    TFWrapper-->>CLI: Return Outputs
```

### Helm Deployment Detail

```mermaid
sequenceDiagram
    participant CLI as CLI Core
    participant HelmWrapper as Helm Wrapper
    participant Helm as Helm Binary
    participant K8s as Kubernetes API (EKS or AKS)
    participant Pods as Application Pods
    
    CLI->>HelmWrapper: Install(app, env, values)
    HelmWrapper->>HelmWrapper: Load Base Chart
    HelmWrapper->>HelmWrapper: Merge Values
    HelmWrapper->>HelmWrapper: Validate Chart
    
    HelmWrapper->>Helm: helm upgrade --install app-name ./chart
    Helm->>K8s: Create/Update Release
    
    K8s->>K8s: Create Deployment
    K8s->>K8s: Create Service
    K8s->>K8s: Create Ingress
    K8s->>K8s: Create ConfigMap
    K8s->>K8s: Create ServiceAccount
    
    K8s->>Pods: Schedule Pods
    Pods->>Pods: Pull Image
    Pods->>Pods: Start Container
    Pods-->>K8s: Pod Running
    
    K8s-->>Helm: Release Deployed
    Helm-->>HelmWrapper: Install Success
    
    HelmWrapper->>Helm: helm status app-name
    Helm->>K8s: Get Release Status
    K8s-->>Helm: Status Info
    Helm-->>HelmWrapper: Status Details
    
    HelmWrapper->>HelmWrapper: Verify Pods Ready
    HelmWrapper-->>CLI: Deployment Complete
```

## Status Command Workflow

```mermaid
flowchart TD
    Start([Developer runs status command]) --> ParseArgs[Parse Command Arguments]
    ParseArgs --> ValidateFlags{Required<br/>Flags Present?}
    
    ValidateFlags -->|No| ShowUsage[Display Usage Help]
    ValidateFlags -->|Yes| SelectProvider{Cloud Provider?}
    
    ShowUsage --> End([Exit with Error])
    
    SelectProvider -->|AWS| CheckAWS[Check AWS Credentials]
    SelectProvider -->|Azure| CheckAzure[Check Azure Credentials]
    
    CheckAWS --> AWSValid{Credentials<br/>Valid?}
    AWSValid -->|No| ShowAWSError[Display AWS Auth Error]
    AWSValid -->|Yes| CheckState[Check Terraform State]
    
    CheckAzure --> AzureValid{Credentials<br/>Valid?}
    AzureValid -->|No| ShowAzureError[Display Azure Auth Error]
    AzureValid -->|Yes| CheckState
    
    ShowAWSError --> End
    ShowAzureError --> End
    
    CheckState --> StateExists{State<br/>Exists?}
    StateExists -->|No| ShowNotFound[Display Environment Not Found]
    StateExists -->|Yes| ReadState[Read Terraform State]
    
    ShowNotFound --> End
    
    ReadState --> ExtractResources[Extract Resource IDs]
    ExtractResources --> CheckNetwork[Check Network Status]
    CheckNetwork --> CheckDB[Check Database Status]
    CheckDB --> CheckNamespace[Check Namespace Status]
    
    CheckNamespace --> NamespaceExists{Namespace<br/>Exists?}
    NamespaceExists -->|No| ShowNSError[Mark Namespace as Missing]
    NamespaceExists -->|Yes| GetPods[Get Pod Status]
    
    ShowNSError --> FormatOutput[Format Status Output]
    
    GetPods --> CheckPodHealth[Check Pod Health]
    CheckPodHealth --> GetIngress[Get Ingress Status]
    GetIngress --> FormatOutput
    
    FormatOutput --> DisplayTable[Display Status Table]
    DisplayTable --> ShowNetworkStatus[Network: OK/Error]
    ShowNetworkStatus --> ShowDBStatus[Database: OK/Error]
    ShowDBStatus --> ShowPodStatus[Pods: X/Y Ready]
    ShowPodStatus --> ShowIngressStatus[Ingress: URL or Pending]
    ShowIngressStatus --> EndSuccess([Exit Success])
```

### Status Check Sequence

```mermaid
sequenceDiagram
    participant CLI as CLI Core
    participant TFWrapper as Terraform Wrapper
    participant Backend as State Backend (S3 or Azure Storage)
    participant CloudUtil as Cloud Provider Utilities
    participant K8s as Kubernetes API
    
    CLI->>TFWrapper: GetState(app, env, provider)
    TFWrapper->>Backend: Read State File
    Backend-->>TFWrapper: State Data
    TFWrapper->>TFWrapper: Parse State
    TFWrapper-->>CLI: Resource IDs
    
    CLI->>CloudUtil: CheckNetwork(network_id)
    CloudUtil->>CloudUtil: Cloud API Call
    CloudUtil-->>CLI: Network Status
    
    CLI->>CloudUtil: CheckDatabase(db_id)
    CloudUtil->>CloudUtil: Cloud API Call
    CloudUtil-->>CLI: Database Status
    
    CLI->>K8s: GetNamespace(namespace)
    K8s-->>CLI: Namespace Info
    
    CLI->>K8s: GetPods(namespace)
    K8s-->>CLI: Pod List
    
    CLI->>K8s: GetIngress(namespace)
    K8s-->>CLI: Ingress Info
    
    CLI->>CLI: Aggregate Status
    CLI->>CLI: Format Table
    CLI->>CLI: Display Output
```

## Destroy Command Workflow

```mermaid
flowchart TD
    Start([Developer runs destroy command]) --> ParseArgs[Parse Command Arguments]
    ParseArgs --> ValidateFlags{Required<br/>Flags Present?}
    
    ValidateFlags -->|No| ShowUsage[Display Usage Help]
    ValidateFlags -->|Yes| CheckConfirm{Confirm<br/>Flag Present?}
    
    ShowUsage --> End([Exit with Error])
    
    CheckConfirm -->|No| PromptConfirm[Prompt for Confirmation]
    CheckConfirm -->|Yes| SelectProvider{Cloud Provider?}
    
    PromptConfirm --> UserConfirms{User<br/>Confirms?}
    UserConfirms -->|No| ShowCancelled[Display Cancelled Message]
    UserConfirms -->|Yes| SelectProvider
    
    ShowCancelled --> End
    
    SelectProvider -->|AWS| CheckAWS[Check AWS Credentials]
    SelectProvider -->|Azure| CheckAzure[Check Azure Credentials]
    
    CheckAWS --> AWSValid{Credentials<br/>Valid?}
    AWSValid -->|No| ShowAWSError[Display AWS Auth Error]
    AWSValid -->|Yes| CheckState[Check Terraform State]
    
    CheckAzure --> AzureValid{Credentials<br/>Valid?}
    AzureValid -->|No| ShowAzureError[Display Azure Auth Error]
    AzureValid -->|Yes| CheckState
    
    ShowAWSError --> End
    ShowAzureError --> End
    
    CheckState --> StateExists{State<br/>Exists?}
    StateExists -->|No| ShowNotFound[Display Environment Not Found]
    StateExists -->|Yes| GetNamespace[Get Namespace Name]
    
    ShowNotFound --> End
    
    GetNamespace --> CheckHelm[Check Helm Release]
    CheckHelm --> HelmExists{Release<br/>Exists?}
    
    HelmExists -->|No| SkipHelm[Skip Helm Uninstall]
    HelmExists -->|Yes| HelmUninstall[Execute Helm Uninstall]
    
    SkipHelm --> TFDestroy[Execute Terraform Destroy]
    
    HelmUninstall --> HelmSuccess{Uninstall<br/>Success?}
    HelmSuccess -->|No| HelmError[Display Helm Error]
    HelmSuccess -->|Yes| VerifyPodsGone[Verify Pods Deleted]
    
    HelmError --> ContinueAnyway{Continue<br/>Anyway?}
    ContinueAnyway -->|No| End
    ContinueAnyway -->|Yes| TFDestroy
    
    VerifyPodsGone --> TFDestroy
    
    TFDestroy --> TFDestroySuccess{Destroy<br/>Success?}
    TFDestroySuccess -->|No| TFDestroyError[Display Terraform Error]
    TFDestroySuccess -->|Yes| VerifyStateGone[Verify Resources Deleted]
    
    TFDestroyError --> ShowManualCleanup[Display Manual Cleanup Instructions]
    ShowManualCleanup --> End
    
    VerifyStateGone --> CalculateSavings[Calculate Cost Savings]
    CalculateSavings --> DisplaySuccess[Display Success Message]
    DisplaySuccess --> ShowSavings[Show Estimated Savings]
    ShowSavings --> LogDestroy[Log Destroy Event]
    LogDestroy --> EndSuccess([Exit Success])
```

### Destroy Sequence with Rollback Handling

```mermaid
sequenceDiagram
    participant CLI as CLI Core
    participant HelmWrapper as Helm Wrapper
    participant K8s as Kubernetes API
    participant TFWrapper as Terraform Wrapper
    participant Cloud as Cloud API (AWS or Azure)
    participant Backend as State Backend
    
    CLI->>CLI: Prompt for Confirmation
    CLI->>HelmWrapper: Uninstall(app, env)
    
    HelmWrapper->>K8s: Delete Release
    K8s->>K8s: Delete Deployment
    K8s->>K8s: Delete Service
    K8s->>K8s: Delete Ingress
    K8s->>K8s: Delete Pods
    K8s-->>HelmWrapper: Release Deleted
    HelmWrapper-->>CLI: Uninstall Success
    
    CLI->>TFWrapper: Destroy(app, env, provider)
    TFWrapper->>Backend: Acquire State Lock
    Backend-->>TFWrapper: Lock Acquired
    
    TFWrapper->>Cloud: Delete Kubernetes Namespace
    Cloud-->>TFWrapper: Namespace Deleted
    
    TFWrapper->>Cloud: Delete Database Instance
    Cloud-->>TFWrapper: Database Deleting
    
    TFWrapper->>Cloud: Wait for Database Deletion
    Cloud-->>TFWrapper: Database Deleted
    
    TFWrapper->>Cloud: Delete Security Groups/NSGs
    Cloud-->>TFWrapper: Security Groups Deleted
    
    TFWrapper->>Cloud: Delete Subnets
    Cloud-->>TFWrapper: Subnets Deleted
    
    TFWrapper->>Cloud: Delete Network (VPC/VNet)
    Cloud-->>TFWrapper: Network Deleted
    
    TFWrapper->>Backend: Delete State File
    Backend-->>TFWrapper: State Deleted
    
    TFWrapper->>Backend: Release Lock
    Backend-->>TFWrapper: Lock Released
    
    TFWrapper-->>CLI: Destroy Complete
    CLI->>CLI: Calculate Savings
    CLI->>CLI: Display Success
```

## Configuration Loading Workflow

```mermaid
flowchart TD
    Start([CLI Starts]) --> CheckConfigFile{Config File<br/>Exists?}
    
    CheckConfigFile -->|No| UseDefaults[Use Default Values]
    CheckConfigFile -->|Yes| ReadFile[Read .devplatform.yaml]
    
    UseDefaults --> ParseCLIFlags[Parse CLI Flags]
    
    ReadFile --> ValidateYAML{Valid<br/>YAML?}
    ValidateYAML -->|No| ShowYAMLError[Display YAML Parse Error]
    ValidateYAML -->|Yes| ValidateSchema{Valid<br/>Schema?}
    
    ShowYAMLError --> End([Exit with Error])
    
    ValidateSchema -->|No| ShowSchemaError[Display Schema Error]
    ValidateSchema -->|Yes| LoadConfigValues[Load Config Values]
    
    ShowSchemaError --> End
    
    LoadConfigValues --> ParseCLIFlags
    ParseCLIFlags --> MergeConfig[Merge Config with CLI Flags]
    MergeConfig --> PrioritizeCLI[CLI Flags Override Config]
    PrioritizeCLI --> ValidateFinal[Validate Final Config]
    ValidateFinal --> ConfigReady([Config Ready])
```

## Error Handling and Rollback Workflow

```mermaid
flowchart TD
    Operation[CLI Operation] --> Execute{Execute<br/>Operation}
    
    Execute -->|Success| LogSuccess[Log Success]
    Execute -->|Error| CaptureError[Capture Error Details]
    
    LogSuccess --> ReturnSuccess([Return Success])
    
    CaptureError --> ClassifyError{Error<br/>Type?}
    
    ClassifyError -->|Auth Error| DisplayAuthHelp[Display Auth Instructions]
    ClassifyError -->|Validation Error| DisplayValidation[Display Validation Message]
    ClassifyError -->|Terraform Error| HandleTFError[Handle Terraform Error]
    ClassifyError -->|Helm Error| HandleHelmError[Handle Helm Error]
    ClassifyError -->|Network Error| HandleNetworkError[Handle Network Error]
    
    DisplayAuthHelp --> LogError[Log Error Details]
    DisplayValidation --> LogError
    
    HandleTFError --> TFRollbackNeeded{Rollback<br/>Needed?}
    TFRollbackNeeded -->|Yes| ExecuteTFRollback[Execute Terraform Destroy]
    TFRollbackNeeded -->|No| LogError
    
    ExecuteTFRollback --> TFRollbackSuccess{Rollback<br/>Success?}
    TFRollbackSuccess -->|Yes| LogRollbackSuccess[Log Rollback Success]
    TFRollbackSuccess -->|No| LogRollbackFailure[Log Rollback Failure]
    
    LogRollbackSuccess --> LogError
    LogRollbackFailure --> ShowManualCleanup[Show Manual Cleanup Instructions]
    ShowManualCleanup --> LogError
    
    HandleHelmError --> HelmRollbackNeeded{Rollback<br/>Needed?}
    HelmRollbackNeeded -->|Yes| ExecuteHelmRollback[Execute Helm Uninstall]
    HelmRollbackNeeded -->|No| LogError
    
    ExecuteHelmRollback --> HelmRollbackSuccess{Rollback<br/>Success?}
    HelmRollbackSuccess -->|Yes| LogRollbackSuccess
    HelmRollbackSuccess -->|No| LogRollbackFailure
    
    HandleNetworkError --> RetryPossible{Retry<br/>Possible?}
    RetryPossible -->|Yes| WaitBackoff[Wait with Exponential Backoff]
    RetryPossible -->|No| LogError
    
    WaitBackoff --> CheckRetries{Max Retries<br/>Reached?}
    CheckRetries -->|No| Execute
    CheckRetries -->|Yes| LogError
    
    LogError --> DisplayUserError[Display User-Friendly Error]
    DisplayUserError --> ReturnError([Return Error])
```

## Concurrent Execution Workflow

```mermaid
flowchart TD
    Dev1Start([Developer 1: create payment-dev]) --> CLI1[CLI Instance 1]
    Dev2Start([Developer 2: create payment-dev]) --> CLI2[CLI Instance 2]
    
    CLI1 --> TF1Init[Terraform Init]
    CLI2 --> TF2Init[Terraform Init]
    
    TF1Init --> TF1Lock[Acquire Lock: payment-dev]
    TF2Init --> TF2Lock[Acquire Lock: payment-dev]
    
    TF1Lock --> Lock1Success{Lock<br/>Acquired?}
    Lock1Success -->|Yes| TF1Apply[Terraform Apply]
    
    TF2Lock --> Lock2Success{Lock<br/>Acquired?}
    Lock2Success -->|No| TF2Wait[Wait for Lock]
    
    TF1Apply --> TF1Complete[Complete Operation]
    TF1Complete --> TF1Release[Release Lock]
    
    TF1Release --> TF2Retry[CLI2 Retries Lock]
    TF2Wait --> TF2Retry
    
    TF2Retry --> Lock2Retry{Lock<br/>Acquired?}
    Lock2Retry -->|Yes| TF2Apply[Terraform Apply]
    Lock2Retry -->|No| TF2Timeout{Timeout<br/>Reached?}
    
    TF2Timeout -->|Yes| TF2Error[Display Lock Timeout Error]
    TF2Timeout -->|No| TF2Wait
    
    TF2Apply --> TF2Complete[Complete Operation]
    TF2Complete --> TF2Release[Release Lock]
    
    TF2Error --> End2([Exit with Error])
    TF2Release --> End2Success([Exit Success])
    TF1Release --> End1([Exit Success])
```

## Version Check Workflow

```mermaid
flowchart TD
    Start([CLI Starts]) --> CheckTerraform[Check Terraform Binary]
    CheckTerraform --> TFExists{Terraform<br/>Found?}
    
    TFExists -->|No| ShowTFError[Display Terraform Not Found Error]
    TFExists -->|Yes| CheckTFVersion[Check Terraform Version]
    
    ShowTFError --> End([Exit with Error])
    
    CheckTFVersion --> TFVersionOK{Version >= 1.5?}
    TFVersionOK -->|No| ShowTFVersionError[Display Terraform Version Error]
    TFVersionOK -->|Yes| CheckHelm[Check Helm Binary]
    
    ShowTFVersionError --> End
    
    CheckHelm --> HelmExists{Helm<br/>Found?}
    HelmExists -->|No| ShowHelmError[Display Helm Not Found Error]
    HelmExists -->|Yes| CheckHelmVersion[Check Helm Version]
    
    ShowHelmError --> End
    
    CheckHelmVersion --> HelmVersionOK{Version >= 3.0?}
    HelmVersionOK -->|No| ShowHelmVersionError[Display Helm Version Error]
    HelmVersionOK -->|Yes| CheckKubectl[Check Kubectl Binary]
    
    ShowHelmVersionError --> End
    
    CheckKubectl --> KubectlExists{Kubectl<br/>Found?}
    KubectlExists -->|No| ShowKubectlError[Display Kubectl Not Found Error]
    KubectlExists -->|Yes| CheckCloudCLI{Check Cloud CLI}
    
    ShowKubectlError --> End
    
    CheckCloudCLI -->|AWS| CheckAWSCLI[Check AWS CLI]
    CheckCloudCLI -->|Azure| CheckAzureCLI[Check Azure CLI]
    CheckCloudCLI -->|Both| CheckAWSCLI
    
    CheckAWSCLI --> AWSCLIExists{AWS CLI<br/>Found?}
    AWSCLIExists -->|No| ShowAWSCLIError[Display AWS CLI Not Found Error]
    AWSCLIExists -->|Yes| CheckAzureCLI
    
    CheckAzureCLI --> AzureCLIExists{Azure CLI<br/>Found?}
    AzureCLIExists -->|No| ShowAzureCLIError[Display Azure CLI Not Found Error]
    AzureCLIExists -->|Yes| AllToolsReady([All Tools Ready])
    
    ShowAWSCLIError --> End
    ShowAzureCLIError --> End
```

## Logging Workflow

```mermaid
flowchart TD
    Event[CLI Event] --> DetermineLevel{Log<br/>Level?}
    
    DetermineLevel -->|Debug| DebugEnabled{Debug<br/>Enabled?}
    DetermineLevel -->|Info| WriteInfo[Write Info Log]
    DetermineLevel -->|Warning| WriteWarning[Write Warning Log]
    DetermineLevel -->|Error| WriteError[Write Error Log]
    
    DebugEnabled -->|Yes| WriteDebug[Write Debug Log]
    DebugEnabled -->|No| Skip[Skip Log]
    
    WriteDebug --> FormatLog[Format Log Entry]
    WriteInfo --> FormatLog
    WriteWarning --> FormatLog
    WriteError --> FormatLog
    
    FormatLog --> AddTimestamp[Add Timestamp]
    AddTimestamp --> AddContext[Add Context Info]
    AddContext --> WriteToFile[Write to Log File]
    
    WriteToFile --> CheckRotation{Log File<br/>Too Large?}
    CheckRotation -->|Yes| RotateLog[Rotate Log File]
    CheckRotation -->|No| CheckConsole{Console<br/>Output?}
    
    RotateLog --> CheckConsole
    
    CheckConsole -->|Yes| WriteConsole[Write to Console]
    CheckConsole -->|No| Done([Done])
    
    WriteConsole --> ColorizeOutput{Colorize?}
    ColorizeOutput -->|Yes| ApplyColors[Apply ANSI Colors]
    ColorizeOutput -->|No| Done
    
    ApplyColors --> Done
    Skip --> Done
```

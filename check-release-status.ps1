# Release Status Checker for v1.0.0
# Checks if the GitHub release is live and accessible

param(
    [string]$Version = "v1.0.0",
    [string]$Owner = "vedantchimote",
    [string]$Repo = "DevPlatform-CLI"
)

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "Release Status Checker" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

Write-Host "Checking release: $Version" -ForegroundColor White
Write-Host "Repository: $Owner/$Repo`n" -ForegroundColor White

# Check 1: Verify tag exists locally
Write-Host "1. Checking local tag..." -ForegroundColor Yellow
$localTag = git tag -l $Version
if ($localTag) {
    Write-Host "   ✓ Tag $Version exists locally" -ForegroundColor Green
} else {
    Write-Host "   ✗ Tag $Version not found locally" -ForegroundColor Red
}

# Check 2: Verify tag exists on remote
Write-Host "`n2. Checking remote tag..." -ForegroundColor Yellow
$remoteTags = git ls-remote --tags origin
if ($remoteTags -match $Version) {
    Write-Host "   ✓ Tag $Version exists on remote" -ForegroundColor Green
} else {
    Write-Host "   ✗ Tag $Version not found on remote" -ForegroundColor Red
}

# Check 3: Test release page accessibility
Write-Host "`n3. Checking release page..." -ForegroundColor Yellow
$releaseUrl = "https://github.com/$Owner/$Repo/releases/tag/$Version"
try {
    $response = Invoke-WebRequest -Uri $releaseUrl -Method Head -UseBasicParsing -ErrorAction Stop
    if ($response.StatusCode -eq 200) {
        Write-Host "   ✓ Release page is accessible" -ForegroundColor Green
        Write-Host "   URL: $releaseUrl" -ForegroundColor Gray
    }
} catch {
    Write-Host "   ⏳ Release page not yet available (workflow may still be running)" -ForegroundColor Yellow
    Write-Host "   URL: $releaseUrl" -ForegroundColor Gray
}

# Check 4: Test checksums file
Write-Host "`n4. Checking release assets..." -ForegroundColor Yellow
$checksumsUrl = "https://github.com/$Owner/$Repo/releases/download/$Version/checksums.txt"
try {
    $response = Invoke-WebRequest -Uri $checksumsUrl -Method Head -UseBasicParsing -ErrorAction Stop
    if ($response.StatusCode -eq 200) {
        Write-Host "   ✓ Release assets are available" -ForegroundColor Green
        Write-Host "   Checksums: $checksumsUrl" -ForegroundColor Gray
    }
} catch {
    Write-Host "   ⏳ Release assets not yet available" -ForegroundColor Yellow
    Write-Host "   This is normal if the workflow is still running" -ForegroundColor Gray
}

# Check 5: GitHub Actions status
Write-Host "`n5. GitHub Actions workflow..." -ForegroundColor Yellow
$actionsUrl = "https://github.com/$Owner/$Repo/actions"
Write-Host "   Visit: $actionsUrl" -ForegroundColor Gray
Write-Host "   Look for 'Release' workflow triggered by tag $Version" -ForegroundColor Gray

# Summary
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "Summary" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

Write-Host "Quick Links:" -ForegroundColor White
Write-Host "  • Release Page: https://github.com/$Owner/$Repo/releases/tag/$Version" -ForegroundColor Gray
Write-Host "  • Actions: https://github.com/$Owner/$Repo/actions" -ForegroundColor Gray
Write-Host "  • Workflow: https://github.com/$Owner/$Repo/actions/workflows/release.yml" -ForegroundColor Gray

Write-Host "`nExpected Assets:" -ForegroundColor White
Write-Host "  • devplatform-cli_${Version}_Linux_x86_64.tar.gz" -ForegroundColor Gray
Write-Host "  • devplatform-cli_${Version}_Linux_arm64.tar.gz" -ForegroundColor Gray
Write-Host "  • devplatform-cli_${Version}_Darwin_x86_64.tar.gz" -ForegroundColor Gray
Write-Host "  • devplatform-cli_${Version}_Darwin_arm64.tar.gz" -ForegroundColor Gray
Write-Host "  • devplatform-cli_${Version}_Windows_x86_64.zip" -ForegroundColor Gray
Write-Host "  • checksums.txt" -ForegroundColor Gray
Write-Host "  • .deb and .rpm packages" -ForegroundColor Gray

Write-Host "`nNext Steps:" -ForegroundColor White
Write-Host "  1. Visit the Actions page to monitor workflow progress" -ForegroundColor Gray
Write-Host "  2. Once complete, verify all assets are present on the release page" -ForegroundColor Gray
Write-Host "  3. Download and test a binary to ensure it works" -ForegroundColor Gray
Write-Host "  4. Verify checksums match" -ForegroundColor Gray

Write-Host "`nNote: The workflow typically takes 10-15 minutes to complete.`n" -ForegroundColor Yellow


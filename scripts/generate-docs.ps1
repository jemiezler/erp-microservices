# Generate Swagger Documentation for ERP Microservices

# Check if swag is installed
if (!(Get-Command swag -ErrorAction SilentlyContinue)) {
    Write-Host "Error: 'swag' is not installed. Please install it using:" -ForegroundColor Red
    Write-Host "go install github.com/swaggo/swag/cmd/swag@latest" -ForegroundColor Cyan
    exit 1
}

# Use PSScriptRoot to get the directory where the script is located
$ScriptDir = $PSScriptRoot
$RootDir = (Get-Item $ScriptDir).Parent.FullName
$ServicesDir = Join-Path $RootDir "apis\services"

# Store original location to return to it later
$OriginalDir = Get-Location

try {
    # HR Service
    $HRDir = Join-Path $ServicesDir "hr-service"
    Write-Host "Generating docs for HR Service in $HRDir..." -ForegroundColor Green
    Set-Location $HRDir
    swag init -g cmd/api/main.go --parseDependency --parseInternal

    # Finance Service
    $FinanceDir = Join-Path $ServicesDir "finance-service"
    Write-Host "Generating docs for Finance Service in $FinanceDir..." -ForegroundColor Green
    Set-Location $FinanceDir
    swag init -g cmd/api/main.go --parseDependency --parseInternal
}
finally {
    # Always return to original directory
    Set-Location $OriginalDir
}

Write-Host "Done!" -ForegroundColor Green

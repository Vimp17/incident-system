#!/usr/bin/env pwsh

# Incident System Health Check Script
# Run: .\scripts\check-health-en.ps1

param(
    [string]$BaseUrl = "http://localhost:8080",
    [string]$ApiKey = "operator-key-secure-change-me"
)

# Set UTF-8 encoding
$OutputEncoding = [System.Text.Encoding]::UTF8
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Incident System Health Check" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

function Invoke-HealthCheck {
    param(
        [string]$Name,
        [string]$Uri,
        [string]$Method = "GET",
        [hashtable]$Headers = @{},
        [string]$Body = $null
    )
    
    try {
        Write-Host "[$Name] Sending request..." -ForegroundColor Gray
        
        $params = @{
            Uri = $Uri
            Method = $Method
            ContentType = "application/json"
            ErrorAction = "Stop"
        }
        
        if ($Headers.Count -gt 0) {
            $params.Headers = $Headers
        }
        
        if ($Body) {
            $params.Body = $Body
        }
        
        $response = Invoke-RestMethod @params
        
        Write-Host "[$Name] ✅ Success" -ForegroundColor Green
        if ($response) {
            Write-Host "   Response: " -NoNewline -ForegroundColor Gray
            $response | ConvertTo-Json -Depth 5 | Write-Host -ForegroundColor White
        }
        return $true
    }
    catch {
        Write-Host "[$Name] ❌ Error: $_" -ForegroundColor Red
        if ($_.Exception.Response) {
            $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
            $reader.BaseStream.Position = 0
            $reader.DiscardBufferedData()
            $errorBody = $reader.ReadToEnd()
            Write-Host "   Error body: $errorBody" -ForegroundColor Red
        }
        return $false
    }
}

# Test 1: Health Check
Write-Host "1. Testing Health Check API..." -ForegroundColor Yellow
$healthResult = Invoke-HealthCheck -Name "Health" -Uri "$BaseUrl/api/v1/system/health"
if (-not $healthResult) {
    Write-Host "❌ Server not responding. Make sure server is running." -ForegroundColor Red
    Write-Host "   Run: go run cmd/server/main.go" -ForegroundColor Yellow
    exit 1
}

Write-Host ""
Write-Host "2. Creating test incident..." -ForegroundColor Yellow
$incidentBody = @{
    user_id = "health_check_operator"
    latitude = 55.7558
    longitude = 37.6173
    title = "Test incident for health check"
    description = "Created by automatic health check script"
    severity = "medium"
    radius = 500
} | ConvertTo-Json -Compress

$incidentResult = Invoke-HealthCheck `
    -Name "Create Incident" `
    -Uri "$BaseUrl/api/v1/incidents" `
    -Method "POST" `
    -Headers @{"X-API-Key" = $ApiKey} `
    -Body $incidentBody

Write-Host ""
Write-Host "3. Testing location check..." -ForegroundColor Yellow
$locationBody = @{
    user_id = "health_check_user"
    latitude = 55.7558
    longitude = 37.6173
} | ConvertTo-Json -Compress

Invoke-HealthCheck `
    -Name "Location Check" `
    -Uri "$BaseUrl/api/v1/location/check" `
    -Method "POST" `
    -Body $locationBody

Write-Host ""
Write-Host "4. Getting incidents list..." -ForegroundColor Yellow
Invoke-HealthCheck `
    -Name "List Incidents" `
    -Uri "$BaseUrl/api/v1/incidents" `
    -Method "GET" `
    -Headers @{"X-API-Key" = $ApiKey}

Write-Host ""
Write-Host "5. Getting statistics..." -ForegroundColor Yellow
Invoke-HealthCheck `
    -Name "Statistics" `
    -Uri "$BaseUrl/api/v1/incidents/stats?minutes=5" `
    -Method "GET" `
    -Headers @{"X-API-Key" = $ApiKey}

Write-Host ""
Write-Host "6. Checking Docker containers..." -ForegroundColor Yellow
try {
    $containers = docker ps --format "{{.Names}}" 2>$null
    
    if ($LASTEXITCODE -eq 0 -and $containers) {
        Write-Host "   Running containers:" -ForegroundColor Green
        $containers | ForEach-Object { Write-Host "   - $_" -ForegroundColor White }
    }
    else {
        Write-Host "   Docker not available or no containers running" -ForegroundColor Yellow
    }
}
catch {
    Write-Host "   Could not check Docker: $_" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Health check completed!" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

Write-Host ""
Write-Host "Recommendations:" -ForegroundColor Yellow
Write-Host "1. Start server: go run cmd/server/main.go" -ForegroundColor Gray
Write-Host "2. Start Docker services: docker-compose up -d postgres redis" -ForegroundColor Gray
Write-Host "3. View logs: docker-compose logs -f" -ForegroundColor Gray
Write-Host "4. Test API with: make test-api" -ForegroundColor Gray
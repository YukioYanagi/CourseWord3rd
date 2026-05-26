# Gateway (Go) + Python transform health check.
# Run from project root:
#   powershell -ExecutionPolicy Bypass -File .\examples\check.ps1
# Optional:
#   $env:GATEWAY = "http://127.0.0.1:8090"; powershell -File .\examples\check.ps1

$ErrorActionPreference = "Stop"
$Gateway = if ($env:GATEWAY) { $env:GATEWAY.TrimEnd("/") } else { "http://127.0.0.1:8080" }
$Python = if ($env:PYTHON_HEALTH) { $env:PYTHON_HEALTH } else { "http://127.0.0.1:5000/health" }

$here = Split-Path -Parent $MyInvocation.MyCommand.Path
$root = Resolve-Path (Join-Path $here "..")
$sampleJson = Join-Path $here "sample.json"

Write-Host "== Python ==" -ForegroundColor Cyan
try {
    $py = Invoke-RestMethod -Uri $Python -Method Get -TimeoutSec 5
    $py | ConvertTo-Json -Compress
} catch {
    Write-Host "Python not reachable ($Python). Start: python -m uvicorn app:app --host 127.0.0.1 --port 5000" -ForegroundColor Yellow
    throw
}

Write-Host "`n== Go health ($Gateway/api/v1/health) ==" -ForegroundColor Cyan
try {
    $r = Invoke-WebRequest -Uri "$Gateway/api/v1/health" -Method Get -TimeoutSec 5
    Write-Host "X-API-Version:" $r.Headers["X-Api-Version"]
    Write-Host $r.Content
} catch {
    Write-Host "Gateway not reachable. Check port or set GATEWAY env." -ForegroundColor Yellow
    throw
}

Write-Host "`n== POST sample.json (format=json) ==" -ForegroundColor Cyan
if (-not (Test-Path $sampleJson)) { throw "Missing file: $sampleJson" }

$boundary = [Guid]::NewGuid().ToString("N")
$LF = "`r`n"
$bodyLines = @(
    "--$boundary",
    'Content-Disposition: form-data; name="format"',
    "",
    "json",
    "--$boundary",
    'Content-Disposition: form-data; name="file"; filename="sample.json"',
    "Content-Type: application/json",
    "",
    ([System.IO.File]::ReadAllText($sampleJson, [System.Text.UTF8Encoding]::new($false))),
    "--$boundary--"
)
$bodyText = $bodyLines -join $LF
$bytes = [System.Text.Encoding]::UTF8.GetBytes($bodyText)
$resp = Invoke-WebRequest -Uri "$Gateway/api/v1/send" -Method Post -ContentType "multipart/form-data; boundary=$boundary" -Body $bytes -TimeoutSec 30
Write-Host $resp.Content

Write-Host "`n== GET /api/v1/received ==" -ForegroundColor Cyan
(Invoke-RestMethod -Uri "$Gateway/api/v1/received" -Method Get) | ConvertTo-Json -Depth 5

Write-Host "`nDone. Web UI: $Gateway/" -ForegroundColor Green

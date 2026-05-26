# Local security checks (Windows). Run from repo root.
$ErrorActionPreference = "Stop"
Set-Location (Split-Path -Parent $PSScriptRoot)

Write-Host "== go mod verify ==" -ForegroundColor Cyan
go mod verify

Write-Host "== go vet ==" -ForegroundColor Cyan
go vet ./...

Write-Host "== gosec ==" -ForegroundColor Cyan
$gobin = Join-Path (go env GOPATH) "bin"
$gosec = Join-Path $gobin "gosec.exe"
if (-not (Test-Path $gosec)) {
    go install github.com/securego/gosec/v2/cmd/gosec@latest
}
& $gosec -conf .gosec.json ./...

Write-Host "== pip-audit (python) ==" -ForegroundColor Cyan
python -m pip install -q pip-audit
python -m pip_audit -r python/requirements.txt

Write-Host "== bandit (python) ==" -ForegroundColor Cyan
python -m pip install -q "bandit[toml]"
python -m bandit -r python/app.py -ll

Write-Host "Done." -ForegroundColor Green

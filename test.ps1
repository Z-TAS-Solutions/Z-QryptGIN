$PiUser = "Z-TAS"
$PiHost = "104.43.91.57"
$RemotePath = "/home/Z-TAS-void/"

Write-Host "Starting Build for Linux (arm64)..." -ForegroundColor Cyan

$env:GOOS = "linux"
$env:GOARCH = "arm64"

go build -ldflags="-s -w" -o zqrypt cmd/api/main.go

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build Failed!" -ForegroundColor Red
    exit
}

Write-Host "Uploading to Pi..." -ForegroundColor Yellow

scp zqrypt "${PiUser}@${PiHost}:${RemotePath}/zqrypt"

Write-Host "Setting executable permissions..." -ForegroundColor Blue
ssh "${PiUser}@${PiHost}" "chmod +x ${RemotePath}/zqrypt"

Write-Host "Done! go run ./zqrypt on your server peasant !." -ForegroundColor Green

$env:GOOS = "windows"
$env:GOARCH = "amd64"


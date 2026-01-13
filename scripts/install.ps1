$ErrorActionPreference = "Stop"

$installDir = "$env:LOCALAPPDATA\git-radar"
$binaryName = "git-radar.exe"

if (Test-Path ".\$binaryName") {
    $binaryPath = ".\$binaryName"
} elseif (Test-Path ".\dist\git-radar-windows-amd64.exe") {
    $binaryPath = ".\dist\git-radar-windows-amd64.exe"
} else {
    Write-Error "Error: Binary not found. Download or build first."
    exit 1
}

Write-Host "Installing git-radar to $installDir..."

New-Item -ItemType Directory -Force -Path $installDir | Out-Null
Copy-Item $binaryPath "$installDir\$binaryName"

$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
    Write-Host "Added $installDir to PATH"
}

Write-Host ""
Write-Host "git-radar installed successfully!" -ForegroundColor Green
Write-Host "Restart your terminal, then run 'git-radar' from any git repository."

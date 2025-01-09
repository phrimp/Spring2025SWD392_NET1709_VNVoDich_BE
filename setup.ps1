# powershell -ExecutionPolicy Bypass -File script.ps1
# Ensure Script is Run as Administrator
if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Host "Please run this script as Administrator." -ForegroundColor Red
    pause
    exit
}

# Update System and Install Required Software
Write-Host "Starting Docker setup..." -ForegroundColor Green

# Check if Chocolatey is installed, install if not
if (!(Get-Command choco -ErrorAction SilentlyContinue)) {
    Write-Host "Chocolatey is not installed. Installing Chocolatey..."
    Set-ExecutionPolicy Bypass -Scope Process -Force
    [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
    Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))
    $env:Path += ";$env:ALLUSERSPROFILE\chocolatey\bin"
} else {
    Write-Host "Chocolatey is already installed."
}

# Install Docker Desktop if not installed
if (!(Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Host "Docker Desktop is not installed. Installing Docker Desktop..."
    choco install docker-desktop -y
    Write-Host "Docker Desktop installed. Please reboot the system and rerun this script." -ForegroundColor Yellow
    pause
    exit
} else {
    Write-Host "Docker Desktop is already installed."
}

# Start Docker Desktop
Write-Host "Starting Docker Desktop..."
try {
    Write-Host "Waiting for 15 seconds"
    Start-Process "C:\Program Files\Docker\Docker\Docker Desktop.exe" -Verb RunAs
    Start-Sleep -Seconds 15  # Allow time for Docker to initialize
} catch {
    Write-Host "Docker installation failed. Exit"
    pause
    exit
}

# Verify Docker Installation
Write-Host "Verifying Docker installation..."
try {
    docker --version
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Docker installation verification failed. Please check Docker installation manually." -ForegroundColor Red
        pause
        exit
    } else {
        Write-Host "Docker is installed and working correctly."
    }
} catch {
    Write-Host "Docker command failed. Ensure Docker Desktop is installed and running." -ForegroundColor Red
    pause
    exit
}

# Install Docker Compose
Write-Host "Installing Docker Compose..."
try {
    choco install docker-compose -y
    Write-Host "Docker Compose installed successfully."
} catch {
    Write-Host "Failed to install Docker Compose. Please check your setup." -ForegroundColor Red
    pause
    exit
}

# Prepare the Application
Write-Host "Preparing application setup..."
$appDir = "$env:USERPROFILE\docker-app"
if (!(Test-Path $appDir)) {
    mkdir $appDir
}

# Copy the docker-compose.yml and related files to the application directory
Write-Host "Copying application files to $appDir..."
$sourceDir = Read-Host "Enter the path to your application files (docker-compose.yml, etc.)"
if (!(Test-Path $sourceDir)) {
    Write-Host "Invalid source directory. Please ensure the path exists." -ForegroundColor Red
    pause
    exit
}
Copy-Item -Path "$sourceDir\*" -Destination $appDir -Recurse -Force

# Navigate to Application Directory and Start Docker Compose
Set-Location -Path $appDir
Write-Host "Starting application with Docker Compose..."
try {
    docker-compose up -d
    Write-Host "Application started successfully."
} catch {
    Write-Host "Failed to start the application. Please check the docker-compose.yml file." -ForegroundColor Red
    pause
    exit
}

# Verify Services
Write-Host "Verifying running containers..."
docker ps

Write-Host "Setup complete! Your application is running." -ForegroundColor Green
pause

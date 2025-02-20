if (Test-Path .env) {
    Get-Content .env | ForEach-Object {
        if ($_ -match '^([^=]+)=(.*)$') {
            [Environment]::SetEnvironmentVariable($matches[1], $matches[2])
        }
    }
}

$REPO_PATH = if ($env:GITHUB_REPO_PATH) { $env:GITHUB_REPO_PATH } else { "https://github.com/phrimp/Spring2025SWD392_NET1709_VNVoDich_BE" }
$BRANCH = if ($env:GITHUB_BRANCH) { $env:GITHUB_BRANCH } else { "master" }
$CHECK_INTERVAL = if ($env:CHECK_INTERVAL) { $env:CHECK_INTERVAL } else { 60 }
$GITHUB_TOKEN = if ($env:GITHUB_TOKEN) { $env:GITHUB_TOKEN } else { "" }
$DOCKER_COMPOSE_FILE = if ($env:DOCKER_COMPOSE_FILE) { $env:DOCKER_COMPOSE_FILE } else { "./docker-compose.yml" }

$dockerProcess = $null

# Function to get the latest commit hash
function Get-LatestCommit {
    if ($GITHUB_TOKEN) {
        $repoPath = $REPO_PATH -replace "https://github.com/", ""
        $url = "https://${GITHUB_TOKEN}@github.com/${repoPath}"
    } else {
        $url = "origin"
    }
    $result = git ls-remote $url -h "refs/heads/$BRANCH"
    if ($result) {
        return $result.Split()[0]
    }
    return $null
}

# Function to log messages with timestamp
function Write-Log {
    param([string]$message)
    Write-Output "[$([DateTime]::Now.ToString('yyyy-MM-dd HH:mm:ss'))] $message"
}

# Function to start Docker
function Start-DockerProcess {
    Write-Log "Starting Docker services..."
    docker-compose -f $DOCKER_COMPOSE_FILE up --build -d
    if ($LASTEXITCODE -ne 0) {
        Write-Log "Error: Docker-compose failed to start"
        return $false
    }
    return $true
}

# Function to stop Docker process
function Stop-DockerProcess {
    Write-Log "Stopping Docker services..."
    docker-compose -f $DOCKER_COMPOSE_FILE down
    if ($LASTEXITCODE -ne 0) {
        Write-Log "Error: Docker-compose failed to stop"
        return $false
    }
    return $true
}


# Configure git if needed
if ($GITHUB_TOKEN) {
    git config core.askPass "echo"
    git config credential.helper store
}

# Store initial commit hash
$CURRENT_COMMIT = Get-LatestCommit
Write-Log "Starting monitor for branch: $BRANCH"
Write-Log "Current commit: $CURRENT_COMMIT"

# Start initial Docker process
$dockerProcess = Start-DockerProcess
Write-Log "Started Docker process with ID: $($dockerProcess.Id)"

$CURRENT_COMMIT = Get-LatestCommit
Write-Log "Starting monitor for branch: $BRANCH"
Write-Log "Current commit: $CURRENT_COMMIT"

# Start initial Docker process
if (-not (Start-DockerProcess)) {
    Write-Log "Failed to start Docker services. Exiting..."
    exit 1
}

while ($true) {
    $LATEST_COMMIT = Get-LatestCommit
    
    if ($CURRENT_COMMIT -ne $LATEST_COMMIT) {
        Write-Log "New commit detected!"
        Write-Log "Previous commit: $CURRENT_COMMIT"
        Write-Log "New commit: $LATEST_COMMIT"
        
        Write-Log "Pulling latest changes..."
        $pullResult = git pull origin $BRANCH
        
        if ($LASTEXITCODE -eq 0) {
            # Stop current Docker process
            if (Stop-DockerProcess) {
                # Start new Docker process
                if (Start-DockerProcess) {
                    $CURRENT_COMMIT = $LATEST_COMMIT
                    Write-Log "Successfully restarted Docker services"
                } else {
                    Write-Log "Failed to restart Docker services"
                }
            } else {
                Write-Log "Failed to stop Docker services"
            }
        } else {
            Write-Log "Error: Git pull failed"
        }
    } else {
        Write-Log "No new commits detected"
    }
    
    Start-Sleep -Seconds $CHECK_INTERVAL
}
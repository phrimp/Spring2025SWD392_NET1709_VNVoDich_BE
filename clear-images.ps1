# For Linux/macOS (save as docker-cleanup.sh):
#!/bin/bash

# Function to print steps and execute commands
execute_step() {
    echo "$1"
    eval "$2"
    echo
}

echo "Starting Docker cleanup process..."
echo "======================================"

# Stop all running containers
execute_step "Stopping all running containers..." \
    "docker stop \$(docker ps -q) 2>/dev/null || echo 'No running containers to stop'"

# Remove all stopped containers
execute_step "Removing stopped containers..." \
    "docker container prune -f"

# Remove all unused images
execute_step "Removing dangling images (unused and untagged)..." \
    "docker image prune -f"

# Remove all unused images (including tagged ones)
execute_step "Removing all unused images..." \
    "docker image prune -a -f"

# Remove unused networks
execute_step "Removing unused networks..." \
    "docker network prune -f"

# Remove unused volumes
execute_step "Removing unused volumes..." \
    "docker volume prune -f"

# Display system status after cleanup
echo "System status after cleanup:"
echo "======================================"
execute_step "Remaining containers:" "docker ps -a"
execute_step "Remaining images:" "docker images"
execute_step "Disk space usage:" "docker system df"

echo "Docker cleanup completed!"

# For Windows (save as docker-cleanup.ps1):
# PowerShell Script

function Write-Step {
    param (
        [string]$Message,
        [string]$Command
    )
    Write-Host "$Message" -ForegroundColor Cyan
    Invoke-Expression $Command
    Write-Host ""
}

Write-Host "Starting Docker cleanup process..." -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Green

# Stop all running containers
Write-Step "Stopping all running containers..." {
    $running = docker ps -q
    if ($running) {
        docker stop $running
    } else {
        Write-Host "No running containers to stop"
    }
}

# Remove all stopped containers
Write-Step "Removing stopped containers..." {
    docker container prune -f
}

# Remove all unused images
Write-Step "Removing dangling images (unused and untagged)..." {
    docker image prune -f
}

# Remove all unused images (including tagged ones)
Write-Step "Removing all unused images..." {
    docker image prune -a -f
}

# Remove unused networks
Write-Step "Removing unused networks..." {
    docker network prune -f
}

# Remove unused volumes
Write-Step "Removing unused volumes..." {
    docker volume prune -f
}

# Display system status after cleanup
Write-Host "System status after cleanup:" -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Green

Write-Step "Remaining containers:" {
    docker ps -a
}

Write-Step "Remaining images:" {
    docker images
}

Write-Step "Disk space usage:" {
    docker system df
}

Write-Host "Docker cleanup completed!" -ForegroundColor Green

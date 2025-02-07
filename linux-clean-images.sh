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

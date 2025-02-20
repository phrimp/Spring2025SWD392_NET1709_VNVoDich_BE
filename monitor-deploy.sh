#!/bin/bash

# Load environment variables from .env file if it exists
if [ -f .env ]; then
  source .env
fi

# Configuration with environment variables and defaults
REPO_PATH=${GITHUB_REPO_PATH:-"https://github.com/phrimp/Spring2025SWD392_NET1709_VNVoDich_BE"}
BRANCH=${GITHUB_BRANCH:-"master"}
CHECK_INTERVAL=${CHECK_INTERVAL:-60}
GITHUB_TOKEN=${GITHUB_TOKEN:-""}
DOCKER_COMPOSE_FILE=${DOCKER_COMPOSE_FILE:-"../docker-compose.yml"}

# Function to get the latest commit hash
get_latest_commit() {
  if [ -n "$GITHUB_TOKEN" ]; then
    git ls-remote "https://${GITHUB_TOKEN}@github.com/${REPO_PATH#https://github.com/}" -h "refs/heads/$BRANCH" | cut -f1
  else
    git ls-remote origin -h "refs/heads/$BRANCH" | cut -f1
  fi
}

# Function to log messages with timestamp
log_message() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

# Initial repository setup
if [ ! -d ".git" ]; then
  log_message "Initializing repository..."
  if [ -n "$GITHUB_TOKEN" ]; then
    git clone "https://${GITHUB_TOKEN}@github.com/${REPO_PATH#https://github.com/}" .
  else
    git clone "$REPO_PATH" .
  fi
fi

# Change to repository directory or create it if it doesn't exist
mkdir -p "$(dirname "$REPO_PATH")" 2>/dev/null
cd "$(dirname "$REPO_PATH")" || {
  log_message "Error: Could not change to repository directory"
  exit 1
}

# Configure git if needed
if [ -n "$GITHUB_TOKEN" ]; then
  git config core.askPass "echo"
  git config credential.helper store
fi

# Store initial commit hash
CURRENT_COMMIT=$(get_latest_commit)
log_message "Starting monitor for branch: $BRANCH"
log_message "Current commit: $CURRENT_COMMIT"

# Main loop
while true; do
  # Fetch latest commit hash
  LATEST_COMMIT=$(get_latest_commit)

  # Check if there are new commits
  if [ "$CURRENT_COMMIT" != "$LATEST_COMMIT" ]; then
    log_message "New commit detected!"
    log_message "Previous commit: $CURRENT_COMMIT"
    log_message "New commit: $LATEST_COMMIT"

    # Pull latest changes
    log_message "Pulling latest changes..."
    if git pull origin "$BRANCH"; then
      # Rebuild and restart Docker containers
      log_message "Rebuilding and restarting Docker containers..."
      if docker-compose -f "$DOCKER_COMPOSE_FILE" up --build -d; then
        log_message "Docker containers successfully rebuilt and restarted"
        CURRENT_COMMIT=$LATEST_COMMIT
      else
        log_message "Error: Docker rebuild failed"
      fi
    else
      log_message "Error: Git pull failed"
    fi
  else
    log_message "No new commits detected"
  fi

  # Wait before next check
  sleep "$CHECK_INTERVAL"
done

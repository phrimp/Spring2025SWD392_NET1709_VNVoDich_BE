
try {
    Get-Process -Name "Docker Desktop"
} catch {
    Write-Host "Starting Docker Desktop..."
    try {
        Write-Host "Waiting for 15 seconds"
        Start-Process "C:\Program Files\Docker\Docker\Docker Desktop.exe" -Verb RunAs
        Start-Sleep -Seconds 15  # Allow time for Docker to initialize
    } catch {
        Write-Host "Docker start failed. Exit"
        pause
        exit
    }
} finally {
    docker-compose up --build
}



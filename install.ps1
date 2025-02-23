#Requires -Version 5.0

$ErrorActionPreference = 'Stop'

# Configuration
$GithubRepo = "conormkelly/yts-cli"
$BinaryName = "yts.exe"
$InstallDir = "$env:LOCALAPPDATA\Programs\yts"

function Write-ColorOutput($ForegroundColor, $Message) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    Write-Output $Message
    $host.UI.RawUI.ForegroundColor = $fc
}

function Get-Architecture {
    $arch = [System.Environment]::GetEnvironmentVariable("PROCESSOR_ARCHITECTURE")
    switch ($arch) {
        "AMD64" { return "x86_64" }
        "ARM64" { return "arm64" }
        default {
            Write-ColorOutput "Red" "Unsupported architecture: $arch"
            exit 1
        }
    }
}

function Get-LatestRelease {
    try {
        $releases = Invoke-RestMethod -Uri "https://api.github.com/repos/$GithubRepo/releases/latest"
        return @{
            Version = $releases.tag_name
            Assets = $releases.assets
        }
    } catch {
        Write-ColorOutput "Red" "Failed to fetch latest release: $_"
        exit 1
    }
}

function Test-Command($Name) {
    return [bool](Get-Command -Name $Name -ErrorAction SilentlyContinue)
}

function Install-YTS {
    Write-ColorOutput "Blue" "Installing YTS CLI..."
    
    # Check if already installed
    if (Test-Path "$InstallDir\$BinaryName") {
        Write-ColorOutput "Yellow" "YTS CLI is already installed. Updating..."
    }
    
    # Get system info
    $arch = Get-Architecture
    Write-Output "Detected architecture: $arch"
    
    # Get latest version
    $release = Get-LatestRelease
    $version = $release.Version
    Write-Output "Latest version: $version"
    
    # Create temp directory
    $tmpDir = New-TemporaryFile | ForEach-Object { 
        Remove-Item $_ -Force
        New-Item -ItemType Directory -Path $_ 
    }
    
    try {
        Push-Location $tmpDir
        
        # Download binary
        $assetName = "yts_Windows_$arch.zip"
        $asset = $release.Assets | Where-Object { $_.name -eq $assetName }
        
        if (-not $asset) {
            Write-ColorOutput "Red" "No release found for Windows $arch"
            exit 1
        }
        
        Write-Output "Downloading $assetName..."
        $progressPreference = 'silentlyContinue'
        Invoke-WebRequest -Uri $asset.browser_download_url -OutFile "yts.zip"
        $progressPreference = 'Continue'
        
        # Verify download
        if (-not (Test-Path "yts.zip")) {
            throw "Download failed"
        }
        
        # Extract binary
        Write-Output "Extracting..."
        Expand-Archive -Path "yts.zip" -DestinationPath "." -Force
        
        # Create installation directory
        if (-not (Test-Path $InstallDir)) {
            New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        }
        
        # Install binary
        Write-Output "Installing to $InstallDir..."
        Copy-Item -Force -Path $BinaryName -Destination "$InstallDir\$BinaryName"
        
        # Update PATH
        $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
        if ($userPath -notlike "*$InstallDir*") {
            Write-Output "Updating PATH..."
            [Environment]::SetEnvironmentVariable(
                "Path",
                "$userPath;$InstallDir",
                "User"
            )
        }
        
        Write-ColorOutput "Green" "YTS CLI $version has been installed successfully!"
        Write-Output "Please restart your terminal, then run 'yts --help' to get started."
        
    } catch {
        Write-ColorOutput "Red" "Installation failed: $_"
        exit 1
    } finally {
        Pop-Location
        Remove-Item -Recurse -Force $tmpDir -ErrorAction SilentlyContinue
    }
}

# Run installation
Install-YTS

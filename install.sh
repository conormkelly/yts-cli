#!/bin/bash

# YTS CLI installer for Unix-like systems
set -euo pipefail

# Configuration
GITHUB_REPO="conormkelly/yts-cli"
BINARY_NAME="yts"
INSTALL_DIR="/usr/local/bin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print error message and exit
error() {
    echo -e "${RED}Error: $1${NC}" >&2
    exit 1
}

# Print warning message
warn() {
    echo -e "${YELLOW}Warning: $1${NC}" >&2
}

# Print info message
info() {
    echo -e "${BLUE}$1${NC}"
}

# Print success message
success() {
    echo -e "${GREEN}$1${NC}"
}

# Check for required commands
check_dependencies() {
    local missing_deps=()
    
    for cmd in curl tar chmod mkdir; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            missing_deps+=("$cmd")
        fi
    done
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        error "Required commands not found: ${missing_deps[*]}"
    fi
}

# Check if we have sudo access (if needed)
check_sudo_access() {
    if [ "$EUID" -eq 0 ]; then
        return 0  # Already root
    fi
    
    if command -v sudo >/dev/null 2>&1; then
        if sudo -v >/dev/null 2>&1; then
            return 0  # Has sudo access
        fi
    fi
    
    return 1  # No sudo access
}

# Function to detect OS and architecture
detect_platform() {
    # Detect OS
    case "$(uname -s | tr '[:upper:]' '[:lower:]')" in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        *)
            error "Unsupported operating system: $(uname -s)"
            ;;
    esac
    
    # Detect architecture
    local arch
    arch=$(uname -m)
    case "${arch}" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            error "Unsupported architecture: ${arch}"
            ;;
    esac
    
    PLATFORM="${OS}_${ARCH}"
    info "Detected platform: ${PLATFORM}"
}

# Function to get latest release version
get_latest_version() {
    info "Fetching latest release information..."
    
    local latest_release
    if ! latest_release=$(curl -sL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest"); then
        error "Failed to fetch latest release information"
    fi
    
    VERSION=$(echo "${latest_release}" | grep -o '"tag_name": "[^"]*' | cut -d'"' -f4)
    if [ -z "${VERSION}" ]; then
        error "Failed to determine latest version"
    fi
    
    info "Latest version: ${VERSION}"
}

# Function to create temporary directory with cleanup
setup_temp_dir() {
    TMP_DIR=$(mktemp -d)
    if [ ! -d "${TMP_DIR}" ]; then
        error "Failed to create temporary directory"
    fi
    
    # Set up cleanup trap
    trap 'rm -rf "${TMP_DIR}"' EXIT
}

# Function to download and verify release
download_release() {
    local download_url="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/yts_${OS}_${ARCH}.tar.gz"
    info "Downloading from: ${download_url}"
    
    if ! curl -sL --retry 3 "${download_url}" -o "${TMP_DIR}/yts.tar.gz"; then
        error "Failed to download release"
    fi
    
    # Extract the archive
    if ! tar xzf "${TMP_DIR}/yts.tar.gz" -C "${TMP_DIR}"; then
        error "Failed to extract archive"
    fi
}

# Function to install the binary
install_binary() {
    info "Installing to ${INSTALL_DIR}..."
    
    # Ensure install directory exists
    if [ ! -d "${INSTALL_DIR}" ]; then
        if ! mkdir -p "${INSTALL_DIR}"; then
            if ! sudo mkdir -p "${INSTALL_DIR}"; then
                error "Failed to create installation directory"
            fi
        fi
    fi
    
    # Copy binary and set permissions
    if [ -w "${INSTALL_DIR}" ]; then
        cp "${TMP_DIR}/yts" "${INSTALL_DIR}/"
        chmod +x "${INSTALL_DIR}/yts"
    else
        if ! check_sudo_access; then
            error "Installation requires sudo access to ${INSTALL_DIR}"
        fi
        sudo cp "${TMP_DIR}/yts" "${INSTALL_DIR}/"
        sudo chmod +x "${INSTALL_DIR}/yts"
    fi
}

# Function to verify installation
verify_installation() {
    if ! command -v yts >/dev/null 2>&1; then
        warn "Installation completed, but 'yts' is not in PATH"
        echo "Add ${INSTALL_DIR} to your PATH or restart your terminal"
    else
        success "YTS CLI has been installed successfully!"
        echo "Run 'yts --help' to get started"
    fi
}

# Main installation flow
main() {
    info "Installing YTS CLI..."
    
    check_dependencies
    detect_platform
    get_latest_version
    setup_temp_dir
    download_release
    install_binary
    verify_installation
}

# Run main installation
main

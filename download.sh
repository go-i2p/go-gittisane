#!/bin/bash

# Set GitHub repo info
OWNER="go-i2p"
REPO="go-gittisane"

echo "Fetching latest release info from GitHub..."

# Get latest release data
RELEASE_DATA=$(curl -s "https://api.github.com/repos/$OWNER/$REPO/releases/latest")

# Extract version number
VERSION=$(echo "$RELEASE_DATA" | grep -Po '"tag_name": "\K.*?(?=")')
echo "Latest version: $VERSION"

# Create downloads directory
mkdir -p downloads
cd downloads

# Download each asset
echo "$RELEASE_DATA" | grep -Po '"browser_download_url": "\K.*?(?=")' | while read -r url; do
    filename=$(basename "$url")
    echo "Downloading $filename..."
    curl -L -o "$filename" "$url"
    
    # Make Linux/macOS binaries executable
    if [[ "$filename" != *".exe" ]]; then
        chmod +x "$filename"
    fi
done

echo "Download complete! Files are in the 'downloads' directory"
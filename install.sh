#!/bin/bash

# Exit on any error
set -e

# Build the Go application
echo "Building mac_agent..."
go build -o mac_agent main.go

# Create the installation directory
INSTALL_DIR="$HOME/Library/mac_agent"
mkdir -p "$INSTALL_DIR"

# Move the built application to the installation directory
mv mac_agent "$INSTALL_DIR/"

# Create the LaunchAgent plist file
PLIST_FILE="$HOME/Library/LaunchAgents/com.example.mac_agent.plist"
cat > "$PLIST_FILE" <<EOL
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.example.mac_agent</string>
    <key>ProgramArguments</key>
    <array>
        <string>$INSTALL_DIR/mac_agent</string>
    </array>
    <key>StartInterval</key>
    <integer>300</integer>
    <key>RunAtLoad</key>
    <true/>
    <key>StandardOutPath</key>
    <string>$HOME/Library/Logs/mac_agent/stdout.log</string>
    <key>StandardErrorPath</key>
    <string>$HOME/Library/Logs/mac_agent/stderr.log</string>
</dict>
</plist>
EOL

# Load the LaunchAgent
launchctl load "$PLIST_FILE"

echo "mac_agent has been installed and scheduled to run every 5 minutes."
echo "Logs will be saved in $HOME/Library/Logs/mac_agent/"
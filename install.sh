#!/bin/bash

# Exit on any error
set -e

# Build the Go application
echo "Building system-agent..."
go build -o system-agent main.go

# Create the installation directory
INSTALL_DIR="$HOME/Library/system-agent"
mkdir -p "$INSTALL_DIR"

# Move the built application to the installation directory
mv system-agent "$INSTALL_DIR/"

# Create the LaunchAgent plist file
PLIST_FILE="$HOME/Library/LaunchAgents/com.example.system-agent.plist"
cat > "$PLIST_FILE" <<EOL
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.example.system-agent</string>
    <key>ProgramArguments</key>
    <array>
        <string>$INSTALL_DIR/system-agent</string>
    </array>
    <key>StartInterval</key>
    <integer>300</integer>
    <key>RunAtLoad</key>
    <true/>
    <key>StandardOutPath</key>
    <string>$HOME/Library/Logs/system-agent/stdout.log</string>
    <key>StandardErrorPath</key>
    <string>$HOME/Library/Logs/system-agent/stderr.log</string>
</dict>
</plist>
EOL

# Load the LaunchAgent
launchctl load "$PLIST_FILE"

echo "system-agent has been installed and scheduled to run every 5 minutes."
echo "Logs will be saved in $HOME/Library/Logs/system-agent/"
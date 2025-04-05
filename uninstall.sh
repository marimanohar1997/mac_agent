#!/bin/bash

# Exit on any error
set -e

AGENT_LABEL="com.example.mac_agent"
PLIST_FILE="$HOME/Library/LaunchAgents/$AGENT_LABEL.plist"
INSTALL_DIR="$HOME/Library/mac_agent"
LOG_DIR="$HOME/Library/Logs/mac_agent"

echo "Uninstalling mac_agent..."

# Unload the Launch Agent if it's loaded
if launchctl list | grep -q "$AGENT_LABEL"; then
    launchctl unload "$PLIST_FILE"
    echo "Launch agent unloaded."
else
    echo "Launch agent not loaded or already removed."
fi

# Remove the plist
if [ -f "$PLIST_FILE" ]; then
    rm "$PLIST_FILE"
    echo "Removed plist file."
else
    echo "Plist file not found."
fi

# Remove the installation directory
if [ -d "$INSTALL_DIR" ]; then
    rm -rf "$INSTALL_DIR"
    echo "Removed installed mac_agent files."
else
    echo "Installation directory not found."
fi

# Remove logs
if [ -d "$LOG_DIR" ]; then
    rm -rf "$LOG_DIR"
    echo "Removed log directory."
else
    echo "Log directory not found."
fi

echo "mac_agent uninstalled successfully!"
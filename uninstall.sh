#!/bin/bash

# Unload the Launch Agent
launchctl remove com.example.system-agent

# Remove the binary and plist file
sudo rm /usr/local/bin/system-agent
rm ~/Library/LaunchAgents/com.example.system-agent.plist

# Optional: Remove log directory
sudo rm -rf /usr/local/var/log

echo "Agent uninstalled!"
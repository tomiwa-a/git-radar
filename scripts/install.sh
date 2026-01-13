#!/bin/bash
set -e

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="git-radar"

if [ -f "./$BINARY_NAME" ]; then
    BINARY_PATH="./$BINARY_NAME"
elif [ -f "./dist/$BINARY_NAME-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)" ]; then
    BINARY_PATH="./dist/$BINARY_NAME-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)"
else
    echo "Error: Binary not found. Run 'make build' first."
    exit 1
fi

echo "Installing $BINARY_NAME to $INSTALL_DIR..."
sudo cp "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo "✓ git-radar installed successfully!"
echo "  Run 'git-radar' from any git repository."

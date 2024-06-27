#!/usr/bin/env bash

# Description: This script is used to run the applications and controllers.
# It is used to start the applications and controllers, and to clean up
# when the script exits.

# Check the number of arguments
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <number of apps>"
    exit 1
fi

# Check if the argument is a positive integer
if ! [[ "$1" =~ ^[0-9]+$ ]]; then
    echo "The argument must be a positive integer"
    exit 1
fi

# Check if the number of apps is at least 1
if [ "$1" -lt 1 ]; then
    echo "The number of apps must be at least 1"
    exit 1
fi

NUM_APPS="$1"

# Start the application
./scripts/v2/startup.sh "$NUM_APPS"

# Start the applications and controllers
for ((i = 0; i < NUM_APPS; i++)); do
    ./scripts/v2/app_terminal.sh "$i" &
    ./scripts/v2/ctl_terminal.sh "$i" &
done

echo "Press Ctrl+C to stop the application..."

# Clean-up when the script exits
trap cleanup EXIT
function cleanup() {
    ./scripts/v2/cleanup.sh "$NUM_APPS"
    local APP_TERMINAL_TITLE="Application"
    local CTL_TERMINAL_TITLE="Controller"
    for ((i = 0; i < NUM_APPS; i++)); do
        pkill -f "xterm -T $APP_TERMINAL_TITLE $i" 2>/dev/null
        pkill -f "xterm -T $CTL_TERMINAL_TITLE $i" 2>/dev/null
    done
    exit 0
}

while true; do
    sleep 1
done

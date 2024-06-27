# Description: This script is used to start the application and open
# terminals for each node.

# Start the application
./scripts/v1/startup.sh
./scripts/v1/app_terminal.sh 0
./scripts/v1/app_terminal.sh 1
./scripts/v1/app_terminal.sh 2

# Clean-up when the script exits
trap './scripts/v1/cleanup.sh' EXIT

# Wait for the user to press Enter
read -p "Press Enter to stop the application..."

# Stop the application
TERMINAL_TITLE="Application"
pkill -f "xterm -T $TERMINAL_TITLE 0" 2>/dev/null
pkill -f "xterm -T $TERMINAL_TITLE 1" 2>/dev/null
pkill -f "xterm -T $TERMINAL_TITLE 2" 2>/dev/null

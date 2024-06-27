# Start the application
./scripts/v3/startup.sh

for ((i = 0; i < 5; i++)); do
    ./scripts/v2/app_terminal.sh "$i" &
    ./scripts/v2/ctl_terminal.sh "$i" &
done
./scripts/v3/net_terminal.sh 0
./scripts/v3/net_terminal.sh 1
./scripts/v3/net_terminal.sh 2
./scripts/v3/net_terminal.sh 3
./scripts/v3/net_terminal.sh 4


# Clean-up when the script exits
trap './scripts/v3/cleanup.sh' EXIT

# Wait for the user to press Enter
read -p "Press Enter to stop the application..."

# Stop the application
TERMINAL_TITLE="Net"
pkill -f "xterm -T $TERMINAL_TITLE 0" 2>/dev/null
pkill -f "xterm -T $TERMINAL_TITLE 1" 2>/dev/null
pkill -f "xterm -T $TERMINAL_TITLE 2" 2>/dev/null
pkill -f "xterm -T $TERMINAL_TITLE 3" 2>/dev/null
pkill -f "xterm -T $TERMINAL_TITLE 4" 2>/dev/null
pkill -f "xterm -T $TERMINAL_TITLE 5" 2>/dev/null
pkill -f "xterm -T $TERMINAL_TITLE 6" 2>/dev/null

APP_TERMINAL_TITLE="Application"
CTL_TERMINAL_TITLE="Controller"
for ((i = 0; i < 5; i++)); do
    pkill -f "xterm -T $APP_TERMINAL_TITLE $i" 2>/dev/null
    pkill -f "xterm -T $CTL_TERMINAL_TITLE $i" 2>/dev/null
done

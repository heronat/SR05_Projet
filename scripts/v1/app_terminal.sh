#!/usr/bin/env bash

# Check the number of arguments
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <app_id>"
    exit 1
fi

APP_ID="$1"
TMP_DIR=/tmp
TERMINAL_TITLE="Application $APP_ID"

if [ ! -d $TMP_DIR ]; then
    echo "The temporary directory does not exist"
    echo "Exiting..."
    exit 1
fi

IN_PIPE="$TMP_DIR/in_A$APP_ID"
ERROR_LOG="$TMP_DIR/error_A$APP_ID"

if [ ! -p "$IN_PIPE" ]; then
    echo "The named pipe does not exist"
    echo "Exiting..."
    exit 1
fi

if [ ! -f "$ERROR_LOG" ]; then
    echo "The error log does not exist"
    echo "Exiting..."
    exit 1
fi

# Start the terminal
xterm -T "$TERMINAL_TITLE" -e "bash -c '
    echo \"Application $APP_ID\"
    tail -f $ERROR_LOG &
    while true; do
        read -p \"Enter a message: \" INPUT
        echo \$INPUT > $IN_PIPE
    done
'" &

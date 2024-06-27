#!/usr/bin/env bash

# Check the number of arguments
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <net_id>"
    exit 1
fi

if ! [[ "$1" =~ ^[0-9]+$ ]]; then
    echo "The argument must be a positive integer"
    exit 1
fi

NET_ID="$1"
TMP_DIR=/tmp
TERMINAL_TITLE="Net $NET_ID"

if [ ! -d $TMP_DIR ]; then
    echo "The temporary directory does not exist"
    echo "Exiting..."
    exit 1
fi

IN_PIPE="$TMP_DIR/in_N$NET_ID"
ERROR_LOG="$TMP_DIR/error_N$NET_ID"

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
    echo -e \"\033[1mNet $NET_ID\033[0m\"
    tail -f $ERROR_LOG &
    while true; do
        read INPUT
        echo \$INPUT > $IN_PIPE
    done
'" &
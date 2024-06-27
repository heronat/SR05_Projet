#!/usr/bin/env bash

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

# Define the path to the temporary directory
TMP_DIR=/tmp

if [ ! -d $TMP_DIR ]; then
    echo "The temporary directory does not exist"
    echo "Exiting..."
    exit 1
fi

APP_CMD=app
CTL_CMD=ctl

# Stop the processes
pkill -f $APP_CMD 2>/dev/null
pkill -f $CTL_CMD 2>/dev/null

for ((i = 0; i < NUM_APPS; i++)); do
    # Define the path to the named pipes
    IN_A="$TMP_DIR/in_A$i"
    IN_C="$TMP_DIR/in_C$i"
    OUT_A="$TMP_DIR/out_A$i"
    OUT_C="$TMP_DIR/out_C$i"
    ERROR_A="$TMP_DIR/error_A$i"
    ERROR_C="$TMP_DIR/error_C$i"

    pkill -f "cat $OUT_A >$IN_C" 2>/dev/null
    pkill -f "tee $OUT_A" 2>/dev/null

    # Remove the named pipes
    rm -f "$IN_A" "$IN_C" "$OUT_A" "$OUT_C" "$ERROR_A" "$ERROR_C"
done

echo "Clean-up complete."

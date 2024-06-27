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

export PATH=$PATH:/usr/local/go/bin

# Define the path to the temporary directory
TMP_DIR=/tmp

if [ ! -d "$TMP_DIR" ]; then
    echo "The temporary directory does not exist"
    echo "Do you want to create it? [y/n]"
    read -r answer
    if [ "$answer" = "y" ]; then
        mkdir "$TMP_DIR"
    else
        echo "Exiting..."
        exit 1
    fi
fi

# Define the paths to the Go programs
APP_GO=app.go
CTL_GO=ctl.go

if [ ! -f "$APP_GO" ]; then
    echo "The Go program $APP_GO does not exist"
    exit 1
fi

if [ ! -f "$CTL_GO" ]; then
    echo "The Go program $CTL_GO does not exist"
    exit 1
fi

# Build the Go programs
go build -o "$TMP_DIR/app" $APP_GO
go build -o "$TMP_DIR/ctl" $CTL_GO
APP="$TMP_DIR/app -id"
CTL="$TMP_DIR/ctl -id"

for ((i = 0; i < NUM_APPS; i++)); do
    # Define the paths to the named pipes
    IN_A="$TMP_DIR/in_A$i"
    IN_C="$TMP_DIR/in_C$i"
    OUT_A="$TMP_DIR/out_A$i"
    OUT_C="$TMP_DIR/out_C$i"
    ERROR_A="$TMP_DIR/error_A$i"
    ERROR_C="$TMP_DIR/error_C$i"

    # Create the named pipes if they do not exist
    rm -f "$IN_A" "$IN_C" "$OUT_A" "$OUT_C"
    mkfifo "$IN_A" "$IN_C" "$OUT_A" "$OUT_C"

    # Start the processes in the background
    $APP "$i" <"$IN_A" >"$OUT_A" 2>"$ERROR_A" &
    $CTL "$i" <"$IN_C" >"$OUT_C" 2>"$ERROR_C" &
done

for ((i = 0; i < NUM_APPS; i++)); do
    NEXT=$((i + 1))
    if [ "$i" -eq "$((NUM_APPS - 1))" ]; then
        NEXT=0
    fi

    IN_A="$TMP_DIR/in_A$i"
    IN_C="$TMP_DIR/in_C$i"
    OUT_A="$TMP_DIR/out_A$i"
    OUT_C="$TMP_DIR/out_C$i"

    IN_A_NEXT="$TMP_DIR/in_A$NEXT"
    IN_C_NEXT="$TMP_DIR/in_C$NEXT"

    cat "$OUT_A" >"$IN_C" &
    cat "$OUT_C" | tee "$IN_A" >"$IN_C_NEXT" &
done

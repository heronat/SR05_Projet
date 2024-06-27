#!/usr/bin/env bash

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

# Define the paths to the named pipes
IN_A0="$TMP_DIR/in_A0"
IN_C0="$TMP_DIR/in_C0"
OUT_A0="$TMP_DIR/out_A0"
OUT_C0="$TMP_DIR/out_C0"
ERROR_A0="$TMP_DIR/error_A0"
ERROR_C0="$TMP_DIR/error_C0"

IN_A1="$TMP_DIR/in_A1"
IN_C1="$TMP_DIR/in_C1"
OUT_A1="$TMP_DIR/out_A1"
OUT_C1="$TMP_DIR/out_C1"
ERROR_A1="$TMP_DIR/error_A1"
ERROR_C1="$TMP_DIR/error_C1"

IN_A2="$TMP_DIR/in_A2"
IN_C2="$TMP_DIR/in_C2"
OUT_A2="$TMP_DIR/out_A2"
OUT_C2="$TMP_DIR/out_C2"
ERROR_A2="$TMP_DIR/error_A2"
ERROR_C2="$TMP_DIR/error_C2"

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

# Create the named pipes if they do not exist
rm -f "$IN_A0" "$IN_C0" "$OUT_A0" "$OUT_C0"
mkfifo "$IN_A0" "$IN_C0" "$OUT_A0" "$OUT_C0"

rm -f "$IN_A1" "$IN_C1" "$OUT_A1" "$OUT_C1"
mkfifo "$IN_A1" "$IN_C1" "$OUT_A1" "$OUT_C1"

rm -f "$IN_A2" "$IN_C2" "$OUT_A2" "$OUT_C2"
mkfifo "$IN_A2" "$IN_C2" "$OUT_A2" "$OUT_C2"

# Start the processes in the background
$APP 0 <"$IN_A0" >"$OUT_A0" 2>"$ERROR_A0" &
$CTL 0 <"$IN_C0" >"$OUT_C0" 2>"$ERROR_C0" &

$APP 1 <"$IN_A1" >"$OUT_A1" 2>"$ERROR_A1" &
$CTL 1 <"$IN_C1" >"$OUT_C1" 2>"$ERROR_C1" &

$APP 2 <"$IN_A2" >"$OUT_A2" 2>"$ERROR_A2" &
$CTL 2 <"$IN_C2" >"$OUT_C2" 2>"$ERROR_C2" &

# Connect the named pipes
cat "$OUT_A0" >"$IN_C0" &
cat "$OUT_C0" | tee "$IN_A0" >"$IN_C1" &
cat "$OUT_A1" >"$IN_C1" &
cat "$OUT_C1" | tee "$IN_A1" >"$IN_C2" &
cat "$OUT_A2" >"$IN_C2" &
cat "$OUT_C2" | tee "$IN_A2" >"$IN_C0" &

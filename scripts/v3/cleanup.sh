#!/usr/bin/env bash

# Define the path to the temporary directory
TMP_DIR=/tmp

if [ ! -d $TMP_DIR ]; then
    echo "The temporary directory does not exist"
    echo "Exiting..."
    exit 1
fi

NET_CMD=net
APP_CMD=app
CTL_CMD=ctl

# Stop the processes
pkill -f $NET_CMD 2>/dev/null
pkill -f $APP_CMD 2>/dev/null
pkill -f $CTL_CMD 2>/dev/null

IN_N="$TMP_DIR/in_N*"
OUT_N="$TMP_DIR/out_N*"
ERROR_N="$TMP_DIR/error_N*"
MAP_N="$TMP_DIR/map_N*"
IN_A="$TMP_DIR/in_A*"
OUT_A="$TMP_DIR/out_A*"
ERROR_A="$TMP_DIR/error_A*"
IN_C="$TMP_DIR/in_C*"
OUT_C="$TMP_DIR/out_C*"
ERROR_C="$TMP_DIR/error_C*"
pkill -f /tmp 2>/dev/null


rm -f $IN_N $OUT_N $ERROR_N $MAP_N $IN_A $OUT_A $ERROR_A $IN_C $OUT_C $ERROR_C

echo "Clean-up complete."

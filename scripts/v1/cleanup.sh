#!/usr/bin/env bash

# Define the path to the temporary directory
TMP_DIR=/tmp

if [ ! -d $TMP_DIR ]; then
    echo "The temporary directory does not exist"
    echo "Exiting..."
    exit 1
fi

# Define the path to the named pipes
IN_A0=$TMP_DIR/in_A0
IN_C0=$TMP_DIR/in_C0
OUT_A0=$TMP_DIR/out_A0
OUT_C0=$TMP_DIR/out_C0
ERROR_A0="$TMP_DIR/error_A0"
ERROR_C0="$TMP_DIR/error_C0"

IN_A1=$TMP_DIR/in_A1
IN_C1=$TMP_DIR/in_C1
OUT_A1=$TMP_DIR/out_A1
OUT_C1=$TMP_DIR/out_C1
ERROR_A1="$TMP_DIR/error_A1"
ERROR_C1="$TMP_DIR/error_C1"

IN_A2=$TMP_DIR/in_A2
IN_C2=$TMP_DIR/in_C2
OUT_A2=$TMP_DIR/out_A2
OUT_C2=$TMP_DIR/out_C2
ERROR_A2="$TMP_DIR/error_A2"
ERROR_C2="$TMP_DIR/error_C2"

APP_CMD=app
CTL_CMD=ctl

# Stop the processes
pkill -f $APP_CMD 2>/dev/null
pkill -f $CTL_CMD 2>/dev/null

pkill -f "cat $OUT_A0 >$IN_C0" 2>/dev/null
pkill -f "cat $OUT_A1 >$IN_C1" 2>/dev/null
pkill -f "cat $OUT_A2 >$IN_C2" 2>/dev/null

pkill -f "tee $OUT_A0" 2>/dev/null
pkill -f "tee $OUT_A1" 2>/dev/null
pkill -f "tee $OUT_A2" 2>/dev/null

# Remove the named pipes
rm -f $TMP_DIR/in* $TMP_DIR/out* $TMP_DIR/error*

echo "Clean-up complete."

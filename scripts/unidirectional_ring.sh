#!/usr/bin/env bash

export PATH=$PATH:/usr/local/go/bin

# Define the path to the temporary directory
TMP_DIR=/tmp

if [ ! -d $TMP_DIR ]; then
    echo "The temporary directory does not exist"
    echo "Do you want to create it? [y/n]"
    read -r answer
    if [ "$answer" = "y" ]; then
        mkdir $TMP_DIR
    else
        echo "Exiting..."
        exit 1
    fi
    exit 1
fi

# Define the path to the named pipes
IN_A=$TMP_DIR/in_A
IN_B=$TMP_DIR/in_B
IN_C=$TMP_DIR/in_C
OUT_A=$TMP_DIR/out_A
OUT_B=$TMP_DIR/out_B
OUT_C=$TMP_DIR/out_C

CTL_GO=ctl.go

if [ ! -f "$CTL_GO" ]; then
    echo "The Go program $CTL_GO does not exist"
    exit 1
fi

# Define the path to the ctl executable
go build -o $TMP_DIR/ctl $CTL_GO
CTL="$TMP_DIR/ctl -id"

# Create the named pipes if they do not exist
rm -f $IN_A $IN_B $IN_C $OUT_A $OUT_B $OUT_C
mkfifo $IN_A $IN_B $IN_C $OUT_A $OUT_B $OUT_C

# Start the processes in the background
$CTL 0 <$IN_A >$OUT_A &
$CTL 1 <$IN_B >$OUT_B &
$CTL 2 <$IN_C >$OUT_C &

# Connect the named pipes
cat $OUT_A >$IN_B &
cat $OUT_B >$IN_C &
cat $OUT_C >$IN_A &

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
PIPE=$TMP_DIR/bidirectional_pipe

CTL_GO=ctl.go

if [ ! -f "$CTL_GO" ]; then
    echo "The Go program $CTL_GO does not exist"
    exit 1
fi

# Define the path to the ctl executable
go build -o $TMP_DIR/ctl $CTL_GO
CTL="$TMP_DIR/ctl -id"

# Create the named pipes if they do not exist
rm -f $PIPE
mkfifo $PIPE

# Start the processes in the background
$CTL 0 <$PIPE | $CTL 1 >$PIPE

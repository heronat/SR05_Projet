#!/usr/bin/env bash

export PATH=$PATH:/usr/local/go/bin

# Define the path to the temporary directory
TMP_DIR="/tmp"

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
IN_N0="$TMP_DIR/in_N0"
IN_N1="$TMP_DIR/in_N1"
IN_N2="$TMP_DIR/in_N2"
IN_N3="$TMP_DIR/in_N3"
IN_N4="$TMP_DIR/in_N4"

OUT_N0="$TMP_DIR/out_N0"
OUT_N1="$TMP_DIR/out_N1"
OUT_N2="$TMP_DIR/out_N2"
OUT_N3="$TMP_DIR/out_N3"
OUT_N4="$TMP_DIR/out_N4"

ERROR_N0="$TMP_DIR/error_N0"
ERROR_N1="$TMP_DIR/error_N1"
ERROR_N2="$TMP_DIR/error_N2"
ERROR_N3="$TMP_DIR/error_N3"
ERROR_N4="$TMP_DIR/error_N4"

MAP_N0="$TMP_DIR/map_N0"
MAP_N1="$TMP_DIR/map_N1"
MAP_N2="$TMP_DIR/map_N2"
MAP_N3="$TMP_DIR/map_N3"
MAP_N4="$TMP_DIR/map_N4"



# Define the paths to the Go programs
APP_GO=app.go
CTL_GO=ctl.go
NET_GO=net.go

if [ ! -f "$NET_GO" ]; then
    echo "The Go program $NET_GO does not exist"
    exit 1
fi

# Build the Go programs
go build -o "$TMP_DIR/net" $NET_GO
go build -o "$TMP_DIR/app" $APP_GO
go build -o "$TMP_DIR/ctl" $CTL_GO


NET="$TMP_DIR/net -id"
APP="$TMP_DIR/app -id"
CTL="$TMP_DIR/ctl -id"



for ((i = 0; i < 5; i++)); do
    # Define the paths to the named pipes
    IN_A="$TMP_DIR/in_A$i"
    IN_C="$TMP_DIR/in_C$i"

    OUT_A="$TMP_DIR/out_A$i"
    OUT_C="$TMP_DIR/out_C$i"

    ERROR_A="$TMP_DIR/error_A$i"
    ERROR_C="$TMP_DIR/error_C$i"

    IN_NET="$TMP_DIR/in_N$i"
    OUT_NET="$TMP_DIR/out_N$i"
    ERROR_NET="$TMP_DIR/error_N$i"

    # Create the named pipes if they do not exist
    rm -f "$IN_A" "$IN_C" "$OUT_A" "$OUT_C" "$IN_NET" "$OUT_NET"
    mkfifo "$IN_A" "$IN_C" "$OUT_A" "$OUT_C" "$IN_NET" "$OUT_NET"

    rm -f "$ERROR_A" "$ERROR_C" "$ERROR_NET"
    touch "$ERROR_A" "$ERROR_C" "$ERROR_NET"

    # Start the processes in the background
    $APP "$i" <"$IN_A" >"$OUT_A" 2>"$ERROR_A" &
    $CTL "$i" <"$IN_C" >"$OUT_C" 2>"$ERROR_C" &
    $NET "$i" <"$IN_NET" >"$OUT_NET" 2>"$ERROR_NET" &
done

OUT_A0="$TMP_DIR/out_A0"
OUT_A1="$TMP_DIR/out_A1"
OUT_A2="$TMP_DIR/out_A2"
OUT_A3="$TMP_DIR/out_A3"
OUT_A4="$TMP_DIR/out_A4"

OUT_C0="$TMP_DIR/out_C0"
OUT_C1="$TMP_DIR/out_C1"
OUT_C2="$TMP_DIR/out_C2"
OUT_C3="$TMP_DIR/out_C3"
OUT_C4="$TMP_DIR/out_C4"

IN_A0="$TMP_DIR/in_A0"
IN_A1="$TMP_DIR/in_A1"
IN_A2="$TMP_DIR/in_A2"
IN_A3="$TMP_DIR/in_A3"
IN_A4="$TMP_DIR/in_A4"

IN_C0="$TMP_DIR/in_C0"
IN_C1="$TMP_DIR/in_C1"
IN_C2="$TMP_DIR/in_C2"
IN_C3="$TMP_DIR/in_C3"
IN_C4="$TMP_DIR/in_C4"

for ((i = 0; i < 5; i++)); do
    IN_A="$TMP_DIR/in_A$i"
    IN_C="$TMP_DIR/in_C$i"
    OUT_A="$TMP_DIR/out_A$i"
    OUT_C="$TMP_DIR/out_C$i"

    IN_A_NEXT="$TMP_DIR/in_A$NEXT"
    IN_C_NEXT="$TMP_DIR/in_C$NEXT"
    IN_NET="$TMP_DIR/in_N$i"
    cat "$OUT_A" >"$IN_C" &
    cat "$OUT_C"| tee "$IN_NET" >"$IN_A" &
done

# Connect the named pipes to the Go programs
cat "$OUT_N0" | tee "$IN_N1" "$IN_N2" "$IN_N3" >"$IN_C0" &
echo "$IN_C0" >> $MAP_N0
echo "$IN_N1" >> $MAP_N0
echo "$IN_N2" >> $MAP_N0
echo "$IN_N3" >> $MAP_N0

cat "$OUT_N1" | tee "$IN_N0" >"$IN_C1" &
echo "$IN_C1" >> $MAP_N1
echo "$IN_N0" >> $MAP_N1

cat "$OUT_N2" | tee "$IN_N0" >"$IN_C2" &
echo "$IN_C2" >> $MAP_N2
echo "$IN_N0" >> $MAP_N2

cat "$OUT_N3" | tee "$IN_N0" "$IN_N4">"$IN_C3" &
echo "$IN_C3" >> $MAP_N3
echo "$IN_N0" >> $MAP_N3
echo "$IN_N4" >> $MAP_N3

cat "$OUT_N4" | tee "$IN_N3" >"$IN_C4" &
echo "$IN_C4" >> $MAP_N4
echo "$IN_N3" >> $MAP_N4

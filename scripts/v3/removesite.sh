#!/usr/bin/env bash 

export PATH=$PATH:/usr/local/go/bin

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <id du site a supprimer>"
    exit 1
fi

if ! [[ "$1" =~ ^[0-9]+$ ]]; then
    echo "The argument must be a positive integer"
    exit 1
fi
MAP_TODELETE="/tmp/map_N$1"
OUT_TODELETE="/tmp/out_N$1"
ERROR_TODELETE="/tmp/error_N$1"
if [ ! -f $MAP_TODELETE ]; then
    echo "The site with the id $1 does not exist"
    exit 1
fi



mapfile -t listeconn < $MAP_TODELETE
nombre_lignes=${#listeconn[@]}
echo "nombre de lignes map a supprimer: $nombre_lignes"

if [ "$nombre_lignes" -eq 2 ]; then
    CAT_TODELETE=$(pgrep -f "cat $OUT_TODELETE")
    TAIL_TODELETE=$(pgrep -f "tail -f $ERROR_TODELETE")
    NET_TODELETE=$(pgrep -f "/tmp/net -id $1")
    CTL_TODELETE=$(pgrep -f "/tmp/ctl -id $1")
    APP_TODELETE=$(pgrep -f "/tmp/app -id $1")
    XTERM_TODELETE=$(pgrep -f "xterm -T Net $1")
    XTERM_CTL_TODELETE=$(pgrep -f "xterm -T Controller $1")
    XTERM_APP_TODELETE=$(pgrep -f "xterm -T Application $1")
    echo "Signal sent"
    sleep 5
    kill $CAT_TODELETE
    echo "kill de $CAT_TODELETE, $OUT_TODELETE"
    kill $TAIL_TODELETE
    echo "kill de $TAIL_TODELETE, $ERROR_TODELETE"
    kill $NET_TODELETE
    echo "kill de $NET_TODELETE"
    kill $CTL_TODELETE
    echo "kill de $CTL_TODELETE"
    kill $APP_TODELETE
    echo "kill de $APP_TODELETE"
    kill $XTERM_TODELETE
    echo "kill de $XTERM_TODELETE"
    kill $XTERM_CTL_TODELETE
    echo "kill de $XTERM_CTL_TODELETE"
    kill $XTERM_APP_TODELETE
    echo "kill de $XTERM_APP_TODELETE"

   
    END_PIPE=$(grep -oP '/tmp/in_N\K\d+' $MAP_TODELETE) #Search where deleted site was connected
    echo "END_PIPE: $END_PIPE"
    sed -i "/\/tmp\/in_N$1/d" "/tmp/map_N$END_PIPE" #Delete deleted site from map file
    rm -f $MAP_TODELETE $OUT_TODELETE $ERROR_TODELETE
    OUT_OLD="/tmp/out_N$END_PIPE"
     #Create new output for the site who stays connected
    PARENT_PROCS=$(pgrep -f "cat $OUT_OLD")

    for PROC in $PARENT_PROCS; do
        echo "site out: $PROC"
        kill $PROC
    done

    MAP_TOCHECK="/tmp/map_N$END_PIPE"
    mapfile -t listeout < $MAP_TOCHECK
    nombre_lignes=${#listeout[@]}
    echo "nombre de lignes map bout du pipe: $nombre_lignes"
    if [ "$nombre_lignes" -eq 1 ]; then
        cat "$OUT_OLD" > "${listeout[0]}" &
    else
        # Créer une liste des fichiers pour tee, sauf le dernier fichier
        tee_files=("${listeout[@]:0:$((nombre_lignes-1))}")
        # Le dernier fichier pour la redirection >
        output="${listeout[$((nombre_lignes-1))]}"
        # Construire la commande avec tee et >
        cat "$OUT_OLD" | tee "${tee_files[@]}" > "$output" &
        echo "Utilisation de la commande 'tee' vers les fichiers : ${tee_files[*]}"
        echo "et la redirection '>' vers $output"
    fi
    

else
    echo "The site with the id $1 is connected to mutliple sites"
    echo " Kill interface"
    XTERM_CTL_TODELETE=$(pgrep -f "xterm -T Controller $1")
    XTERM_APP_TODELETE=$(pgrep -f "xterm -T Application $1")
    kill $XTERM_CTL_TODELETE
    echo "kill de $XTERM_CTL_TODELETE"
    kill $XTERM_APP_TODELETE
    echo "kill de $XTERM_APP_TODELETE"
fi

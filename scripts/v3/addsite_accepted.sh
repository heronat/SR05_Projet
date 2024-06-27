#!/usr/bin/env bash

export PATH=$PATH:/usr/local/go/bin

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <debut du pipe(site original)> <fin du pipe(site nouveau)"
    exit 1
fi

if ! [[ "$1" =~ ^[0-9]+$ ]]; then
    echo "The argument must be a positive integer"
    exit 1
fi

if ! [[ "$2" =~ ^[0-9]+$ ]]; then
    echo "The argument must be a positive integer"
    exit 1
fi

DEBUT_PIPE="$1"
FIN_PIPE="$2"
TMP_DIR=/tmp
IN_NEW="$TMP_DIR/in_N$FIN_PIPE"
OUT_NEW="$TMP_DIR/out_N$FIN_PIPE"
ERROR_NEW="$TMP_DIR/error_N$FIN_PIPE"
IN_OLD="$TMP_DIR/in_N$DEBUT_PIPE"
OUT_OLD="$TMP_DIR/out_N$DEBUT_PIPE"
MAP_NEW="$TMP_DIR/map_N$FIN_PIPE"
MAP_OLD="$TMP_DIR/map_N$DEBUT_PIPE"
IN_CNEW="$TMP_DIR/in_C$FIN_PIPE"

echo "Chemin vers le nouveau site : $IN_NEW"
if [ -f $MAP_NEW ]; then
    echo "Le site avec l'id $2 existe déjà"
    exit 1
fi

NET_GO=net.go

if [ ! -f "$NET_GO" ]; then
    echo "The Go program $NET_GO does not exist"
    exit 1
fi


#Listing des canaux de communication
touch "$MAP_NEW"
echo "$IN_CNEW" >"$MAP_NEW"
echo "$IN_OLD" >>"$MAP_NEW"
#Ajout du nouveau site dans les sorties du parent
echo "$IN_NEW" >>"$MAP_OLD"

mapfile -t lignes < "$MAP_OLD"
nombre_lignes=${#lignes[@]}
echo "Debug : ${lignes[@]}"

PARENT_PROCS=$(pgrep -f "cat $OUT_OLD")

for PROC in $PARENT_PROCS; do
    echo "site out: $PROC"
    kill $PROC
done


# Rediriger la sortie du site original vers le site nouveau, en plus des anciens sites
if [ "$nombre_lignes" -eq 1 ]; then
    cat "$OUT_OLD" > "${lignes[0]}" &
else
    # Créer une liste des fichiers pour tee, sauf le dernier fichier

    tee_files=("${lignes[@]:0:$((nombre_lignes-1))}")
    # Le dernier fichier pour la redirection >
    output="${lignes[$((nombre_lignes-1))]}"

    # Construire la commande avec tee et >
    cat "$OUT_OLD" | tee "${tee_files[@]}" > "$output" &
    echo "Utilisation de la commande 'tee' vers les fichiers : ${tee_files[*]}"
    echo "et la redirection '>' vers $output"
fi


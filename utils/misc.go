package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var NB_SITES = 5 // Nombre de sites dans l'anneau

// Attendre un temps aléatoire
func WaitRandomTime(minMs int, maxMs int) {
	if minMs > maxMs {
		minMs, maxMs = maxMs, minMs
	}
	duration_ms := rand.Intn(maxMs-minMs) + minMs
	time.Sleep(time.Duration(duration_ms) * time.Millisecond)
}

// Convertir une liste d'entiers en chaîne de caractères
func IntSliceToString(slice []int, sep string) string {
	var strSlice []string
	for _, i := range slice {
		strSlice = append(strSlice, strconv.Itoa(i))
	}
	return strings.Join(strSlice, sep)
}

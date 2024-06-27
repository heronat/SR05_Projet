package display

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// Codes pour le terminal
const (
	CYAN    string = "\033[1;36m"
	ROUGE   string = "\033[1;31m"
	ORANGE  string = "\033[38;5;208m"
	VERT    string = "\033[1;32m"
	MAGENTA string = "\033[1;35m"
	RESET   string = "\033[0m"
)

// Récupération du PID
var pid int = os.Getpid()

// Initialisation des logs
var stderr = log.New(os.Stderr, "", 0)
var logger = log.New(os.Stderr, "", log.LstdFlags)

// Display affiche un message sur la sortie standard d'erreur
func Display(msg string) {
	stderr.Printf("%s\n", msg)
}

func Log(msg string) {
	logger.Printf("%s\n", msg)
}

func display(msg string, id int, color string, symbol string) {
	var callerFct string = "???"
	var callerFile string = "???"
	var callerLine string = "???"

	// Récupération de la ligne et du fichier
	// 2 pour remonter de 2 niveaux dans la pile d'appels
	pc, file, line, ok := runtime.Caller(2)
	if ok {
		callerFct = runtime.FuncForPC(pc).Name()
		callerFile = strings.Split(file, "/")[len(strings.Split(file, "/"))-1]
		callerLine = strconv.Itoa(line)
	}

	stderr.Printf("%s%s [%d/%d] %s:%s:%-4.4s %s%s\n", color, symbol, pid, id,
		callerFile, callerFct, callerLine, msg, RESET)
}

// Info affiche un message d'information
func Info(msg string, id int) {
	display(msg, id, CYAN, "i")
}

// Warning affiche un message d'avertissement
func Warning(msg string, id int) {
	display(msg, id, ORANGE, "*")
}

// Error affiche un message d'erreur
func Error(msg string, id int) {
	display(msg, id, ROUGE, "!")
}

// Success affiche un message de succès
func Success(msg string, id int) {
	display(msg, id, VERT, "+")
}

// Debug affiche un message de débogage
func Debug(msg string, id int) {
	display(msg, id, MAGENTA, "d")
}

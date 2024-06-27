package utils

import (
	"fmt"
	//"os"
	"strconv"
	"strings"
)

// Parse une paire clé-valeur
// >>> parseKeyValue("=Type=demandeSC")
// "Type", "demandeSC", nil
func parseKeyValue(str string) (string, string, error) {
	if len(str) < 3 {
		return "", "", fmt.Errorf("parseKeyValue: string %s too short", str)
	}
	var keyValueSeparator = str[0:1]
	field := strings.SplitN(str[1:], keyValueSeparator, 2)
	if len(field) != 2 {
		return "", "", fmt.Errorf("parseKeyValue: invalid format %s", str)
	}
	return field[0], field[1], nil
}

// Trouve les champs d'une chaîne de caractères
// >>> splitFields("/=Type=demandeSC/=Sender=0/=Receiver=1...")
// ["=Type=demandeSC", "=Sender=0", "=Receiver=1", ...]
func splitFields(
	str string, fieldSep string,
	startEscChar string, closingEscChar string,
) ([]string, error) {
	var fields []string
	var field string = ""
	var escSeqCount int = 0
	for i := 1; i < len(str); i++ {
		if str[i:i+1] == fieldSep && escSeqCount == 0 {
			fields = append(fields, field)
			field = ""
			continue
		}
		if str[i:i+1] == startEscChar {
			escSeqCount++
		} else if str[i:i+1] == closingEscChar {
			escSeqCount--
		}
		field += str[i : i+1]
	}
	fields = append(fields, field)
	if escSeqCount != 0 {
		return nil, fmt.Errorf("splitFields: invalid format %s", str)
	}
	return fields, nil
}

// Trouve la valeur d'une clé dans une chaîne de caractères
// >>> findKeyValue("/=Type=demandeSC/=Sender=0/=Receiver=1...", "Type")
// "demandeSC", nil
func findKeyValue(str string, key string) (string, error) {
	if len(str) < 3 {
		return "", fmt.Errorf("findKeyValue: string %s too short", str)
	}
	var fieldSep = str[0:1]
	var startEscChar string = "{"
	var closingEscChar string = "}"

	fields, err := splitFields(str, fieldSep, startEscChar, closingEscChar)
	if err != nil {
		return "", err
	}

	for _, field := range fields {
		k, v, err := parseKeyValue(field)
		if err != nil {
			return "", err
		}
		if k == key {
			return v, nil
		}
	}
	return "", fmt.Errorf("findKeyValue: key %s not found", key)
}

// Parse un tableau d'entiers
// >>> parseArrayInt("ArrayInt{/=0=1/=1=2/=2=3}")
// [1, 2, 3], nil
func parseArrayInt(str string) ([]int, error) {
	var array []int

	str = strings.TrimPrefix(str, "ArrayInt{")
	str = strings.TrimSuffix(str, "}")

	fields := strings.Split(str, "/")
	for _, field := range fields {
		if field == "" {
			continue
		}
		_, value, err := parseKeyValue(field)
		if err != nil {
			return nil, err
		}
		intValue, _ := strconv.Atoi(value)
		
		//fmt.Fprintln(os.Stderr, "int value  : "+ strconv.Itoa(intValue))
		array = append(array, intValue)
	}

	return array, nil
}

// Parse une matrice d'entiers
// >>> parseMatrixInt("MatrixInt{/=0=ArrayInt{/=0=1/=1=2/=2=3}/=1=ArrayInt{/=0=4/=1=5/=2=6}}")
// [[1, 2, 3], [4, 5, 6]], nil
func parseMatrixInt(str string) ([][]int, error) {
	var matrix [][]int
	
	str = strings.TrimPrefix(str, "MatrixInt{")
	str = strings.TrimSuffix(str, "}")
	
	fields, err := splitFields(str, "/", "{", "}")
	if err != nil {
		return nil, err
	}
	for _, field := range fields {
		if field == "" {
			continue
		}
		_, value, err := parseKeyValue(field)
		if err != nil {
			return nil, err
		}
		array, _ := parseArrayInt(value)
		matrix = append(matrix, array)
	}
	//fmt.Fprintln(os.Stderr, matrix)
	return matrix, nil
}

// Parse un état
// >>> parseState("State{/=SiteActions=0/=BankActions=0}")
// &State{SiteActions: 0, BankActions: 0}
func parseState(str string) (*State, error) {
	var state State

	if str == "" {
		return nil, nil
	}

	str = strings.TrimPrefix(str, "State{")
	str = strings.TrimSuffix(str, "}")

	siteActionsStr, err := findKeyValue(str, "SiteActions")
	if err != nil {
		return nil, err
	}
	state.SiteActions, _ = strconv.Atoi(siteActionsStr)

	bankActionsStr, err := findKeyValue(str, "BankActions")
	if err != nil {
		return nil, err
	}
	state.BankActions, _ = strconv.Atoi(bankActionsStr)

	return &state, nil
}

// Parse un tableau d'états
// >>> parseArrayState("ArrayState{/=0=State{/=SiteActions=0/=BankActions=0}...}")
// &[&State{SiteActions: 0, BankActions: 0}, ...], nil
func parseArrayState(str string) ([]*State, error) {
	var states []*State

	str = strings.TrimPrefix(str, "ArrayState{")
	str = strings.TrimSuffix(str, "}")

	fields, err := splitFields(str, "/", "{", "}")
	if err != nil {
		return nil, err
	}
	for _, field := range fields {
		if field == "" {
			continue
		}
		_, value, err := parseKeyValue(field)
		if err != nil {
			return nil, err
		}
		state, err := parseState(value)
		if err != nil {
			return nil, err
		}
		
		states = append(states, state)
	}
	//fmt.Fprintln(os.Stderr, states)
	return states, nil
}

// Parse un tableau de messages
// >>> parseArrayMessage("ArrayMessage{/=0=Message{/=Type=demandeSC/=Sender=0/=VectorClock=0,0,0}...}")
// &[&Message{Type: DEMANDE_SC, Sender: 0, VectorClock: [0, 0, 0], ...}, ...], nil
func parseArrayMessage(str string) ([]*Message, error) {
	var messages []*Message

	str = strings.TrimPrefix(str, "ArrayMessage{")
	str = strings.TrimSuffix(str, "}")

	fields, err := splitFields(str, "/", "{", "}")
	if err != nil {
		return nil, err
	}
	for _, field := range fields {
		if field == "" {
			continue
		}
		_, value, err := parseKeyValue(field)
		if err != nil {
			return nil, err
		}
		msg, err := ParseMessage(value)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// Parse un état global
// >>> parseGlobalState("GlobalState{/=States=ArrayState{/=0=/=1=/=2=}...}")
// &GlobalState{States: [nil, nil, nil], ...}, nil
func parseGlobalState(str string) (*GlobalState, error) {
	if str == "" {
		return nil, nil
	}
	var globalState GlobalState
	//fmt.Fprintln(os.Stderr, msg.GlobalState)
	str = strings.TrimPrefix(str, "GlobalState{")
	str = strings.TrimSuffix(str, "}")
	
	statesStr, err := findKeyValue(str, "States")
	if err != nil {
		return nil, err
	}
	states, _ := parseArrayState(statesStr)
	//fmt.Fprintln(os.Stderr, len(states))
	copy(globalState.States[:], states)

	matrixClockStr, err := findKeyValue(str, "MatrixClock")
	if err != nil {
		return nil, err
	}
	matrixClock, _ := parseMatrixInt(matrixClockStr)
	//fmt.Fprintln(os.Stderr, matrixClock)
	globalState.MatrixClock = make([][]int, NB_SITES)
	for i := 0; i < NB_SITES; i++ {
		globalState.MatrixClock[i] = make([]int, len(matrixClock[i]))
		//fmt.Fprintln(os.Stderr,len(matrixClock[i]))
		copy(globalState.MatrixClock[i][:], matrixClock[i])
	}

	prepostMessagesStr, err := findKeyValue(str, "PrepostMessages")
	if err != nil {
		return nil, err
	}
	prepostMessages, _ := parseArrayMessage(prepostMessagesStr)
	globalState.PrepostMessages = prepostMessages

	return &globalState, nil
}

// Parse un message
// >>> ParseMessage("Message{/=Type=demandeSC/=Sender=0/=VectorClock=VectorClock{/=0=0/=1=0/=2=0}...}")
// &Message{Type: DEMANDE_SC, Sender: 0, VectorClock: [0, 0, 0], ...}, nil
func ParseMessage(str string) (*Message, error) {
	var msg Message

	str = strings.TrimPrefix(str, "Message{")
	str = strings.TrimSuffix(str, "}")

	typeStr, err := findKeyValue(str, "Type")
	if err != nil {
		return nil, err
	}
	msg.Type = MessageType(typeStr)

	senderStr, err := findKeyValue(str, "Sender")
	if err != nil {
		return nil, err
	}
	msg.Sender, _ = strconv.Atoi(senderStr)

	receiverStr, err := findKeyValue(str, "Receiver")
	if err != nil {
		return nil, err
	}
	msg.Receiver, _ = strconv.Atoi(receiverStr)

    transmitterStr, err := findKeyValue(str, "Transmitter")
	if err != nil {
		return nil, err
	}
	msg.Transmitter, _ = strconv.Atoi(transmitterStr)

    parentStr, err := findKeyValue(str, "Parent")
	if err != nil {
		return nil, err
	}
	msg.Parent, _ = strconv.Atoi(parentStr)

	stampStr, err := findKeyValue(str, "Stamp")
	if err != nil {
		return nil, err
	}
	msg.Stamp, _ = strconv.Atoi(stampStr)

	vectClockStr, err := findKeyValue(str, "VectorClock")
	if err != nil {
		return nil, err
	}
	vectClock, _ := parseArrayInt(vectClockStr)
	msg.VectorClock = make([]int, len(vectClock))
	copy(msg.VectorClock[:], vectClock)


	siteListStr, err := findKeyValue(str, "SiteList")
	if err != nil {
		return nil, err
	}
	siteList, _ := parseArrayInt(siteListStr)
	//fmt.Fprintln(os.Stderr, " len sitelistavcopie : "+ strconv.Itoa(len(siteList)))
	//fmt.Fprintln(os.Stderr, "len sitelistaprescopie : "+ strconv.Itoa(len(msg.SiteList)))
	msg.SiteList = make([]int, len(siteList))
	copy(msg.SiteList[:], siteList)
	
	colorStr, err := findKeyValue(str, "Color")
	if err != nil {
		return nil, err
	}
	msg.Color = Color(colorStr)

	stateStr, err := findKeyValue(str, "State")
	if err != nil {
		return nil, err
	}
	state, _ := parseState(stateStr)
	msg.State = state

	globalStateStr, err := findKeyValue(str, "GlobalState")
	if err != nil {
		return nil, err
	}
	globalState, _ := parseGlobalState(globalStateStr)
	msg.GlobalState = globalState

	balanceStr, err := findKeyValue(str, "Balance")
	if err != nil {
		return nil, err
	}
	balance, _ := strconv.Atoi(balanceStr)
	msg.Balance = balance

	eluStr, err := findKeyValue(str, "Elu")
	if err != nil {
		return nil, err
	}
	msg.Elu, _ = strconv.Atoi(eluStr)

	return &msg, nil
}

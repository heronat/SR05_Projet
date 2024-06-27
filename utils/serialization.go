package utils
//import "fmt"
//import "os"
import "strconv"
//import "fmt"
//import "os"
// Sérialise une paire clé-valeur
// >>> serializeKeyValue("Type", "demandeSC")
// "/=Type=demandeSC"
func serializeKeyValue(key string, value string) string {
	var fieldSeparator = "/"
	var keyValueSeparator = "="
	return fieldSeparator + keyValueSeparator + key + keyValueSeparator + value
}

// Sérialise un tableau d'entiers
// >>> serializeArrayInt([1, 2, 3])
// "ArrayInt{/=0=1/=1=2/=2=3}"
func serializeArrayInt(array []int) string {
	var str string = "ArrayInt{"
	//fmt.Fprintln(os.Stderr, array)
	for i, value := range array {
		str += serializeKeyValue(strconv.Itoa(i), strconv.Itoa(value))
	}
	str += "}"
	//fmt.Fprintln(os.Stderr, str)
	return str
}

// Sérialise une matrice d'entiers
// >>> serializeMatrixInt([[1, 2, 3], [4, 5, 6]])
// "MatrixInt{/=0=ArrayInt{/=0=1/=1=2/=2=3}/=1=ArrayInt{/=0=4/=1=5/=2=6}}"
func serializeMatrixInt(matrix [][]int) string {
	var str string = "MatrixInt{"
	for i, row := range matrix {
		str += serializeKeyValue(strconv.Itoa(i), serializeArrayInt(row))
	}
	str += "}"
	return str
}

// Sérialise un état
// >>> (&State{SiteActions: 0, BankActions: 0}).serializeState()
// "State{/=SiteActions=0/=BankActions=0}"
func (state *State) serializeState() string {
	if state == nil {
		return ""
	}
	var str string = "State{"
	str += serializeKeyValue("SiteActions", strconv.Itoa(state.SiteActions))
	str += serializeKeyValue("BankActions", strconv.Itoa(state.BankActions))
	str += "}"
	return str
}

// Sérialise un tableau d'états
// >>> serializeArrayState(&[NB_SITES]*State{&State{SiteActions: 0, BankActions: 0}, ...})
// "ArrayState{/=0=State{/=SiteActions=0/=BankActions=0}/=1=State{/=SiteActions=0/=BankActions=0}...}"
func serializeArrayState(states []*State) string {
	var str string = "ArrayState{"
	for i, state := range states {
		str += serializeKeyValue(strconv.Itoa(i), state.serializeState())
	}
	str += "}"
	return str
}

// Sérialise un tableau de messages
// >>> serializeArrayMessage(&[NB_SITES]*Message{&Message{Type: DEMANDE_SC, Sender: 0, VectorClock: [0, 0, 0], ...}, ...})
// "ArrayMessage{/=0=Message{/=Type=demandeSC/=Sender=0/=VectorClock=0,0,0}...}"
func serializeArrayMessage(messages []*Message) string {
	var str string = "ArrayMessage{"
	for i, msg := range messages {
		str += serializeKeyValue(strconv.Itoa(i), msg.SerializeMessage())
	}
	str += "}"
	return str
}

// Sérialise un état global
// >>> (&GlobalState{States: [nil, nil, nil], MatrixClock: [[0, 0], [0, 0]]}).serializeGlobalState()
// "GlobalState{/=States=ArrayState{/=0=/=1=/=2=}/=MatrixClock{/=0=ArrayInt{/=0=0/=1=0}/=1=ArrayInt{/=0=0/=1=0}}}"
func (globalS *GlobalState) serializeGlobalState() string {
	if globalS == nil {
		return ""
	}
	var str string = "GlobalState{"
	str += serializeKeyValue("States", serializeArrayState(globalS.States[:]))
	var matrixInt [][]int
	for i := 0; i < NB_SITES; i++ {
		matrixInt = append(matrixInt, globalS.MatrixClock[i][:])
		//fmt.Fprintln(os.Stderr, matrixInt)
	}
	str += serializeKeyValue("MatrixClock", serializeMatrixInt(matrixInt))
	str += serializeKeyValue("PrepostMessages", serializeArrayMessage(globalS.PrepostMessages))
	str += "}"
	return str
}

// Sérialise un message
// >>> (&Message{Type: DEMANDE_SC, Sender: 0, VectorClock: [0, 0, 0], ...}).SerializeMessage()
// "Message{/=Type=demandeSC/=Sender=0/=VectorClock=0,0,0/=...}"
func (msg *Message) SerializeMessage() string {
	var str string = "Message{"
	str += serializeKeyValue("Type", string(msg.Type))
	//fmt.Fprintln(os.Stderr, msg.Type)
	str += serializeKeyValue("Sender", strconv.Itoa(msg.Sender))
	str += serializeKeyValue("Receiver", strconv.Itoa(msg.Receiver))
    str += serializeKeyValue("Transmitter", strconv.Itoa(msg.Transmitter))
    str += serializeKeyValue("Parent", strconv.Itoa(msg.Parent))
	str += serializeKeyValue("Stamp", strconv.Itoa(msg.Stamp))
	str += serializeKeyValue("VectorClock", serializeArrayInt(msg.VectorClock[:]))
	str += serializeKeyValue("SiteList", serializeArrayInt(msg.SiteList[:]))
	str += serializeKeyValue("Color", string(msg.Color))
	str += serializeKeyValue("State", msg.State.serializeState())
	str += serializeKeyValue("GlobalState", msg.GlobalState.serializeGlobalState())
	str += serializeKeyValue("Balance", strconv.Itoa(msg.Balance))
    str += serializeKeyValue("Elu", strconv.Itoa(msg.Elu))
	str += "}"
	//fmt.Fprintln(os.Stderr, str)
	return str
}

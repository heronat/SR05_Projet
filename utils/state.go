package utils

type State struct {
	SiteActions int
	BankActions int
}

type GlobalState struct {
	States          []*State        // États des sites
	MatrixClock     [][]int // Tableau d'horloges vectorielles
	PrepostMessages []*Message              // Messages prepost
}

// Initialisation de l'état
func NewState(siteActions int, bankActions int) *State {
	return &State{SiteActions: siteActions, BankActions: bankActions}
}

// Initialisation de l'état global
func NewGlobalState() *GlobalState {
	var states []*State
	var matrixClock [][]int
	var prepostMessages []*Message
	for i := 0; i < NB_SITES; i++ {
		states = append(states, nil)
		matrixClock = append(matrixClock, make([]int, NB_SITES))
	} // Initialisation de la matrice d'horloges
	prepostMessages = nil
	/*for i := 0; i < NB_SITES; i++ {
		states[i] = nil
		for j := 0; j < NB_SITES; j++ {
			matrixClock[i][j] = 0
		} // Initialisation de la matrice d'horloges
	}
	prepostMessages = nil*/

	return &GlobalState{
		States:          states,
		MatrixClock:     matrixClock,
		PrepostMessages: prepostMessages,
	}
}

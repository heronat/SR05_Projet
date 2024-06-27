package utils

import "fmt"
//import "os"
// Type énuméré pour les types de messages
type MessageType string

const (
	// Algorithme de file d'attente répartie
	DEMANDE_SC MessageType = "demandeSC"
	FIN_SC     MessageType = "finSC"
	DEBUT_SC   MessageType = "debutSC"
	LIBERATION MessageType = "liberation"
	ACCUSE     MessageType = "accuse"
	REQUETE    MessageType = "requete"

	// Réplication de données
	ACTUALISATION MessageType = "actualisation"
    ACTUALISATION_GLOBAL_STATE MessageType = "actualisationGlobalState"
    ACTUALISATION_INSTANTANE MessageType = "actualisationInstantane"

	// Construction d'instantanés
	DEMANDE_SNAPSHOT MessageType = "demandeSnapshot"
	FIN_SNAPSHOT     MessageType = "finSnapshot"
	ETAT             MessageType = "etat"
	PREPOST          MessageType = "prepost"
	SIGNAL_SNAPSHOT  MessageType = "signal"
    REVERT_INSTANTANE_APP MessageType = "revertInstantaneApp"
    REVERT_INSTANTANE_CTL MessageType = "revertInstantaneCtl"

	// Gestion du Membership 
	AJOUT MessageType = "ajout"
	SUPPRESSION MessageType = "suppression"
    DEMANDE_AJOUT MessageType = "demandeAjout"
    ACCEPTATION_AJOUT MessageType = "acceptationAjout"
    PREVENTION_VOISINS MessageType = "preventionVoisins"
	NBSITESTRANSMISSION MessageType = "nbSitesTransmission"

    // Élection
	DEMANDE_ADMISSION MessageType = "demandeAdmission"
	DEMANDE_DEPART    MessageType = "demandeDepart"
	ELECTION_BLEU     MessageType = "electionBleu"  // Propagation de la candidature
	ELECTION_ROUGE    MessageType = "electionRouge" // Confirmation des résultats


	//Echanges CTL NET
	LISTE_CTL MessageType = "listeCtl"
	DEMANDE_STATE MessageType = "demandeState"
	RETOUR_STATE MessageType = "retourState"
    SUPPRESSION_CTL MessageType = "suppressionCtl"

	//Echange CTL APP
	NBSITESAPP MessageType = "nbSitesApp"
	NEWSITE_APP MessageType = "newSiteApp"
    DELETESITE_APP MessageType = "deleteSiteApp"
)

// Type énuméré pour les couleurs des messages
type Color string

const (
	BLANC Color = "blanc"
	ROUGE Color = "rouge"
)

type Message struct {
	Type        MessageType   // Type de message
	Sender      int           // Site émetteur
	Receiver    int           // Site destinataire
	Transmitter int           // Site qui nous transmet le message
	Parent      int           // Parent du site qui nous transmet le message
	Stamp       int           // Estampille logique
	VectorClock []int // Horloge vectorielle
	SiteList   []int         // Liste des sites
	Color       Color         // Couleur du message
	State       *State        // État du site
	GlobalState *GlobalState  // État global
	Balance     int           // Différence entre le nombre de messages envoyés et reçus
	Elu         int           // Site elu ELECTION
}

// Créer un message
// >>> NewMessage(DEMANDE_SC, 0, 1, 0, ...)
// &Message{Type: DEMANDE_SC, Sender: 0, Receiver: 1, Stamp: 0, ...}
func NewMessage(
	Type MessageType, Sender int, Receiver int, Transmitter int, Parent int,
	Stamp int, VectorClock []int, SiteList []int, Color Color,
	State *State, GlobalState *GlobalState, Balance int, Elu int,
) *Message {
	return &Message{
		Type:        Type,
		Sender:      Sender,
		Receiver:    Receiver,
        Transmitter: Transmitter,
        Parent:      Parent,
		Stamp:       Stamp,
		VectorClock: VectorClock,
		SiteList:    SiteList,
		Color:       Color,
		State:       State,
		GlobalState: GlobalState,
		Balance:     Balance,
        Elu:         Elu,
	}
}

// Envoyer un message
func SendMessage(msg *Message) {
	// WaitRandomTime(100, 1000) // Simulation d'un traitement long
	fmt.Println(msg.SerializeMessage())
	//fmt.Fprintln(os.Stderr, msg.SerializeMessage())
}

// Recevoir un message
func ReceiveMessage() (*Message, error) {
	var rawMsg string
	_, err := fmt.Scanln(&rawMsg)
	
	if err != nil {
		return nil, err
	}
	// WaitRandomTime(100, 1000) // Simulation d'un traitement long
	return ParseMessage(rawMsg)
}

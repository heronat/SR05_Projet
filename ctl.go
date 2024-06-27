package main

import (
	"flag"
	"src/display"
	. "src/utils"
	"strconv"
	"time"
)

var sitelist []int

func (ctl *Controller) constructionInstantane(msg *Message) {
	i := ctl.Id
	hi := ctl.Stamp
	Vi := ctl.VectorClock
	couleur := ctl.Color
	initiateur := ctl.Initiate
	bi := ctl.Balance
	EGi := ctl.GlobalState
	NbEtatsAttendus := ctl.NbStates
	NbMsgAttendus := ctl.NbMsgs

	j := msg.Sender
	c := msg.Color
	s := msg.State
	EG := msg.GlobalState
	bj := msg.Balance

	switch msg.Type {
	case DEMANDE_SNAPSHOT:
		display.Info("Demande d'instantané reçue", i)
		couleur = ROUGE
		initiateur = true
		EGi.States[i] = s
		EGi.MatrixClock[i] = Vi
		bi = 0
		NbEtatsAttendus = NB_SITES - 1
		NbMsgAttendus = bi
		//Envoyer un signal aux autres sites pour les prévenir de la demande d'instantane
		for _, j := range sitelist {
			if j != i {
				display.Info("Envoi d'un signal d'instantané à "+strconv.Itoa(j), i)
				SendMessage(NewMessage(SIGNAL_SNAPSHOT, i, j, -1, -1, hi, Vi, nil, couleur, nil, nil, 1, -1))
			}
		}
	case ETAT:
		display.Info("Réception d'un état local et d'un bilan", i)
		//cette partie n'est traitée que par le site qui a effectué la demande de snapshot
		if initiateur {
			display.Info("Collecte des états des autres sites", i)
			display.Info("Quantité d'action du site "+strconv.Itoa(j)+" : "+strconv.Itoa(EG.States[j].SiteActions), i)
			display.Info("Quantité d'action de la banque "+strconv.Itoa(j)+" : "+strconv.Itoa(EG.States[j].BankActions), i)
			EGi.States[j] = EG.States[j]
			EGi.MatrixClock[j] = EG.MatrixClock[j]
			NbEtatsAttendus--
			NbMsgAttendus = NbMsgAttendus + bj
			if NbEtatsAttendus == 0 && NbMsgAttendus == 0 {
				display.Success("Fin de l'instantané", i)
				couleur = BLANC
				for _, k := range sitelist {
					display.Info("Quantité d'action du site "+strconv.Itoa(k)+" : "+strconv.Itoa(EGi.States[k].SiteActions), i)
					display.Info("Quantité d'action de la banque "+strconv.Itoa(k)+" : "+strconv.Itoa(EGi.States[k].BankActions), i)
					display.Info("Horloge vectorielle du site "+strconv.Itoa(k)+" : "+IntSliceToString(EGi.MatrixClock[k][:], ","), i)
					//On envoie un message à tous les controleurs pour les prévenir de la fin de la snapshot
					SendMessage(NewMessage(FIN_SNAPSHOT, i, k, -1, -1, hi, Vi, nil, BLANC, nil, EGi, -1, -1))
				}
				display.Info("Nombre de messages prépost : "+strconv.Itoa(len(EGi.PrepostMessages)), i)
			}
		} else {
			// envoyer([état] EG, bilan) sur l'anneau
			display.Info("Retransmission de l'état local et du bilan", i)
			SendMessage(msg)
		}
	case PREPOST:
		// Réception d'un message prépost retransmis sur l'anneau par son destinataire pour l'initiateur
		display.Info("Réception d'un message prépost", i)
		if initiateur {
			display.Info("Collecte des messages prépost des autres sites", i)
			NbMsgAttendus--
			EGi.PrepostMessages = append(EGi.PrepostMessages, EG.PrepostMessages...)
			display.Info("NbEtatsAttendus "+strconv.Itoa(NbEtatsAttendus)+" NbMsgAttendus "+strconv.Itoa(NbMsgAttendus), i)
			if NbEtatsAttendus == 0 && NbMsgAttendus == 0 {
				display.Success("Fin de l'instantané", i)
				couleur = BLANC
				for _, k := range sitelist {
					display.Info("Quantité d'action du site "+strconv.Itoa(k)+" : "+strconv.Itoa(EGi.States[k].SiteActions), i)
					display.Info("Quantité d'action de la banque "+strconv.Itoa(k)+" : "+strconv.Itoa(EGi.States[k].BankActions), i)
					display.Info("Horloge vectorielle du site "+strconv.Itoa(k)+" : "+IntSliceToString(EGi.MatrixClock[k][:], ","), i)
					//On envoie un message à tous les controleurs pour les prévenir de la fin de la snapshot
					SendMessage(NewMessage(FIN_SNAPSHOT, i, k, -1, -1, hi, Vi, nil, BLANC, nil, EGi, -1, -1))
				}
			}
		} else {
			SendMessage(NewMessage(PREPOST, i, -1, -1, -1, hi, Vi, nil, couleur, s, EGi, -1, -1))
		}
	case FIN_SNAPSHOT:
		display.Info("Fin de l'instantané reçu, remise à zéro des variables", i)
		//Quand l'instantane est finie, on diffuse la sauvegarde créée à l'application et on réinitialise les valeurs pour une prochaine snapshot
		couleur = BLANC
		initiateur = false
		EGi = EG
		SendMessage(NewMessage(ACTUALISATION_GLOBAL_STATE, i, i, -1, -1, hi, Vi, nil, couleur, nil, EGi, -1, -1))

	case REVERT_INSTANTANE_APP:
		display.Info("Demande de revert d'instantane reçue de la part de l'application", i)
		for _, k := range sitelist {
			SendMessage(NewMessage(REVERT_INSTANTANE_CTL, i, k, -1, -1, hi, Vi, nil, "", nil, EG, -1, -1))
		}

	case REVERT_INSTANTANE_CTL:
		display.Info("Envoi de l'instantane à reset à l'application", i)
		SendMessage(NewMessage(ACTUALISATION_INSTANTANE, i, i, -1, -1, hi, Vi, nil, "", nil, EG, -1, -1))

	default:
		// Réception d'un message de l'application de base
		display.Info("Réception d'un message de l'application de base", i)
		if msg.Type == SIGNAL_SNAPSHOT {
			display.Info("Réception d'un signal d'instantané", i)
			bi = 0
		} else {
			bi--
		}
		if c == ROUGE && couleur == BLANC {
			display.Info("Première réception d'un message rouge", i)
			couleur = ROUGE
			// EGi.States[i] est l'état local du site i mis à jour par l'application
			// de base lors de la réception d'un message de fin de SC
			display.Info("Quantité d'action du site "+strconv.Itoa(i)+" : "+strconv.Itoa(EGi.States[i].SiteActions), i)
			display.Info("Quantité d'action de la banque "+strconv.Itoa(i)+" : "+strconv.Itoa(EGi.States[i].BankActions), i)
			EGi.MatrixClock[i] = Vi
			display.Info("Envoi d'un état local et d'un bilan"+strconv.Itoa(bi), i)
			SendMessage(NewMessage(ETAT, i, -1, -1, -1, hi, Vi, nil, couleur, nil, EGi, bi, -1))
		}
		if c == BLANC && couleur == ROUGE {
			display.Info("Réception postclic d'un message envoyé préclic", i)
			EGi.PrepostMessages = append(EGi.PrepostMessages, msg)
			EGi.States[i] = msg.State
			SendMessage(NewMessage(PREPOST, i, -1, -1, -1, hi, Vi, nil, couleur, nil, EGi, -1, -1))
		}
	}

	//on sauvegarde les valeurs modifiées dans le controleur
	ctl.Stamp = hi
	ctl.VectorClock = Vi
	ctl.Color = couleur
	ctl.Initiate = initiateur
	ctl.Balance = bi
	ctl.GlobalState = EGi
	ctl.NbStates = NbEtatsAttendus
	ctl.NbMsgs = NbMsgAttendus
}

//fonction qui vérifie si un message possède l'horloge la plus petite parmis toutes celles du site
func minTabRequete(tab []struct {
	Type  MessageType
	Stamp int
}, i int) bool {
	if tab[i].Type != REQUETE {
		return false
	}
	for tmp := 0; tmp < NB_SITES; tmp++ {
        if tab[tmp].Stamp > -1 {
            if tab[i].Stamp > tab[tmp].Stamp {
                return false
            } else if tab[i].Stamp == tab[tmp].Stamp && i > tmp {
                return false
            }
        }
	}
	return true
}

func (ctl *Controller) fileAttenteRepartie(msg *Message) {
	i := ctl.Id
	tab := ctl.QueueCS
	hi := ctl.Stamp
	Vi := ctl.VectorClock
	h := msg.Stamp
	Vm := msg.VectorClock
	j := msg.Sender

	//EGi := ctl.GlobalState
	couleur := ctl.Color
	bilan := ctl.Balance

	switch msg.Type {

	case DEMANDE_SC:
		display.Info("Demande de SC reçue", i)
		hi++
		display.Info("vecteur traite : "+IntSliceToString(Vi, ","), i)
		display.Info("ctl correspondant : "+IntSliceToString(ctl.VectorClock, ","), i)
		Vi[i]++
		tab[i] = struct {
			Type  MessageType
			Stamp int
		}{Type: REQUETE, Stamp: hi}
		// envoyer([requête] hi) à tous les autres sites
		display.Info("[DEBUG] PRELIST", i)
		for _, k := range sitelist {
			if k != i {
				display.Info("Envoi d'une requête à "+strconv.Itoa(k), i)
				bilan++
				SendMessage(NewMessage(REQUETE, i, k, -1, -1, hi, Vi, nil, couleur, nil, nil, -1, -1))
			}
		}
	case FIN_SC:
		display.Info("Fin de SC reçue", i)
		//Si on reçoit une fin de section, cela implique généralement que des valeurs d'actions ont changé dans notre application, donc on sauvegarde le nouvel état localement
		ctl.GlobalState.States[i] = msg.State
		hi++
		Vi[i]++
		tab[i] = struct {
			Type  MessageType
			Stamp int
		}{Type: LIBERATION, Stamp: ctl.Stamp}
		// envoyer([libération] hi) à tous les autres sites
		for _, k := range sitelist {
			if k != i {
				display.Info("Envoi d'une libération à "+strconv.Itoa(k), i)
				bilan++
				SendMessage(NewMessage(LIBERATION, i, k, i, -1, hi, Vi, nil, couleur, msg.State, nil, -1, -1))
			}
		}
	case REQUETE:
		display.Debug("Requête reçue", i)
		hi = MaxStamp(hi, h) + 1
		Vi = MaxVectorClock(Vi, Vm)
		Vi[i]++
		tab[j] = struct {
			Type  MessageType
			Stamp int
		}{Type: REQUETE, Stamp: h}
		display.Info("Envoi d'un accusé de réception à "+strconv.Itoa(j), i)
		bilan++
		SendMessage(NewMessage(ACCUSE, i, j, -1, -1, hi, Vi, nil, couleur, nil, nil, 1, -1))
		// L'arrivée du message pourrait permettre de satisfaire une éventuelle demande de Si
		if minTabRequete(tab, i) {
			for k := range tab {
				display.Info("Indice du tableau "+strconv.Itoa(k)+" valeur de l'estampille: "+strconv.Itoa(tab[k].Stamp), i)
			}
			// envoyer([débutSC]) à l'application de base
			display.Info("Envoi d'un début de SC à l'application de base", i)
			bilan++
			SendMessage(NewMessage(DEBUT_SC, i, i, -1, -1, hi, Vi, nil, couleur, nil, nil, -1, -1))
		}
	case LIBERATION:
		display.Info("Libération reçue", i)
		hi = MaxStamp(hi, h) + 1
		Vi = MaxVectorClock(Vi, Vm)
		Vi[i]++
		tab[j] = struct {
			Type  MessageType
			Stamp int
		}{Type: LIBERATION, Stamp: h}
		bilan++
		ctl.GlobalState.States[i].BankActions = msg.State.BankActions
		SendMessage(NewMessage(ACTUALISATION, i, i, -1, -1, hi, Vi, nil, couleur, msg.State, nil, -1, -1))
		if minTabRequete(tab, i) {
			for k := range tab {
				display.Info("Indice du tableau "+strconv.Itoa(k)+" valeur de l'estampille: "+strconv.Itoa(tab[k].Stamp), i)
			}
			display.Info("Envoi d'un début de SC à l'application de base", i)
			SendMessage(NewMessage(DEBUT_SC, i, i, -1, -1, hi, Vi, nil, couleur, nil, nil, -1, -1))
		}
	case ACCUSE:
		display.Info("Accusé de réception reçu", i)
		hi = MaxStamp(hi, h) + 1
		Vi = MaxVectorClock(Vi, Vm)
		Vi[i]++
		if tab[j].Type != REQUETE {
			tab[j] = struct {
				Type  MessageType
				Stamp int
			}{Type: ACCUSE, Stamp: h}
		}

		for k := range tab {
			display.Debug("Indice du tableau "+strconv.Itoa(k)+" valeur de l'estampille: "+strconv.Itoa(tab[k].Stamp), i)
		}
		if minTabRequete(tab, i) {
			for k := range tab {
				display.Info("Indice du tableau "+strconv.Itoa(k)+" valeur de l'estampille: "+strconv.Itoa(tab[k].Stamp), i)
			}
			display.Info("Envoi d'un début de SC à l'application de base", i)
			bilan++
			SendMessage(NewMessage(DEBUT_SC, i, i, -1, -1, hi, Vi, nil, couleur, nil, nil, -1, -1))
		}
	}

	ctl.Stamp = hi
	ctl.VectorClock = Vi
	ctl.QueueCS = tab

	ctl.Color = couleur
	ctl.Balance = bilan
}

func (ctl *Controller) horlogeLogique(msg *Message) {
	i := ctl.Id
	hi := ctl.Stamp
	Vi := ctl.VectorClock
	h := msg.Stamp
	Vm := msg.VectorClock

	display.Info("Mise à jour de l'horloge logique", i)

	display.Debug("Estampille actuelle : "+strconv.Itoa(hi), i)
	hi = MaxStamp(hi, h) + 1
	display.Debug("Estampille mise à jour : "+strconv.Itoa(hi), i)

	display.Debug("Horloge vectorielle actuelle controleur: "+IntSliceToString(Vi[:], ","), i)
	display.Debug("Horloge vectorielle message : "+IntSliceToString(Vm[:], ","), i)
	Vi = MaxVectorClock(Vi, Vm)
	//Vi[i]++
	display.Debug("Horloge vectorielle mise à jour : "+IntSliceToString(Vi[:], ","), i)

	ctl.Stamp = hi
	ctl.VectorClock = Vi
}

type Controller struct {
	Id int // Identifiant du site

	Stamp   int // Horloge logique
	QueueCS []struct {
		Type  MessageType
		Stamp int
	} // Tableau des requêtes

	Color       Color
	Initiate    bool         // Initiateur
	Balance     int          // Différence entre le nombre de messages envoyés et reçus
	GlobalState *GlobalState // Etat global
	NbStates    int          // Nombre d'états attendus
	NbMsgs      int          // Nombre de messages attendus

	VectorClock []int // Horloge vectorielle
}

func NewController(i int) *Controller {
	tab := []struct {
		Type  MessageType
		Stamp int
	}{}
	// Initialisation du tableau des requêtes
	// Ajouter {LIBERATION, -1} pour chaque i
	// Cette valeur sera écrasée à la première demande de section critique sauf si un site ne fait plus parti du réseau
	for i := 0; i < NB_SITES; i++ {
		tab = append(tab, struct {
			Type  MessageType
			Stamp int
		}{Type: LIBERATION, Stamp: -1})
	}

	vectClock := make([]int, NB_SITES)

	return &Controller{
		Id:          i,
		Stamp:       0,
		QueueCS:     tab,
		Color:       BLANC,
		Initiate:    false,
		Balance:     0,
		GlobalState: NewGlobalState(),
		NbStates:    0,
		NbMsgs:      0,
		VectorClock: vectClock,
	}
}

func main() {
	p_id := flag.Int("id", -1, "Identifiant de l'application de base")
	flag.Parse()
	ctl := NewController(*p_id)
	varIteration := -1
	if *p_id > 4 {
		for {
			msg, err := ReceiveMessage()
			display.Debug("Reception d'un message....", *p_id)
			if err != nil {
				display.Error("Le format du message reçu est incorrect", *p_id)
				time.Sleep(5 * time.Second)
				continue
			}
			/*if msg.Type == NBSITESTRANSMISSION{
				NB_SITES = msg.Balance
				display.Success("Recepetion du nb de sites du net :"+strconv.Itoa(NB_SITES), *p_id)
				ctl.VectorClock = make([]int, len(msg.VectorClock))
				copy(ctl.VectorClock, msg.VectorClock)
				break
			}*/

			if msg.Type == LISTE_CTL && msg.Sender == *p_id && msg.Receiver == *p_id {
				display.Info("Liste des controleurs reçue", *p_id)
				display.Info("Liste des controleurs : "+IntSliceToString(msg.SiteList, ","), *p_id)

				NB_SITES = len(msg.VectorClock)
				ctl = NewController(*p_id)
				sitelist = make([]int, len(msg.SiteList))
				copy(sitelist[:], msg.SiteList)
				display.Info("Vector clock msg: "+IntSliceToString(msg.VectorClock, ","), *p_id)
				ctl.VectorClock = make([]int, len(msg.VectorClock))
				copy(ctl.VectorClock, msg.VectorClock)
				display.Info("ctl vectorclock : "+IntSliceToString(ctl.VectorClock, ","), *p_id)
				varIteration = len(msg.SiteList) - 1
				display.Error("Nb de retour attendu : "+strconv.Itoa(varIteration), *p_id)

				//On partage nbsites a l'app
				SendMessage(NewMessage(NBSITESAPP, *p_id, *p_id, -1, -1, 0, nil, msg.SiteList, "", nil, nil, NB_SITES, -1))
				display.Success("Envoi du nb de sites à l'app", *p_id)

				for _, k := range sitelist {
					if k != *p_id {
						display.Info("Demande du state de l'app pour l'envoyer à la nouvelle app", *p_id)
						SendMessage(NewMessage(DEMANDE_STATE, *p_id, k, -1, -1, 0, ctl.VectorClock[:], sitelist, "", nil, nil, 0, -1))
					}
				}
				break
			}
		}
		if *p_id < 0 || *p_id > NB_SITES {
			display.Error("L'identifiant doit être compris entre 0 et "+
				strconv.Itoa(NB_SITES), -1)
			return
		}
		for k := 0; k < NB_SITES; k++ {
			ctl.GlobalState.States[k] = NewState(20, 100)
			//display.Debug("Le site "+strconv.Itoa(k)+" a "+strconv.Itoa(ctl.GlobalState.States[k].SiteActions)+" actions de site et "+strconv.Itoa(ctl.GlobalState.States[k].BankActions)+" actions de banque", *p_id)
		}

		// On separe en deux boucles for pour etre sur qu'on ne recoit pas de retour state AVANT LISTE_CTL
		for {
			msg, err := ReceiveMessage()
			if err != nil {
				display.Error("Le format du message reçu est incorrect", *p_id)
				time.Sleep(5 * time.Second)
				continue
			}
			if msg.Type == RETOUR_STATE {
				display.Success("Retour d'état reçu de "+strconv.Itoa(msg.Sender), *p_id)
				varIteration--
				if varIteration == 0 {
					SendMessage(NewMessage(ACTUALISATION, ctl.Id, ctl.Id, -1, -1, 0, ctl.VectorClock, nil, "", msg.State, nil, -1, -1))
					//display.Success("On a recu un retour d'etat", *p_id)
					//display.Debug("[DEBUG] nombre de sites dans tabrequete: "+strconv.Itoa(len(ctl.QueueCS)), *p_id)
					//display.Debug("[DEBUG] nombre de sites dans newglobalstate: "+strconv.Itoa(len(ctl.GlobalState.States)), *p_id)
					//display.Debug("[DEBUG] nombre de sites dans newglobalstatematrix: "+strconv.Itoa(len(ctl.GlobalState.MatrixClock)), *p_id)
					break
				}

				continue
			}

		}
	}

	for k := 0; k < NB_SITES; k++ {
		ctl.GlobalState.States[k] = NewState(20, 100)
		//display.Debug("Le site "+strconv.Itoa(k)+" a "+strconv.Itoa(ctl.GlobalState.States[k].SiteActions)+" actions de site et "+strconv.Itoa(ctl.GlobalState.States[k].BankActions)+" actions de banque", *p_id)
	}

	// Initialisation du contrôleur
	// Au debut tout le monde a 20 actions de site et 100 actions de banque

	for {
		// Attente d'un message
		//display.Debug("Entrez un message", *p_id)
		msg, err := ReceiveMessage()
		//display.Debug("Vi reçu: " + IntSliceToString(msg.VectorClock, ","), *p_id)
		if err != nil {
			display.Error("Le format du message reçu est incorrect", *p_id)
			time.Sleep(5 * time.Second)
			continue
		}
		if msg.Type == LISTE_CTL && msg.Receiver != *p_id {
			display.Info("Liste des controleurs pour qqn dautres", *p_id)
			continue
		}
		if msg.Receiver != *p_id && // Le message n'est pas pour moi
			msg.Receiver != -1 { // ni pour tout le monde
			//display.Info("Message non destiné à ce site : "+strconv.Itoa(msg.Receiver), *p_id)
			continue
		}

		ctl.horlogeLogique(msg)         // Mise à jour de l'horloge logique
		ctl.constructionInstantane(msg) // Constructions d'instantanés
		if msg.Type == LISTE_CTL && msg.Sender == *p_id && msg.Receiver == *p_id {
			display.Info("Liste des controleurs reçue", *p_id)
			display.Info("Liste des controleurs : "+IntSliceToString(msg.SiteList, ","), *p_id)
			sitelist = make([]int, len(msg.SiteList))
			copy(sitelist[:], msg.SiteList)
			continue
		}
		if msg.Type == DEMANDE_STATE {

			// 0 par défaut, à corriger pour avoir l'estampille du nouveau site add
			ctl.VectorClock = append(ctl.VectorClock, 0)


			//On ajoute le nouveau state du nouveau site
			state := NewState(20, 100)
			ctl.GlobalState.States = append(ctl.GlobalState.States, state)

			//On ajoute l'estampille du novueau site dans les matrices saved
			for _, k := range sitelist {
				ctl.GlobalState.MatrixClock[k] = append(ctl.GlobalState.MatrixClock[k], 0)
			}

			//On ajoute le vecteur du nouveau site dans les matrices saved
			ctl.GlobalState.MatrixClock = append(ctl.GlobalState.MatrixClock, msg.VectorClock)
			sitelist = make([]int, len(msg.SiteList))
			copy(sitelist, msg.SiteList)
			NB_SITES += 1
			// On ajoute la case du nouveau site dans queueCS
			ctl.QueueCS = append(ctl.QueueCS, struct {
				Type  MessageType
				Stamp int
			}{Type: LIBERATION, Stamp: 0})

			/*display.Debug("[DEBUG] nombre de sites dans tabrequete: "+ strconv.Itoa(len(ctl.QueueCS)), *p_id)
			display.Debug("[DEBUG] nombre de sites dans newglobalstate: "+ strconv.Itoa(len(ctl.GlobalState.States)), *p_id)
			display.Debug("[DEBUG] nombre de sites dans newglobalstatematrix: "+ strconv.Itoa(len(ctl.GlobalState.MatrixClock)), *p_id)*/

			for _, k := range sitelist {
				display.Debug("matrix clock : "+IntSliceToString(ctl.GlobalState.MatrixClock[k], ","), *p_id)
			}
			SendMessage(NewMessage(RETOUR_STATE, ctl.Id, msg.Sender, -1, -1, 0, ctl.VectorClock, nil, "", ctl.GlobalState.States[ctl.Id], nil, 0, -1))
			SendMessage(NewMessage(NEWSITE_APP, ctl.Id, ctl.Id, -1, -1, 0, ctl.VectorClock, sitelist, "", nil, nil, 0, -1))

		}
		if msg.Type == RETOUR_STATE {
			display.Success("Retour d'état reçu de "+strconv.Itoa(msg.Sender), *p_id)
			ctl.GlobalState.States[msg.Sender] = msg.State
		}
		if msg.Type == SUPPRESSION_CTL {
			//Traiter la suppression
			//             ctl.VectorClock = make([]int, len(msg.VectorClock))
			//             copy(ctl.VectorClock, msg.VectorClock)
			sitelist = make([]int, len(msg.SiteList))
			copy(sitelist, msg.SiteList)
			display.Info("Suppression d'un site"+strconv.Itoa(msg.Transmitter), *p_id)
			ctl.QueueCS[msg.Transmitter].Type = LIBERATION //l'id du site qui a été supprimé est stocké dans l'attribut Parent du message
			ctl.QueueCS[msg.Transmitter].Stamp = -1

			display.Info("Suppression d'un site", *p_id)
			SendMessage(NewMessage(DELETESITE_APP, ctl.Id, ctl.Id, -1, -1, 0, ctl.VectorClock, sitelist, "", nil, nil, 0, -1))
		}
		if msg.Type == DEMANDE_SNAPSHOT || // Le message est déjà traité
			msg.Type == ETAT || // par l'algorithme de construction d'instantanés
			msg.Type == PREPOST || msg.Type == SIGNAL_SNAPSHOT ||
			msg.Type == FIN_SNAPSHOT {
			continue
		}
		if msg.Type == DEBUT_SC ||
			msg.Type == ACTUALISATION { // Le message est de type inconnu
			display.Warning("Type de message inconnu : "+string(msg.Type), *p_id)
			continue
		}

		if msg.Type == ELECTION_BLEU || msg.Type == ELECTION_ROUGE {
			display.Info("Message destine au net", *p_id)
			continue
		}

		if msg.Type == DEMANDE_DEPART {
			display.Info("Reception d'une demande de depart du reseau par l'application", *p_id)
			display.Info("Transmission au net", *p_id)
			SendMessage(msg)

		}

		display.Success("Message reçu de "+strconv.Itoa(msg.Sender)+" : "+string(msg.Type), *p_id)

		ctl.fileAttenteRepartie(msg) // Algorithme de file d'attente répartie
	}
}

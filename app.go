package main

import (
	"flag"
	"fmt"
	"strconv"
	"sync"
	"time"

	"src/display"
	. "src/utils"
)

// Fonction pour gérer les messages utilisateur
func (app *Application) handleMessages() {
	id := app.Id
	userChan := app.UserChan
	messagesChan := app.MessagesChan
	mutex := app.Mutex
	p_siteActions := &app.State.SiteActions
	p_bankActions := &app.State.BankActions

	for {
		// Affichage du menu
		display.Display("[1] Consulter le nombre d'actions sur votre compte")
		display.Display("[2] Consulter le nombre d'actions disponibles à la banque")
		display.Display("[3] Acheter des actions")
		display.Display("[4] Vendre des actions")
		display.Display("[5] Quitter le réseau")
		if id == 0 {
			display.Display("[6] Sauvegarder l'état global")
			display.Display("[7] Afficher le dernier état global sauvegardé")
			display.Display("[8] Revenir au dernier état global sauvegardé")
		}
		// Conversion de la chaîne en entier
		choice, err := strconv.Atoi(<-userChan)
		if err == nil {
			//display.Info("Vous avez tapé le nombre "+strconv.Itoa(choice), id)
		} else {
			display.Error("Impossible de convertir la chaîne en entier", id)
		}

		mutex.Lock()

		switch choice {
		case 1:
			display.Display("Vous avez actuellement " +
				strconv.Itoa(*p_siteActions) + " actions \n")
		case 2:
			display.Display("Il y a actuellement " +
				strconv.Itoa(*p_bankActions) + " actions disponibles à la banque \n")
		case 3:
			SendMessage(NewMessage(DEMANDE_SC, id, id, -1, -1, 0, make([]int, NB_SITES), nil, "", nil, nil, -1, -1))
			mutex.Unlock()
			display.Info("NBSITES: "+strconv.Itoa(NB_SITES), id)
			display.Info("Attente de la section critique \n", id)
			for {
				// Attente de la section critique
				message := <-messagesChan
				if message.Type == DEBUT_SC {
					break
				}
			} // Fin de l'attente de la section critique
			display.Success("Section critique obtenue \n", id)

			display.Display("Il y a actuellement " +
				strconv.Itoa(*p_bankActions) + " actions disponibles à la banque")
			display.Display("Vous avez actuellement " +
				strconv.Itoa(*p_siteActions) + " actions sur votre compte")

			display.Display("Combien d'actions voulez-vous acheter ?")
			amountBuy, err := strconv.Atoi(<-userChan)
			for err != nil || amountBuy > *p_bankActions {
				if err != nil {
					display.Display("Veuillez entrer un nombre entier.")
				} else {
					display.Display("Il n'y a pas assez d'actions disponibles à la banque.")
				}
				amountBuy, err = strconv.Atoi(<-userChan)
			}

			mutex.Lock()

			*p_bankActions -= amountBuy
			*p_siteActions += amountBuy

			SendMessage(NewMessage(FIN_SC, id, id, -1, -1, 0, make([]int, NB_SITES), nil, "", NewState(*p_siteActions, *p_bankActions), nil, -1, -1))

		case 4:
			//display.Debug("Horloge vectorielle actuelle app: "+IntSliceToString(Vi[:], ","), i)
			SendMessage(NewMessage(DEMANDE_SC, id, id, -1, -1, 0, make([]int, NB_SITES), nil, "", nil, nil, -1, -1))
			mutex.Unlock()

			display.Info("Attente de la section critique \n", id)
			for {
				// Attente de la section critique
				message := <-messagesChan
				if message.Type == DEBUT_SC {
					break
				}
			} // Fin de l'attente de la section critique
			display.Success("Section critique obtenue \n", id)

			display.Display("Il y a actuellement " +
				strconv.Itoa(*p_bankActions) + " actions disponibles à la banque")
			display.Display("Vous avez actuellement " +
				strconv.Itoa(*p_siteActions) + " actions sur votre compte")

			display.Display("Combien d'actions voulez-vous vendre ?")
			amountSell, err := strconv.Atoi(<-userChan)
			for err != nil || amountSell > *p_siteActions {
				if err != nil {
					display.Display("Veuillez entrer un nombre entier.")
				} else {
					display.Display("Vous n'avez pas assez d'actions sur votre compte.")
				}
				amountSell, err = strconv.Atoi(<-userChan)
			}

			mutex.Lock()

			*p_bankActions += amountSell
			*p_siteActions -= amountSell

			//display.Info("Attente de la section critique \n" +strconv.Itoa(NB_SITES), id)
			SendMessage(NewMessage(FIN_SC, id, id, -1, -1, 0, make([]int, NB_SITES), nil, "", NewState(*p_siteActions, *p_bankActions), nil, -1, -1))

		case 5:
			display.Error("Vous avez choisi de quitter le réseau. Action executé dans 5 secondes ...", id)
			SendMessage(NewMessage(DEMANDE_DEPART, id, id, id, -1, 0, make([]int, NB_SITES), app.liste_sites, "", nil, nil, NB_SITES, -1))
			break
		case 6:
			if id == 0 {
				display.Success("Envoi d'une demande de sauvegarde de l'état global \n", id)
				SendMessage(NewMessage(DEMANDE_SNAPSHOT, id, id, -1, -1, 0, make([]int, NB_SITES), nil, BLANC, NewState(*p_siteActions, *p_bankActions), nil, 1, -1))
			}
		case 7:
			if id == 0 {
				for k := 0; k < NB_SITES; k++ {
					display.Display("Quantité d'action du site " + strconv.Itoa(k) + " : " + strconv.Itoa(app.GlobalState.States[k].SiteActions))
					display.Display("Quantité d'action de la banque " + strconv.Itoa(k) + " : " + strconv.Itoa(app.GlobalState.States[k].BankActions) + "\n")
				}
			}
		case 8:
			if id == 0 {
				SendMessage(NewMessage(DEMANDE_SC, id, id, -1, -1, 0, make([]int, NB_SITES), nil, "", nil, nil, -1, -1))
				mutex.Unlock()

				display.Info("Attente de la section critique \n", id)
				for {
					// Attente de la section critique
					message := <-messagesChan
					if message.Type == DEBUT_SC {
						break
					}
				} // Fin de l'attente de la section critique
				display.Success("Section critique obtenue \n", id)

				mutex.Lock()
				SendMessage(NewMessage(REVERT_INSTANTANE_APP, id, id, -1, -1, 0, make([]int, NB_SITES), nil, "", nil, app.GlobalState, -1, -1))

				SendMessage(NewMessage(FIN_SC, id, id, -1, -1, 0, make([]int, NB_SITES), nil, "", NewState(app.GlobalState.States[id].SiteActions, app.GlobalState.States[id].BankActions), nil, -1, -1))
			}
		default:
			display.Error("Choix inconnu \n", id)
		}

		// app.State.SiteActions = *siteActions
		// app.State.BankActions = *bankActions

		mutex.Unlock()
	}
}

// Fonction pour lire un message du contrôleur
func (app *Application) readControllerMessage(msg *Message) {
	id := app.Id
	messagesChan := app.MessagesChan
	siteActions := app.State.SiteActions
	bankActions := app.State.BankActions

	if msg.Receiver != id {
		//display.Error("Message non destiné à ce site", id)
		return
	}

	switch msg.Type {
	case DEBUT_SC:
		//display.Success("Début de la section critique", id)
		messagesChan <- msg
	case ACTUALISATION:
		//display.Success("Actualisation du nombre d'actions de la banque", id)
		bankActions = msg.State.BankActions
		display.Display("Le nombre d'actions disponibles à la banque a été actualisé à " +
			strconv.Itoa(bankActions))
	case ACTUALISATION_GLOBAL_STATE:
		//display.Success("Actualisation de l'etat global", id)
		app.GlobalState = msg.GlobalState

	case ACTUALISATION_INSTANTANE:
		//display.Success("Actualisation de l'etat global à son état antérieur", id)
		display.Display("Un retour en arrière a été demandé par un administrateur pour cause de problème interne")
		siteActions = msg.GlobalState.States[id].SiteActions
		display.Display("Votre compte d'actions a été retourné à : " + strconv.Itoa(siteActions) + " \n")
		bankActions = msg.GlobalState.States[id].BankActions

	case NEWSITE_APP:
		display.Success("App informe de l'ajout d'un nouveau site", id)
		app.liste_sites = make([]int, len(msg.SiteList))
		copy(app.liste_sites, msg.SiteList)
		NB_SITES += 1

	case DELETESITE_APP:
		display.Success("App informe de la suppression d'un site", id)
		//         app.liste_sites = make([]int, len(msg.SiteList))
		// 		copy(app.liste_sites, msg.SiteList)
		// 		NB_SITES = len(msg.SiteList)
	default:
		//display.Error("Type de message inconnu", id)
	}

	app.State.SiteActions = siteActions
	app.State.BankActions = bankActions
}

func (app *Application) readMessages() {
	var rawMsg string

	id := app.Id
	userChan := app.UserChan
	mutex := app.Mutex

	for {
		//display.Info("Attente d'un message...", id)
		_, err := fmt.Scanln(&rawMsg)

		mutex.Lock()

		if err != nil {
			//display.Error("Erreur lors de la lecture du message", id)
			mutex.Unlock()
			continue
		}

		msg, err := ParseMessage(rawMsg)
		if err != nil {
			display.Warning("Impossible de parser le message", id)
			//display.Info("Il s'agit peut-être d'un message utilisateur", id)
			userChan <- rawMsg
		} else {
			//display.Success("Message reçu du contrôleur", id)
			app.readControllerMessage(msg)
		}

		mutex.Unlock()
		rawMsg = ""
	}
}

type Application struct {
	Id          int    // Identifiant du site
	State       *State // État de l'application
	GlobalState *GlobalState
	liste_sites []int
	// Channels
	UserChan     chan string   // Channel pour les messages utilisateur
	MessagesChan chan *Message // Channel pour les messages du contrôleur

	// Mutex
	Mutex *sync.Mutex // Mutex pour les sections critiques
}

// Créer une nouvelle application
func NewApplication(id int, siteActions int, bankActions int, lsite []int) *Application {
	return &Application{
		Id:           id,
		State:        NewState(siteActions, bankActions),
		GlobalState:  NewGlobalState(),
		liste_sites:  lsite,
		UserChan:     make(chan string, 1<<16),
		MessagesChan: make(chan *Message, 1<<16),
		// IMPORTANT: Le buffer du channel doit être assez grand pour éviter un deadlock
		// lors de l'envoi de messages utilisateur
		// "Successfully detecting a deadlock would be akin to solving the halting problem."

		Mutex: &sync.Mutex{},
	}
}

func main() {
	p_id := flag.Int("id", -1, "Identifiant de l'application de base")
	flag.Parse()
	siteActions := 20
	bankActions := 100
	lsitenew := []int{0, 1, 2, 3, 4}
	if *p_id > 4 {
		for {
			msg, err := ReceiveMessage()
			if err != nil {
				display.Error("Erreur lors de la lecture du message", -1)
				continue
			}
			if msg.Type == NBSITESAPP {
				NB_SITES = msg.Balance
				lsitenew = make([]int, len(msg.SiteList))
				copy(lsitenew, msg.SiteList)
				display.Success("Nombre de sites : "+strconv.Itoa(NB_SITES), *p_id)
				break
			}
		}

	}
	if *p_id < 0 || *p_id > NB_SITES {
		display.Error("L'identifiant doit être compris entre 0 et "+
			strconv.Itoa(NB_SITES-1), -1)
		return
	}

	siteActions = 20
	bankActions = 100

	// Créer une nouvelle application
	app := NewApplication(*p_id, siteActions, bankActions, lsitenew)

	for k := 0; k < 5; k++ {
		app.GlobalState.States[k] = NewState(siteActions, bankActions)
	}

	defer close(app.MessagesChan)
	defer close(app.UserChan)

	go app.handleMessages()
	go app.readMessages()

	for {
		time.Sleep(time.Duration(60) * time.Second)
	} // Pour attendre la fin des goroutines...
	// select {} // Attendre indéfiniment
}

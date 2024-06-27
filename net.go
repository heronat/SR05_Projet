package main

import (
	"flag"
	"os"
	"os/exec"
	"src/display"

	//"strings"
	"bufio"
	. "src/utils"
	"strconv"
	"time"
	//"reflect"
)

var NB_ROUTEURS = 5
var liste_sites []int
var p_nb_neighbours = 1

type Net struct {
	Id                int
	Active            bool
	Vi                []int
	SitetoAdd         int
	Parent            int  // Parent du routeur
	TypeElection      bool // Type d'élection, true si ajout, false si suppression
	NbVoisinsAttendus int  // Nombre de voisins attendus
	Elu               int  // Site élu
}

func NewNet(id int, nbNeighbours int, vecteur []int) *Net {
	return &Net{
		Id:                id,
		Active:            true,
		Vi:                vecteur,
		SitetoAdd:         -1,
		Parent:            -1,
		TypeElection:      false,
		NbVoisinsAttendus: nbNeighbours,
		Elu:               int(^uint(0) >> 1),
	}
}

// Algorithme 22 : Élection par extinction de vagues
func (net *Net) election(msg *Message) {
	i := net.Id
	parent := net.Parent
	nbVoisinsAttendus := net.NbVoisinsAttendus
	elu := net.Elu

	j := msg.Sender
	k := msg.Elu

	switch msg.Type {
	case DEMANDE_ADMISSION, DEMANDE_DEPART:
		display.Success("Réception d'une demande d'admission/de départ de "+strconv.Itoa(msg.Sender), i)
		if parent == -1 {
			// Le site n'a pas encore été atteint par la vague ; il peut encore se déclarer candidat.
			elu = i
			parent = i
			net.SitetoAdd = j

			if msg.Type == DEMANDE_ADMISSION {
				net.TypeElection = true
			} else if msg.Type == DEMANDE_DEPART {
				net.TypeElection = false
				net.SitetoAdd = i
			}

			display.Debug("Parent : "+strconv.Itoa(parent)+"\tNbVoisinsAttendus : "+strconv.Itoa(nbVoisinsAttendus)+"\tElu : "+strconv.Itoa(elu), i)
			display.Info("Envoi d'une candidature du site "+strconv.Itoa(i), i)
			SendMessage(NewMessage(ELECTION_BLEU, i, -1, i, -1, -1, net.Vi, nil, "", nil, nil, -1, elu))
		} else {
			display.Warning("Le site a déjà été atteint par la vague, il ne peut pas se porter candidat", i)
		}
	case ELECTION_BLEU:
		display.Success("Réception d'une candidature de "+strconv.Itoa(j), i)
		if elu > k {
			// Première vague reçue, ou vague dont l’identité de l’élu est plus petite que la précédente.
			elu = k
			// On oublie la vague en cours (s’il y en avait une).
			parent = j
			nbVoisinsAttendus--
			if nbVoisinsAttendus > 0 {
				// On diffuse la nouvelle vague.
				display.Info("Diffusion de la vague du site "+strconv.Itoa(k), i)
				display.Debug("Parent : "+strconv.Itoa(parent)+"\tNbVoisinsAttendus : "+strconv.Itoa(nbVoisinsAttendus)+"\tElu : "+strconv.Itoa(elu), i)
				SendMessage(NewMessage(ELECTION_BLEU, i, -1, i, j, -1, net.Vi, nil, "", nil, nil, -1, elu))
			} else {
				// La vague remonte vers Sj.
				display.Info("Fin de la vague du site "+strconv.Itoa(k), i)
				SendMessage(NewMessage(ELECTION_ROUGE, i, j, i, -1, -1, net.Vi, nil, "", nil, nil, -1, elu))
			}
		} else {
			if elu == k {
				// Message appartenant à la même vague, mais Si est déjà au courant.
				// On renvoie un message rouge pour que la vague puisse remonter vers son élu.
				display.Info("Fin de la vague du site "+strconv.Itoa(k), i)
				SendMessage(NewMessage(ELECTION_ROUGE, i, j, i, -1, -1, net.Vi, nil, "", nil, nil, -1, elu))
			} else {
				// Message appartenant à une vague plus ancienne.
				// On ignore le message.
				display.Warning("Message appartenant à une vague plus ancienne", i)
			}
		}
	case ELECTION_ROUGE:
		if elu == k {
			// Seuls les messages de retour appartenant à la vague en cours sont acceptés.
			nbVoisinsAttendus--
			if nbVoisinsAttendus == 0 {
				if elu == i {
					// Fin de l'algorithme. Le site élu est Si.
					display.Success("Fin de l'algorithme. Le site élu est "+strconv.Itoa(elu), i)
					if net.TypeElection == true {
						display.Success("Acceptation du nouveau site dans le réseau ", i)
						//lancer add_site_accepted
						cmd := exec.Command("./scripts/v3/addsite_accepted.sh", strconv.Itoa(i), strconv.Itoa(net.SitetoAdd))
						if err := cmd.Run(); err != nil {
							display.Error("Erreur lors de l'appel du script addsite_accecepted.sh", i)
							return
						}
						time.Sleep(1 * time.Second)
						display.Error("Len Vi : "+strconv.Itoa(len(net.Vi)), i)
						p_nb_neighbours += 1
						net.Vi = append(net.Vi, 0)
						display.Debug("Envoi du vecteur clock vers ACCEPTATION_AJOUT: "+IntSliceToString(net.Vi, ","), i)
						SendMessage(NewMessage(ACCEPTATION_AJOUT, i, net.SitetoAdd, i, -1, 0, net.Vi, liste_sites, "", nil, nil, NB_SITES, -1))
					} else if net.TypeElection == false {
						display.Success("Suppression du site dans le réseau ", i)
						//lancer remove_site
						SendMessage(NewMessage(PREVENTION_VOISINS, i, -1, i, -1, 0, net.Vi, liste_sites, "", nil, nil, NB_SITES, -1))
						display.Error("Suppression acceptee transmise par  "+strconv.Itoa(msg.Sender), i)
						net.Active = false
						net.Vi = removeVectorClock(net.Vi, i)
						liste_sites = removeSite(liste_sites, i)
						for _, lsite := range liste_sites {

							if lsite != i {
								display.Success("Confirmation de la suppression a "+strconv.Itoa(lsite), i)
								SendMessage(NewMessage(SUPPRESSION, i, lsite, i, -1, -1, net.Vi, liste_sites, "", nil, nil, -1, 0))
							}
						}
						time.Sleep(1 * time.Second)
						cmd := exec.Command("./scripts/v3/removesite.sh", strconv.Itoa(net.SitetoAdd))
						if err := cmd.Run(); err != nil {
							display.Error("Erreur lors de l'appel du script removesite.sh", net.SitetoAdd)
							return
						}
					}

				} else {
					// La vague remonte vers le parent.
					display.Info("La vague remonte vers le parent", i)
					SendMessage(NewMessage(ELECTION_ROUGE, i, parent, i, -1, -1, net.Vi, nil, "", nil, nil, -1, elu))
				}
			}
		}
	}

	net.Parent = parent
	net.NbVoisinsAttendus = nbVoisinsAttendus
	net.Elu = elu
}

// a mettre dans utils
func removeVectorClock(slice []int, s int) []int {
	if s >= 0 && s < len(slice) {
		slice[s] = -1
	}
	return slice
}

func removeSite(slice []int, n int) []int {
	result := []int{}
	for _, value := range slice {
		if value != n {
			result = append(result, value)
		}
	}
	return result
}

func getFileLineNumber(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		display.Error("Erreur lors de l'ouverture du fichier : "+err.Error(), -1)
		return -1
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		display.Error("Erreur lors de la lecture du fichier : "+err.Error(), -1)
		return -1
	}

	return lineCount
}

func main() {
	p_id := flag.Int("id", -1, "Identifiant du routeur")
	flag.Parse()
    //initialisation des net de base
	if *p_id >= 0 && *p_id < 5 {
		p_nb_neighbours = getFileLineNumber("/tmp/map_N"+strconv.Itoa(*p_id)) - 1
	}
	display.Info("Nombre de voisins"+strconv.Itoa(p_nb_neighbours), *p_id)

	Vi := make([]int, NB_SITES)
	liste_sites = append(liste_sites, 0, 1, 2, 3, 4)
	if *p_id > 4 {
        //Ajout d'un nouveau site dans le réseau
		display.Info("Demande en cours du nombre de sites au parent", *p_id)
		SendMessage(NewMessage(DEMANDE_ADMISSION, *p_id, -1, *p_id, -1, 0, nil, nil, "", nil, nil, -1, -1))

	} else {
		display.Info("Nombre Site ", NB_SITES)

		if *p_id < 0 || *p_id >= NB_SITES {
			display.Error("L'identifiant doit être compris entre 0 et "+
				strconv.Itoa(NB_SITES-1), -1)
			return
		}

	}
	net := NewNet(*p_id, p_nb_neighbours, Vi)

	if *p_id < 5 {
		display.Debug("On partage la liste de sites au ctlr", *p_id)
		SendMessage(NewMessage(LISTE_CTL, *p_id, *p_id, -1, -1, 0, net.Vi, liste_sites, "", nil, nil, -1, -1))
	}

	// Boucle de traitement des messages reçus
	for {

		display.Info("Entrez un message:", *p_id)
		if net.Active == false {
			display.Error("Le site est en veille", *p_id) // On ne traite pas les messages si le site est en veille, on les renvoit juste sur la sortie
		}

		msg, err := ReceiveMessage()

		if err != nil {
			display.Error("Le format du message reçu est incorrect : "+string(err.Error()), *p_id)
			//fmt.Println(msg)
			//time.Sleep(5*time.Second)
			continue
		}

        //Afficher la VectorClock du message reçu pour vérifier la cohérence
		display.Debug("Vi reçu: "+IntSliceToString(msg.VectorClock, ","), *p_id)

		if msg.Parent == *p_id {
			display.Info("Type message rebond: "+string(msg.Type), *p_id)
			display.Info("Message de rebond de "+strconv.Itoa(msg.Transmitter), *p_id)
			continue
		}

		if msg.Sender == msg.Receiver && msg.Sender != *p_id {
			display.Info("Message entre un net annexe et son ctl, qui est à drop", *p_id)
			continue
		}

		if msg.Receiver != *p_id && msg.Receiver != -1 {
			display.Info("Message non destiné à ce site : "+strconv.Itoa(msg.Receiver), *p_id)
			display.Info("Transmission du message "+string(msg.Type)+" à "+strconv.Itoa(msg.Receiver), *p_id)
			display.Info("Je recois : ("+strconv.Itoa(msg.Sender)+", "+strconv.Itoa(msg.Receiver)+", "+strconv.Itoa(msg.Transmitter)+", "+strconv.Itoa(msg.Parent)+")", *p_id)
			msg.Parent = msg.Transmitter
			msg.Transmitter = *p_id
			display.Info("Je transmets : ("+strconv.Itoa(msg.Sender)+", "+strconv.Itoa(msg.Receiver)+", "+strconv.Itoa(msg.Transmitter)+", "+strconv.Itoa(msg.Parent)+")", *p_id)
			SendMessage(msg) // On renvoie le message dans le réseau car il ne nous est pas destiné
			continue
		}

		if msg.Type == AJOUT {
			if net.Active {
				display.Success("Message d'ajout de l'utilisateur: "+strconv.Itoa(msg.Sender), *p_id)

                //Un site a été ajouté dans le réseau, on met à jour notre liste de sites et la taille de notre horloge vectorielle
				net.Vi = make([]int, len(msg.VectorClock))
				copy(net.Vi, msg.VectorClock)
				liste_sites = make([]int, len(msg.SiteList))
				copy(liste_sites, msg.SiteList)

				display.Success("Nouveau site ajoute a la liste :"+strconv.Itoa(msg.Sender), *p_id)
				NB_SITES += 1
				NB_ROUTEURS += 1

				//On réinitialise les valeurs de l'election
				net.Parent = -1
				net.NbVoisinsAttendus = p_nb_neighbours
				net.Elu = int(^uint(0) >> 1)
				display.Info("RESET des valeurs de l election "+strconv.Itoa(msg.Sender), *p_id)

				display.Success("Nombre sites du pdv de "+strconv.Itoa(*p_id)+" : "+strconv.Itoa(NB_SITES), *p_id)


			}
		} else if msg.Type == SUPPRESSION {
			if net.Active {
				display.Error("Message de suppression de l'utilisateur: "+strconv.Itoa(msg.Sender), *p_id)

                //Un site a été supprimé du réseau, on met à jour notre liste de sites et la taille de notre horloge vectorielle
				net.Vi = make([]int, len(msg.VectorClock))
				net.NbVoisinsAttendus = p_nb_neighbours
				copy(net.Vi, msg.VectorClock)
				liste_sites = make([]int, len(msg.SiteList))
				copy(liste_sites, msg.SiteList)

                //On réinitialise les valeurs de l'election
				net.Parent = -1
				net.NbVoisinsAttendus = p_nb_neighbours
				net.Elu = int(^uint(0) >> 1)
				display.Info("RESET des valeurs de l election "+strconv.Itoa(msg.Sender), *p_id)

                //On prévient notre controleur de la suppression d'un site
				SendMessage(NewMessage(SUPPRESSION_CTL, *p_id, *p_id, msg.Sender, -1, 0, net.Vi, liste_sites, "", nil, nil, NB_SITES, -1))
			}
		} else if msg.Type == ACCEPTATION_AJOUT {
            //Message reçu chez le nouveau site pour le prévenir que sa demande d'ajout a été acceptée
			NB_SITES = msg.Balance + 1
			NB_ROUTEURS = NB_ROUTEURS + 1
			display.Debug("Attente de 3 secondes ...", *p_id)
			time.Sleep(1 * time.Second)
			display.Success("Acceptation du nouveau site dans le réseau", *p_id)
			display.Success("Nombre Site NEW :  "+strconv.Itoa(NB_SITES), *p_id)

			net.Vi = make([]int, len(msg.VectorClock))
			copy(net.Vi, msg.VectorClock)

			net.Parent = -1
			net.NbVoisinsAttendus = p_nb_neighbours
			net.Elu = int(^uint(0) >> 1)
			display.Info("RESET des valeurs de l election pour l'acceptation ajout"+strconv.Itoa(msg.Sender), *p_id)

			if *p_id < 0 || *p_id > NB_SITES {
				display.Error("L'identifiant doit être compris entre 0 et "+
					strconv.Itoa(NB_SITES-1), -1)
				return
			}


			liste_sites = make([]int, len(msg.SiteList))

			copy(liste_sites, msg.SiteList)          // Un nouveau site doit s'init avec lal iste des sites courants, sinon pb
			liste_sites = append(liste_sites, *p_id) // On ajoute le nv site à la liste des sites

			SendMessage(NewMessage(LISTE_CTL, *p_id, *p_id, -1, -1, 0, net.Vi, liste_sites, "", nil, nil, NB_SITES, -1))

			for _, lsite := range liste_sites {
                //On prévient tous les autres sites qu'un nouveau site a rejoint le réseau
				if lsite != *p_id {
					display.Success("Envoi de l'ajout à "+strconv.Itoa(lsite), *p_id)
					SendMessage(NewMessage(AJOUT, *p_id, lsite, *p_id, -1, 0, net.Vi, liste_sites, "", nil, nil, -1, -1)) // On envoie l'ajout à tous les sites, on partage la nouvelle liste
				}
			}

		} else if msg.Type == PREVENTION_VOISINS {
            //Un de nos voisins a été supprimé, on ajuste notre nombre de voisins local
			p_nb_neighbours -= 1
			display.Info("Suppression d'un voisin", *p_id)
			display.Info("Nouveau nombre de voisins"+strconv.Itoa(p_nb_neighbours), *p_id)
		} else if msg.Type == REQUETE || msg.Type == DEMANDE_SC || msg.Type == FIN_SC || msg.Type == LIBERATION || msg.Type == ACCUSE || msg.Type == DEBUT_SC || msg.Type == DEMANDE_STATE || msg.Type == RETOUR_STATE {
            //Ce message vient du controleur
			display.Info("Message entre controleur", *p_id)
			msg.Parent = msg.Transmitter
			msg.Transmitter = *p_id
			SendMessage(msg)
		} else {
			if net.Active {
				display.Info("Message reçu de l'utilisateur: "+strconv.Itoa(msg.Sender), *p_id)
				net.election(msg)
			}
		}

	}

}

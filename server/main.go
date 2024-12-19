package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Struct représentant une excuse avec sa catégorie
type Excuse struct {
	Category string `json:"category"`
	Excuse   string `json:"excuse"`
}

var (
	excuses     []Excuse
	excusesLock sync.RWMutex
	excusesPath string // Variable globale pour le chemin du fichier excuses.json
)

// Fonction pour obtenir le répertoire du fichier source
func getSourceDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Impossible de déterminer le répertoire source")
	}
	return filepath.Dir(filename)
}

// Fonction pour charger les excuses depuis le fichier JSON
func loadExcuses(filename string) error {
	log.Printf("Tentative d'ouverture de %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Erreur lors de l'ouverture de %s : %v", filename, err)
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var loadedExcuses []Excuse
	if err := decoder.Decode(&loadedExcuses); err != nil {
		log.Printf("Erreur lors du décodage de %s : %v", filename, err)
		return err
	}

	excusesLock.Lock()
	excuses = loadedExcuses
	excusesLock.Unlock()

	log.Printf("Chargement des excuses réussi, total : %d", len(excuses))
	return nil
}

// Fonction pour sauvegarder les excuses dans le fichier JSON
func saveExcuses(filename string) error {
	excusesLock.RLock()
	defer excusesLock.RUnlock()

	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Erreur lors de la création de %s : %v", filename, err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(excuses); err != nil {
		log.Printf("Erreur lors de l'encodage de %s : %v", filename, err)
		return err
	}

	log.Printf("Sauvegarde des excuses dans %s réussie", filename)
	return nil
}

// Handler pour récupérer une excuse aléatoire, éventuellement filtrée par catégorie
func getRandomExcuse(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	excusesLock.RLock()
	defer excusesLock.RUnlock()

	var filtered []Excuse
	if category != "" {
		for _, e := range excuses {
			if strings.EqualFold(e.Category, category) {
				filtered = append(filtered, e)
			}
		}
	} else {
		filtered = excuses
	}

	if len(filtered) == 0 {
		http.Error(w, "Aucune excuse disponible pour cette catégorie", http.StatusNotFound)
		return
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(filtered))
	excuse := filtered[randomIndex]

	response := map[string]string{
		"category": excuse.Category,
		"excuse":   excuse.Excuse,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handler pour ajouter une nouvelle excuse
func addExcuse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var newExcuse Excuse
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newExcuse); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	if newExcuse.Category == "" || newExcuse.Excuse == "" {
		http.Error(w, "Catégorie et excuse sont requises", http.StatusBadRequest)
		return
	}

	excusesLock.Lock()
	excuses = append(excuses, newExcuse)
	excusesLock.Unlock()

	// Sauvegarder les excuses mises à jour dans le fichier JSON
	err := saveExcuses(excusesPath)
	if err != nil {
		http.Error(w, "Erreur lors de la sauvegarde", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newExcuse)
}

// Handler pour récupérer toutes les catégories uniques
func getCategories(w http.ResponseWriter, r *http.Request) {
	excusesLock.RLock()
	defer excusesLock.RUnlock()

	categoryMap := make(map[string]bool)
	for _, e := range excuses {
		categoryMap[e.Category] = true
	}

	var categories []string
	for category := range categoryMap {
		categories = append(categories, category)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func main() {
	// Définir un drapeau pour le chemin du fichier excuses.json
	var excusesFile string
	flag.StringVar(&excusesFile, "excuses", "excuses.json", "Chemin vers le fichier excuses.json")
	flag.Parse()

	// Obtenir le répertoire du fichier source
	sourceDir := getSourceDir()

	// Construire le chemin complet vers excuses.json
	excusesPath = filepath.Join(sourceDir, excusesFile)
	log.Printf("Chemin vers excuses.json : %s", excusesPath)

	// Charger les excuses au démarrage
	log.Println("Début du chargement des excuses.")
	err := loadExcuses(excusesPath)
	if err != nil {
		log.Fatalf("Erreur lors du chargement des excuses : %v", err)
	}
	log.Println("Fin du chargement des excuses.")

	// Servir les fichiers statiques depuis le répertoire 'static'
	staticDir := filepath.Join(sourceDir, "static")
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fs)

	// Endpoint API pour obtenir une excuse
	http.HandleFunc("/api/excuse", getRandomExcuse)

	// Endpoint API pour ajouter une nouvelle excuse
	http.HandleFunc("/api/excuse/add", addExcuse)

	// Endpoint API pour obtenir toutes les catégories
	http.HandleFunc("/api/categories", getCategories)

	port := ":8080"
	log.Printf("Serveur démarré sur le port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Erreur du serveur : %v", err)
	}
}

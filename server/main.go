package main

import (
    "encoding/json"
    "flag"
    "log"
    "math/rand"
    "net/http"
    "os"
    "path/filepath"
    "time"
)

var excuses []string

// Charge les excuses depuis le fichier JSON
func loadExcuses(filename string) error {
    log.Printf("Tentative d'ouverture de %s", filename)
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&excuses); err != nil {
        return err
    }
    log.Printf("Chargement des excuses réussi, total : %d", len(excuses))
    return nil
}

// Handler pour récupérer une excuse aléatoire
func getRandomExcuse(w http.ResponseWriter, r *http.Request) {
    if len(excuses) == 0 {
        http.Error(w, "Aucune excuse disponible", http.StatusInternalServerError)
        return
    }

    rand.Seed(time.Now().UnixNano())
    randomIndex := rand.Intn(len(excuses))
    excuse := excuses[randomIndex]

    response := map[string]string{"excuse": excuse}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    // Définir un drapeau pour le chemin du fichier excuses.json
    var excusesFile string
    flag.StringVar(&excusesFile, "excuses", "excuses.json", "Chemin vers le fichier excuses.json")
    flag.Parse()

    // Obtenir le répertoire de travail actuel
    cwd, err := os.Getwd()
    if err != nil {
        log.Fatalf("Impossible d'obtenir le répertoire de travail : %v", err)
    }

    // Construire le chemin complet vers excuses.json
    excusesPath := filepath.Join(cwd, excusesFile)
    log.Printf("Chemin vers excuses.json : %s", excusesPath)

    // Charger les excuses au démarrage
    err = loadExcuses(excusesPath)
    if err != nil {
        log.Fatalf("Erreur lors du chargement des excuses : %v", err)
    }

    http.HandleFunc("/excuse", getRandomExcuse)

    port := ":8080"
    log.Printf("Serveur démarré sur le port %s", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatalf("Erreur du serveur : %v", err)
    }
}

// client/main.go
package main

import (
    "bufio"
    "bytes"
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
)

// Struct pour la réponse de l'excuse
type ExcuseResponse struct {
    Category string `json:"category"`
    Excuse   string `json:"excuse"`
}

func getCategories(baseURL string) ([]string, error) {
    url := fmt.Sprintf("%s/api/categories", baseURL)

    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("Erreur lors de la requête : %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("Erreur du serveur (%d) : %s", resp.StatusCode, string(body))
    }

    var categories []string
    if err := json.NewDecoder(resp.Body).Decode(&categories); err != nil {
        return nil, fmt.Errorf("Erreur lors du décodage de la réponse : %v", err)
    }

    return categories, nil
}

// Fonction pour générer une excuse
func getExcuse(baseURL string, category string) error {
    // Construire l'URL avec le paramètre de catégorie si fourni
    url := fmt.Sprintf("%s/api/excuse", baseURL)
    if category != "" {
        // Remplacer les espaces par %20 pour les URL
        url += fmt.Sprintf("?category=%s", strings.ReplaceAll(category, " ", "%20"))
    }

    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("Erreur lors de la requête : %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("Erreur du serveur (%d) : %s", resp.StatusCode, string(body))
    }

    var excuseResp ExcuseResponse
    if err := json.NewDecoder(resp.Body).Decode(&excuseResp); err != nil {
        return fmt.Errorf("Erreur lors du décodage de la réponse : %v", err)
    }

    fmt.Printf("Catégorie : %s\nExcuse : %s\n", excuseResp.Category, excuseResp.Excuse)
    return nil
}

// Fonction pour ajouter une nouvelle excuse
func addExcuse(baseURL string) error {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Entrez la catégorie de l'excuse : ")
    category, err := reader.ReadString('\n')
    if err != nil {
        return fmt.Errorf("Erreur de lecture : %v", err)
    }
    category = strings.TrimSpace(category)

    fmt.Print("Entrez le texte de l'excuse : ")
    excuse, err := reader.ReadString('\n')
    if err != nil {
        return fmt.Errorf("Erreur de lecture : %v", err)
    }
    excuse = strings.TrimSpace(excuse)

    if category == "" || excuse == "" {
        return fmt.Errorf("Catégorie et excuse ne peuvent pas être vides")
    }

    newExcuse := ExcuseResponse{
        Category: category,
        Excuse:   excuse,
    }

    data, err := json.Marshal(newExcuse)
    if err != nil {
        return fmt.Errorf("Erreur lors du marshalling : %v", err)
    }

    url := fmt.Sprintf("%s/api/excuse/add", baseURL)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
    if err != nil {
        return fmt.Errorf("Erreur lors de la requête : %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("Erreur du serveur (%d) : %s", resp.StatusCode, string(body))
    }

    fmt.Println("Excuse ajoutée avec succès !")
    return nil
}

func main() {
    // Définir les flags pour le mode et l'URL de base
    mode := flag.String("mode", "generate", "Mode d'opération : generate ou add")
    baseURL := flag.String("url", "http://localhost:8080", "URL de base du serveur")
    flag.Parse()

    switch *mode {
    case "generate":
        // Récupérer les catégories disponibles
        categories, err := getCategories(*baseURL)
        if err != nil {
            log.Printf("Erreur : %v\n", err)
            return
        }

        // Ajouter une option "Toutes les catégories"
        categories = append([]string{"Toutes les catégories"}, categories...)

        // Afficher les catégories
        fmt.Println("Sélectionnez une catégorie :")
        for i, category := range categories {
            fmt.Printf("%d. %s\n", i+1, category)
        }

        // Lire le choix de l'utilisateur
        reader := bufio.NewReader(os.Stdin)
        fmt.Printf("Entrez le numéro de la catégorie (1-%d) : ", len(categories))
        choiceStr, err := reader.ReadString('\n')
        if err != nil {
            log.Printf("Erreur de lecture : %v\n", err)
            return
        }
        choiceStr = strings.TrimSpace(choiceStr)

        // Convertir le choix en entier
        var choice int
        _, err = fmt.Sscanf(choiceStr, "%d", &choice)
        if err != nil || choice < 1 || choice > len(categories) {
            log.Printf("Choix invalide.")
            return
        }

        // Déterminer la catégorie sélectionnée
        selectedCategory := ""
        if choice != 1 { // 1 correspond à "Toutes les catégories"
            selectedCategory = categories[choice-1]
        }

        // Générer l'excuse avec la catégorie sélectionnée
        err = getExcuse(*baseURL, selectedCategory)
        if err != nil {
            log.Printf("Erreur : %v\n", err)
        }
    case "add":
        err := addExcuse(*baseURL)
        if err != nil {
            log.Printf("Erreur : %v\n", err)
        }
    default:
        fmt.Println("Mode inconnu. Utilisez 'generate' ou 'add'.")
    }
}

package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
)

type ExcuseResponse struct {
    Excuse string `json:"excuse"`
}

func getExcuse() error {
    url := "http://localhost:8080/excuse"

    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("Erreur lors de la requête : %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("Erreur du serveur : %s", body)
    }

    var excuseResp ExcuseResponse
    if err := json.NewDecoder(resp.Body).Decode(&excuseResp); err != nil {
        return fmt.Errorf("Erreur lors du décodage de la réponse : %v", err)
    }

    fmt.Printf("Excuse obtenue : %s\n", excuseResp.Excuse)
    return nil
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Println("=== Client d'Excuses ===")
    fmt.Println("Commands:")
    fmt.Println("  generate - Obtenir une excuse aléatoire")
    fmt.Println("  exit     - Quitter le client")
    fmt.Println("========================")

    for {
        fmt.Print("> ")
        if !scanner.Scan() {
            break
        }
        input := strings.TrimSpace(scanner.Text())

        switch strings.ToLower(input) {
        case "generate":
            err := getExcuse()
            if err != nil {
                log.Printf("Erreur : %v\n", err)
            }
        case "exit", "quit":
            fmt.Println("Au revoir !")
            return
        case "":
            // Ignorer les entrées vides
            continue
        default:
            fmt.Println("Commande inconnue. Utilisez 'generate' ou 'exit'.")
        }
    }

    if err := scanner.Err(); err != nil {
        log.Printf("Erreur de lecture : %v\n", err)
    }
}

package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/<votre-utilisateur-github>/go_course_manager/server/handler"
)

func main() {
    // Initialisation de la base de données
    handler.InitDB()

    // Définition des routes
    http.HandleFunc("/", handler.HomeHandler)
    http.HandleFunc("/ajouter-cours", handler.AddCourseHandler)
    http.HandleFunc("/telecharger-cours", handler.DownloadCourseHandler)

    fmt.Println("Le serveur est lancé sur http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

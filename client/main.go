package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
)

func main() {
    fmt.Println("Téléchargement des cours en PDF...")
    resp, err := http.Get("http://localhost:8080/telecharger-cours")
    if err != nil {
        fmt.Println("Erreur :", err)
        return
    }
    defer resp.Body.Close()

    outFile, err := os.Create("courses.pdf")
    if err != nil {
        fmt.Println("Erreur :", err)
        return
    }
    defer outFile.Close()

    _, err = io.Copy(outFile, resp.Body)
    if err != nil {
        fmt.Println("Erreur :", err)
        return
    }
    fmt.Println("Téléchargement terminé : courses.pdf")
}

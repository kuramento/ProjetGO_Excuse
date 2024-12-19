package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

type ExcuseResponse struct {
    Excuse string `json:"excuse"`
}

func main() {
    url := "http://localhost:8080/excuse"

    resp, err := http.Get(url)
    if err != nil {
        log.Fatalf("Erreur lors de la requête : %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        log.Fatalf("Erreur du serveur : %s", body)
    }

    var excuseResp ExcuseResponse
    if err := json.NewDecoder(resp.Body).Decode(&excuseResp); err != nil {
        log.Fatalf("Erreur lors du décodage de la réponse : %v", err)
    }

    fmt.Printf("Excuse obtenue : %s\n", excuseResp.Excuse)
}

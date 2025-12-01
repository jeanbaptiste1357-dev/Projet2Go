package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
)

// Struct qui représente l'API principale
type API struct {
    Artists   string `json:"artists"`
    Locations string `json:"locations"`
    Dates     string `json:"dates"`
    Relation  string `json:"relation"`
}

type Artist struct {
    ID         int      `json:"id"`
    Name       string   `json:"name"`
    Members    []string `json:"members"`
    Image      string   `json:"image"`
    FirstAlbum string   `json:"firstAlbum"`
    CreationDate int    `json:"creationDate"`
}


func main() {
    // 1) On récupère l'API principale
    resp, err := http.Get("https://groupietrackers.herokuapp.com/api")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    // 2) On lit les données
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    // 3) On transforme le JSON en struct Go
    var api API
    err = json.Unmarshal(body, &api)
    if err != nil {
        log.Fatal(err)
    }

    // 4) On vérifie que ça marche
    fmt.Println("URL des artistes :", api.Artists)

	// 5) On récupère les artistes
respArtists, err := http.Get(api.Artists)
if err != nil {
    log.Fatal(err)
}
defer respArtists.Body.Close()

bodyArtists, err := io.ReadAll(respArtists.Body)
if err != nil {
    log.Fatal(err)
}

var artists []Artist
err = json.Unmarshal(bodyArtists, &artists)
if err != nil {
    log.Fatal(err)
}

// 6) Afficher le premier artiste pour vérifier
fmt.Println("Premier artiste :", artists[0].Name)

}

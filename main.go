package main

import (
    "encoding/json"
    "fmt"
    "html/template"
    "io"
    "log"
    "net/http"
)

type API struct {
    Artists   string `json:"artists"`
    Locations string `json:"locations"`
    Dates     string `json:"dates"`
    Relation  string `json:"relation"`
}

type Artist struct {
    ID           int      `json:"id"`
    Name         string   `json:"name"`
    Members      []string `json:"members"`
    Image        string   `json:"image"`
    FirstAlbum   string   `json:"firstAlbum"`
    CreationDate int      `json:"creationDate"`
}

func main() {
    // Quand l'utilisateur va sur http://localhost:8080/artists
    http.HandleFunc("/artists", artistsHandler)

    fmt.Println("Serveur lancé sur : http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func artistsHandler(w http.ResponseWriter, r *http.Request) {

    // 1) On récupère les URLs de l'API
    resp, err := http.Get("https://groupietrackers.herokuapp.com/api")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    var api API
    json.Unmarshal(body, &api)

    // 2) On récupère les artistes
    respArtists, err := http.Get(api.Artists)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer respArtists.Body.Close()

    bodyArtists, _ := io.ReadAll(respArtists.Body)

    var artists []Artist
    json.Unmarshal(bodyArtists, &artists)

    // 3) On charge le template
    tmpl, err := template.ParseFiles("templates/artists.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 4) On envoie les artistes dans le HTML
    tmpl.Execute(w, artists)

    http.HandleFunc("/artist", artistHandler)

}

func artistHandler(w http.ResponseWriter, r *http.Request) {
    // 1) Récupérer l’ID dans l’URL
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "Missing id", http.StatusBadRequest)
        return
    }


    // 2) Lire l’API principale
    resp, err := http.Get("https://groupietrackers.herokuapp.com/api")
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    defer resp.Body.Close()
    body, _ := io.ReadAll(resp.Body)

    var api API
    json.Unmarshal(body, &api)

    // 3) Récupérer la liste des artistes
    respArtists, err := http.Get(api.Artists)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    defer respArtists.Body.Close()
    bodyArtists, _ := io.ReadAll(respArtists.Body)

    var artists []Artist
    json.Unmarshal(bodyArtists, &artists)

    // 4) Trouver celui qui correspond à l’ID
    var selected Artist
    found := false

    for _, a := range artists {
        if fmt.Sprint(a.ID) == id {
            selected = a
            found = true
            break
        }
    }

    if !found {
        http.Error(w, "Artist not found", 404)
        return
    }

    // 5) Charger le template
    tmpl, err := template.ParseFiles("templates/artist.html")
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    // 6) Envoyer les données au template
    tmpl.Execute(w, selected)

}

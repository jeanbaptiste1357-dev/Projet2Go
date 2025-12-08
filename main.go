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

type Relation struct {
    Index          int                 `json:"index"`
    DatesLocations map[string][]string `json:"datesLocations"`
}

type ArtistPage struct {
    Artist   Artist
    Concerts Relation
}


func main() {
    // Quand l'utilisateur va sur http://localhost:8080/artists
    http.HandleFunc("/artists", artistsHandler)

    fmt.Println("Serveur lancé sur : http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
    // 1) Récupérer l'id dans l'URL
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "ID manquant", http.StatusBadRequest)
        return
    }

    // 2) Appeler l'API principale pour obtenir les URLs
    respAPI, err := http.Get("https://groupietrackers.herokuapp.com/api")
    if err != nil {
        http.Error(w, "Erreur API: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer respAPI.Body.Close()
    bodyAPI, err := io.ReadAll(respAPI.Body)
    if err != nil {
        http.Error(w, "Lecture API: "+err.Error(), http.StatusInternalServerError)
        return
    }

    var api API
    if err := json.Unmarshal(bodyAPI, &api); err != nil {
        http.Error(w, "Parse API: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // 3) Récupérer la liste des artistes
    respArtists, err := http.Get(api.Artists)
    if err != nil {
        http.Error(w, "Erreur artistes: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer respArtists.Body.Close()
    bodyArtists, err := io.ReadAll(respArtists.Body)
    if err != nil {
        http.Error(w, "Lecture artistes: "+err.Error(), http.StatusInternalServerError)
        return
    }

    var artists []Artist
    if err := json.Unmarshal(bodyArtists, &artists); err != nil {
        http.Error(w, "Parse artistes: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // 4) Trouver l'artiste demandé
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
        http.Error(w, "Artiste introuvable", http.StatusNotFound)
        return
    }

    // 5) Récupérer les relations (dates + lieux) — NOTE: on déclare 'relations' ici
    respRelation, err := http.Get(api.Relation)
    if err != nil {
        http.Error(w, "Erreur relations: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer respRelation.Body.Close()
    bodyRelation, err := io.ReadAll(respRelation.Body)
    if err != nil {
        http.Error(w, "Lecture relations: "+err.Error(), http.StatusInternalServerError)
        return
    }

    var relations []Relation
    if err := json.Unmarshal(bodyRelation, &relations); err != nil {
        http.Error(w, "Parse relations: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // 6) Trouver la relation correspondant à l'artiste sélectionné
    var rel Relation
    relFound := false
    for _, r := range relations {
        if r.Index == selected.ID {
            rel = r
            relFound = true
            break
        }
    }
    if !relFound {
        // Si aucune relation trouvée, on peut laisser rel vide (map nil) ou initialiser vide
        rel = Relation{
            Index:          selected.ID,
            DatesLocations: map[string][]string{},
        }
    }

    // 7) Préparer les données pour le template et l'exécuter
    data := ArtistPage{
        Artist:   selected,
        Concerts: rel,
    }

    tmpl, err := template.ParseFiles("templates/artist.html")
    if err != nil {
        http.Error(w, "Template: "+err.Error(), http.StatusInternalServerError)
        return
    }

    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Render template: "+err.Error(), http.StatusInternalServerError)
        return
    }
}

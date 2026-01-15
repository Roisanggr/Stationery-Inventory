package routes

import (
    "encoding/json"
    "net/http"

    "ATK-Backend/models"
    "github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
    r.HandleFunc("/api/atk", GetATK).Methods("GET")
    r.HandleFunc("/api/atk", CreateATK).Methods("POST")
    r.HandleFunc("/api/atk/{id}", UpdateATK).Methods("PUT")
    r.HandleFunc("/api/atk/{id}", DeleteATK).Methods("DELETE")
}

func GetATK(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    items, err := models.Get()
    if err != nil {
        http.Error(w, "failed to get data", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(items)
}

func CreateATK(w http.ResponseWriter, r *http.Request) {
    var newItem models.ATK
    if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }
    if newItem.Qty < 0 {
        http.Error(w, "Qty tidak boleh negatif", http.StatusBadRequest)
        return
    }
    created, err := models.Post(newItem)
    if err != nil {
        http.Error(w, "failed to create item", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(created)
}

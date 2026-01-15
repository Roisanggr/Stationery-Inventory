package routes

import (
    "encoding/json"
    "net/http"
    "strconv"

    "ATK-Backend/models"
    "github.com/gorilla/mux"
)

func UpdateATK(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }

    var updated models.ATK
    if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }
    if updated.Qty < 0 {
        http.Error(w, "Qty tidak boleh negatif", http.StatusBadRequest)
        return
    }

    res, ok, err := models.Update(id, updated)
    if err != nil {
        http.Error(w, "failed to update", http.StatusInternalServerError)
        return
    }
    if !ok {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(res)
}

func DeleteATK(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }

    ok, err := models.Delete(id)
    if err != nil {
        http.Error(w, "failed to delete", http.StatusInternalServerError)
        return
    }
    if !ok {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

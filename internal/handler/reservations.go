package handler

import (
    "encoding/json"
    "net/http"

)

// CreateReservationHandler обрабатывает запрос на создание нового бронирования
func CreateReservationHandler(db *storage.PostgresDB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var res model.Reservation
        if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        // Проверка на пересечение бронирований
        conflict, err := db.IsConflict(res.RoomID, res.StartTime, res.EndTime)
        if err != nil {
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        }
        if conflict {
            http.Error(w, "Reservation conflict", http.StatusConflict)
            return
        }

        // Создание бронирования
        if err := db.CreateReservation(&res); err != nil {
            http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
    }
}

// GetReservationsHandler обрабатывает запрос на получение всех бронирований для зала
func GetReservationsHandler(db *storage.PostgresDB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        roomID := chi.URLParam(r, "room_id")

        reservations, err := db.GetReservations(roomID)
        if err != nil {
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(reservations)
    }
}

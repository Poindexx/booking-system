package model

import "time"

// Reservation представляет бронирование конференц-зала
type Reservation struct {
    ID        int       `json:"id"`
    RoomID    string    `json:"room_id"`
    StartTime time.Time `json:"start_time"`
    EndTime   time.Time `json:"end_time"`
}

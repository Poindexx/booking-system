package storage

import (
    "context"
    "time"
    "github.com/jackc/pgx/v4/pgxpool"
)

type PostgresDB struct {
    Pool *pgxpool.Pool
}

// NewPostgresDB создает подключение к базе данных
func NewPostgresDB() (*PostgresDB, error) {
    pool, err := pgxpool.Connect(context.Background(), "postgres://user:password@localhost:5432/booking_db")
    if err != nil {
        return nil, err
    }

    return &PostgresDB{Pool: pool}, nil
}

// CreateReservation создает новое бронирование
func (db *PostgresDB) CreateReservation(res *model.Reservation) error {
    _, err := db.Pool.Exec(context.Background(),
        `INSERT INTO reservations (room_id, start_time, end_time) VALUES ($1, $2, $3)`,
        res.RoomID, res.StartTime, res.EndTime)
    return err
}

// GetReservations возвращает все бронирования для заданного зала
func (db *PostgresDB) GetReservations(roomID string) ([]model.Reservation, error) {
    rows, err := db.Pool.Query(context.Background(),
        `SELECT id, room_id, start_time, end_time FROM reservations WHERE room_id=$1`, roomID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var reservations []model.Reservation
    for rows.Next() {
        var res model.Reservation
        err := rows.Scan(&res.ID, &res.RoomID, &res.StartTime, &res.EndTime)
        if err != nil {
            return nil, err
        }
        reservations = append(reservations, res)
    }

    return reservations, nil
}

// IsConflict проверяет, есть ли конфликт с существующими бронированиями
func (db *PostgresDB) IsConflict(roomID string, startTime, endTime time.Time) (bool, error) {
    var count int
    err := db.Pool.QueryRow(context.Background(),
        `SELECT COUNT(*) FROM reservations WHERE room_id=$1 AND ($2 < end_time AND $3 > start_time)`,
        roomID, startTime, endTime).Scan(&count)
    if err != nil {
        return false, err
    }

    return count > 0, nil
}

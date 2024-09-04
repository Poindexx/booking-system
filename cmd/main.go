package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/Poindexx/booking-system/tree/main/internal/handler"

)

func main() {
	// Настройка подключения к базе данных
	db, err := storage.NewPostgresDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Инициализация роутера
	r := chi.NewRouter()

	// Регистрация обработчиков
	r.Post("/reservations", handler.CreateReservationHandler(db))
	r.Get("/reservations/{room_id}", handler.GetReservationsHandler(db))

	// Запуск сервера
	log.Println("Starting server on :8080...")
	http.ListenAndServe(":8080", r)
}

// @title TestEffectiveMobile API
// @version 1.0
// @host localhost:8085
// @BasePath /
package main

import (
	_ "TestEffectiveMobile/docs"
	"TestEffectiveMobile/internal/config"
	"TestEffectiveMobile/internal/handler"
	"TestEffectiveMobile/internal/repository"
	"TestEffectiveMobile/internal/usecase"
	"TestEffectiveMobile/migrations"
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	const op = "cmd.server.main"

	// Загрузка конфигурации из .env
	config.LoadEnv()

	// Подключение к Postgres
	connStr := config.GetDBConnectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error(op, "Ошибка подключения к БД", slog.String("error", err.Error()))
		return
	}
	defer db.Close()

	// Проверяем соединение
	if err = db.Ping(); err != nil {
		slog.Error(op, "Ошибка соединения с БД", slog.String("error", err.Error()))
		return
	}
	slog.Info("Успешное подключение к DB")

	// Запуск миграций
	if err = migrations.RunMigration(db); err != nil {
		os.Exit(1)
	}

	// Инициализация всех слоёв
	songRepo := repository.NewSongRepositoryyy()
	songUC := usecase.NewSongUseCase(songRepo, config.GetExternalAPIURL())
	songHandler := handler.NewSongHandler(songUC)

	// Настройка маршрутов
	r := mux.NewRouter()

	r.HandleFunc("/songs", songHandler.ListSongs).Methods("GET")             // Получение списка песен с фильтрацией и пагинацией
	r.HandleFunc("/songs/{id}", songHandler.DeleteSong).Methods("DELETE")    // Удаление песни
	r.HandleFunc("/songs/{id}", songHandler.UpdateSong).Methods("PUT")       // Изменение данных песни
	r.HandleFunc("/songs", songHandler.CreateSong).Methods("POST")           // Добавление новой песни с обогащения
	r.HandleFunc("/songs/{id}/text", songHandler.GetSongText).Methods("GET") // Получение текста песни с пагинацией

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler) // Маршрут для Swagger UI

	// Запуск HTTP-сервера
	port := os.Getenv("PORT")
	slog.Info("Сервер запускается на порту: " + port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error(op, "Ошибка сервера", slog.String("error", err.Error()))
		return
	}
}

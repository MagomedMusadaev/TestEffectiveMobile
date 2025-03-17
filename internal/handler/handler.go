package handler

import (
	"TestEffectiveMobile/internal/entities"
	"TestEffectiveMobile/internal/usecase"
	"encoding/json"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"strconv"
)

type SongHandler interface {
	ListSongs(w http.ResponseWriter, r *http.Request)
	DeleteSong(w http.ResponseWriter, r *http.Request)
	UpdateSong(w http.ResponseWriter, r *http.Request)
	CreateSong(w http.ResponseWriter, r *http.Request)
	GetSongText(w http.ResponseWriter, r *http.Request)
}

type songHandler struct {
	useCase usecase.SongUseCase
}

func NewSongHandler(useCase usecase.SongUseCase) SongHandler {
	return &songHandler{
		useCase: useCase,
	}
}

// ListSongs godoc
// @Summary Получение списка песен
// @Description Получает список песен с возможностью фильтрации по группе и названию, а также с пагинацией.
// @Tags songs
// @Produce json
// @Param group query string false "Название группы"
// @Param song_title query string false "Название песни"
// @Param limit query int false "Лимит записей (по умолчанию 11)"
// @Param offset query int false "Сдвиг записей"
// @Success 200 {array} entities.Song "Список песен"
// @Failure 500 {object} entities.ErrorResponse "Ошибка сервера"
// @Router /songs [get]
func (h *songHandler) ListSongs(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handler.ListSongs"

	query := r.URL.Query()
	filter := make(map[string]string)

	if group := query.Get("group"); group != "" {
		filter["group_name"] = group
	}
	if song := query.Get("song_title"); song != "" {
		filter["song_title"] = song
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 11
	}
	offset, _ := strconv.Atoi(query.Get("offset"))

	songs, err := h.useCase.ListSongs(filter, limit, offset)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(songs)
}

// DeleteSong godoc
// @Summary Удаление песни
// @Description Удаляет песню по идентификатору.
// @Tags songs
// @Param id path int true "ID песни"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} entities.ErrorResponse "Неверный ID"
// @Failure 500 {object} entities.ErrorResponse "Internal Server Error"
// @Router /songs/{id} [delete]
func (h *songHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handler.DeleteSong"

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		slog.Error(op, "Ошибка парсинга id с параметров", slog.String("error", err.Error()))
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if err = h.useCase.DeleteSong(id); err != nil {
		slog.Error(op, "Ошибка при удаления песни", slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UpdateSong godoc
// @Summary Обновление данных песни
// @Description Обновляет данные песни по идентификатору. Передаётся JSON объект песни.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param song body entities.Song true "Обновленные данные песни"
// @Success 200 {string} string "OK"
// @Failure 400 {object} entities.ErrorResponse "Неверный ID или Bad Request"
// @Failure 500 {object} entities.ErrorResponse "Internal Server Error"
// @Router /songs/{id} [put]
func (h *songHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handler.UpdateSong"

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		slog.Error(op, "Ошибка парсинга id с параметров", slog.String("error", err.Error()))
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var song entities.Song
	if err = json.NewDecoder(r.Body).Decode(&song); err != nil {
		slog.Error(op, "Ошибка декодинга данных", slog.String("error", err.Error()))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	song.ID = id
	if err = h.useCase.UpdateSong(song); err != nil {
		slog.Error(op, "Ошибка обновления данных песни", slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// CreateSong godoc
// @Summary Добавление новой песни
// @Description Создает новую песню, выполняет вызов внешнего API для обогащения данных и сохраняет в базу.
// @Tags songs
// @Accept json
// @Produce json
// @Param song body entities.Song true "Данные новой песни"
// @Success 201 {object} entities.Song "Созданная песня"
// @Failure 400 {object} entities.ErrorResponse "Bad Request"
// @Failure 500 {object} entities.ErrorResponse "Internal Server Error"
// @Router /songs [post]
func (h *songHandler) CreateSong(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handler.CreateSong"

	var song entities.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		slog.Error(op, "Ошибка декодинга данных", slog.String("error", err.Error()))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	id, err := h.useCase.CreateSong(song)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	song.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}

// GetSongText godoc
// @Summary Получение текста песни с пагинацией куплетов
// @Description Возвращает текст песни, разделенный на куплеты с пагинацией. Параметры versePage и versePageSize управляют выводом куплетов.
// @Tags songs
// @Produce plain
// @Param id path int true "ID песни"
// @Param versePage query int false "Номер страницы куплетов (по умолчанию 1)"
// @Param versePageSize query int false "Количество куплетов на странице (по умолчанию 5)"
// @Success 200 {string} string "Текст песни"
// @Failure 400 {object} entities.ErrorResponse "Неверный ID"
// @Failure 404 {object} entities.ErrorResponse "Песня не найдена"
// @Failure 500 {object} entities.ErrorResponse "Internal Server Error"
// @Router /songs/{id}/text [get]
func (h *songHandler) GetSongText(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handler.GetSongText"

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		slog.Error(op, "Ошибка парсинга id с параметров", slog.String("error", err.Error()))
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	song, err := h.useCase.GetSongByID(id)
	if err != nil {
		http.Error(w, "Песня не найдена", http.StatusNotFound)
		return
	}

	// Формируем пагинацию
	versePage, _ := strconv.Atoi(r.URL.Query().Get("versePage"))
	if versePage == 0 {
		versePage = 1
	}
	versePageSize, _ := strconv.Atoi(r.URL.Query().Get("versePageSize"))
	if versePageSize == 0 {
		versePageSize = 5
	}

	text, err := h.useCase.GetSongText(song, versePage, versePageSize)
	if err != nil {
		slog.Error(op, "Ошибка получения куплетов", slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(text))
}

package usecase

import (
	"TestEffectiveMobile/internal/entities"
	"TestEffectiveMobile/internal/repository"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

type SongUseCase interface {
	ListSongs(filter map[string]string, limit, offset int) ([]entities.Song, error)
	DeleteSong(id int) error
	UpdateSong(song entities.Song) error
	CreateSong(song entities.Song) (int, error)
	GetSongByID(id int) (*entities.Song, error)
	GetSongText(song *entities.Song, versePage, versePageSize int) (string, error)
}

type songUseCase struct {
	repo           repository.SongRepository
	externalAPIURL string
}

func NewSongUseCase(repo repository.SongRepository, externalAPIURL string) SongUseCase {
	return &songUseCase{
		repo:           repo,
		externalAPIURL: externalAPIURL,
	}
}

func (u *songUseCase) ListSongs(filter map[string]string, limit, offset int) ([]entities.Song, error) {
	return u.repo.ListSongs(filter, limit, offset)
}

func (u *songUseCase) DeleteSong(id int) error {
	return u.repo.DeleteSong(id)
}

func (u *songUseCase) UpdateSong(song entities.Song) error {
	return u.repo.UpdateSong(song)
}

func (u *songUseCase) CreateSong(song entities.Song) (int, error) {
	const op = "internal.useCase.CreateSong"

	// Вызываем внешний API для обогащения
	enriched, err := u.getSongInfo(song.Group, song.Title)
	if err != nil {
		slog.Error(op, "Ошибка внешнего API", slog.String("error", err.Error()))
		return 0, err
	}
	song.ReleaseDate = enriched.ReleaseDate
	song.Text = enriched.Text
	song.Link = enriched.Link

	return u.repo.CreateSong(song)
}

func (u *songUseCase) getSongInfo(group, song string) (*entities.ExternalSongInfo, error) {
	const op = "internal.useCase.GetSongInfo"

	params := url.Values{}
	params.Add("group", group)
	params.Add("song", song)
	fullURL := fmt.Sprintf("%s?%s", u.externalAPIURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		slog.Error(op, "Ошибка запроса обогащения", slog.String("error", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("внешнее API вернуло статус %d", resp.StatusCode)
		slog.Error(op, "Не тот статус код", slog.String("error", err.Error()))
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(op, "Ошибка парсинга тела ответа", slog.String("error", err.Error()))
		return nil, err
	}

	var info entities.ExternalSongInfo
	if err = json.Unmarshal(body, &info); err != nil {
		slog.Error(op, "Ошибка анмаршалинга данных", slog.String("error", err.Error()))
		return nil, err
	}

	return &info, nil
}

func (u *songUseCase) GetSongByID(id int) (*entities.Song, error) {
	return u.repo.GetSongByID(id)
}

func (u *songUseCase) GetSongText(song *entities.Song, versePage, versePageSize int) (string, error) {
	const op = "internal.useCase.GetSongText"

	if song.Text == "" {
		slog.Info("Отсутствует текст песни", "songID", song.ID)
		return "", nil
	}
	verses := strings.Split(song.Text, "\n")
	totalVerses := len(verses)
	slog.Debug("Начало пагинации куплетов",
		"songID", song.ID,
		"totalVerses", totalVerses,
		"versePage", versePage,
		"versePageSize", versePageSize,
	)

	start := (versePage - 1) * versePageSize
	end := start + versePageSize

	if start >= totalVerses {
		slog.Info("Индекс начала пагинации превышает общее количество куплетов",
			"songID", song.ID,
			"start", start,
			"totalVerses", totalVerses,
		)
		return "", nil
	}
	if end > totalVerses {
		end = totalVerses
		slog.Debug("Корректировка индекса окончания пагинации",
			"songID", song.ID,
			"adjustedEnd", end,
		)
	}
	paginated := verses[start:end]
	result := strings.Join(paginated, "\n")
	slog.Debug("Результат пагинации куплетов",
		"songID", song.ID,
		"start", start,
		"end", end,
		"result", result,
	)
	return result, nil
}

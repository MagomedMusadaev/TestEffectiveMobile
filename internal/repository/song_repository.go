package repository

import (
	"TestEffectiveMobile/internal/entities"
	"database/sql"
	"fmt"
	"log/slog"
)

type SongRepository interface {
	ListSongs(filter map[string]string, limit, offset int) ([]entities.Song, error)
	DeleteSong(id int) error
	UpdateSong(song entities.Song) error
	CreateSong(song entities.Song) (int, error)
	GetSongByID(id int) (*entities.Song, error)
}

type songRepository struct {
	db *sql.DB
}

func NewSongRepository(db *sql.DB) SongRepository {
	return &songRepository{
		db: db,
	}
}

func (r *songRepository) ListSongs(filter map[string]string, limit, offset int) ([]entities.Song, error) {
	const op = "internal.repository.ListSongs"

	query := `SELECT id, group_name, song_title, release_date, text, link FROM songs WHERE 1=1`

	// Формируем строку запроса к DB
	args := make([]interface{}, 0)
	i := 1
	for key, value := range filter {
		query += fmt.Sprintf(" AND %s=$%d", key, i)
		args = append(args, value)
		i++
	}
	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, limit, offset)

	// Делаем запрос к DB
	rows, err := r.db.Query(query, args...)
	if err != nil {
		slog.Error(op, "Ошибка запроса к DB", slog.String("error", err.Error()))
		return nil, err
	}

	songs := make([]entities.Song, 0)
	for rows.Next() {
		var song entities.Song
		if err = rows.Scan(
			&song.ID,
			&song.Group,
			&song.Title,
			&song.ReleaseDate,
			&song.Text,
			&song.Link,
		); err != nil {
			slog.Error(op, "Ошибка сканирования результата", slog.String("error", err.Error()))
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil
}

func (r *songRepository) DeleteSong(id int) error {
	const op = "internal.repository.DeleteSong"

	query := `DELETE FROM songs WHERE id = $1`

	if _, err := r.db.Exec(query, id); err != nil {
		slog.Error(op, "Ошибка при удалении записи с DB", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *songRepository) UpdateSong(song entities.Song) error {
	const op = "internal.repository.UpdateSong"

	query := `UPDATE songs SET group_name=$1, song_title=$2, release_date=$3, text=$4, link=$5 WHERE id=$6`
	_, err := r.db.Exec(query, song.Group, song.Title, song.ReleaseDate, song.Text, song.Link, song.ID)
	if err != nil {
		slog.Error(op, "Ошибка при изменении данных в DB", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *songRepository) CreateSong(song entities.Song) (int, error) {
	const op = "internal.repository.CreateSong"

	query := `INSERT INTO songs (group_name, song_title, release_date, text, link)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var id int
	if err := r.db.QueryRow(query,
		song.Group,
		song.Title,
		song.ReleaseDate,
		song.Text,
		song.Link).Scan(&id); err != nil {
		slog.Error(op, "Ошибка изменения данных", slog.String("error", err.Error()))
		return 0, nil
	}
	return id, nil
}

func (r *songRepository) GetSongByID(id int) (*entities.Song, error) {
	const op = "internal.repository.GetSongByID"

	query := `SELECT * FROM WHERE id=$1`

	row := r.db.QueryRow(query, id)

	var song entities.Song
	if err := row.Scan(&song.ID,
		&song.Group,
		&song.Title,
		&song.ReleaseDate,
		&song.Text,
		&song.Link,
	); err != nil {
		slog.Error(op, "Ошибка парсинга данных", slog.String("error", err.Error()))
		return nil, err
	}
	return &song, nil
}

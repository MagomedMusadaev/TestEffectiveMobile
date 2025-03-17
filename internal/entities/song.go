package entities

// Song представляет информацию о песне.
// @Description Структура для представления песни, которая включает 6 полей.
// swagger:model Song
type Song struct {
	// ID уникальный идентификатор песни.
	//
	// required: true
	//
	// example: 1
	ID int `json:"id"`

	// Group название группы или исполнителя.
	//
	// required: true
	//
	// example: "The Beatles"
	Group string `json:"group"`

	// Title название песни.
	//
	// required: true
	//
	// example: "Hey Jude"
	Title string `json:"song"`

	// ReleaseDate дата выпуска песни в формате YYYY-MM-DD.
	//
	// example: "2023-01-01"
	ReleaseDate string `json:"releaseDate,omitempty"`

	// Text текст песни.
	//
	// example: "Hey, Jude, don't make it bad..."
	Text string `json:"text,omitempty"`

	// Link ссылка на дополнительную информацию о песне.
	//
	// example: "https://example.com/song-info"
	Link string `json:"link,omitempty"`
}

func NewSong(
	id int,
	group, title, realiseDate, text, link string,
) *Song {
	return &Song{
		ID:          id,
		Group:       group,
		Title:       title,
		ReleaseDate: realiseDate,
		Text:        text,
		Link:        link,
	}
}

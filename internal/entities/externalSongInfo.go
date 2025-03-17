package entities

// ExternalSongInfo представляет внешнюю информацию о песне.
//
// swagger:model ExternalSongInfo
type ExternalSongInfo struct {
	// ReleaseDate дата выпуска песни в формате YYYY-MM-DD.
	//
	// required: true
	// example: 2023-01-01
	ReleaseDate string `json:"releaseDate"`
	// Text текст песни.
	//
	// required: true
	Text string `json:"text"`
	// Link ссылка на дополнительную информацию о песне.
	//
	// required: true
	// example: https://example.com/song-info
	Link string `json:"link"`
}

func NewEternalSongInfo(releaseDate, text, link string) *ExternalSongInfo {
	return &ExternalSongInfo{
		ReleaseDate: releaseDate,
		Text:        text,
		Link:        link,
	}
}

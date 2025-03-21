package models

type Song struct {
	BandName    string `json:"band_name" validate:"required" db:"band"`
	SongTitle   string `json:"song_title" validate:"required" db:"song"`
	ReleaseDate string `json:"release_date,omitempty" db:"release"`
	Lyrics      string `json:"lyrics,omitempty" db:"lyrics"`
	Link        string `json:"link,omitempty" db:"link"`
}

type SongFilter struct {
	BandName    string `json:"band_name,omitempty" db:"band"`
	SongTitle   string `json:"song_title,omitempty" db:"song"`
	ReleaseDate string `json:"release_date,omitempty" db:"release"`
	Later       bool   `json:"bigger,omitempty"`
	Lyrics      string `json:"lyrics,omitempty" db:"lyrics"`
	Page        int    `json:"page,omitempty"`
	PageSize    int    `json:"page_size,omitempty"`
}

type SongLyrics struct {
	BandName  string `json:"band_name" validate:"required"`
	SongTitle string `json:"song_title" validate:"required"`
	Verse     int    `json:"verse" validate:"required"`
}

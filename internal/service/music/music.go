package music

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/stepan41k/Testovoe/internal/domain/models"
	"github.com/stepan41k/Testovoe/internal/lib/api/logger/sl"
	"github.com/stepan41k/Testovoe/internal/storage"
)

type Music interface {
	GetSongs(ctx context.Context, song models.SongFilter) (songs []models.Song, err error)
	GetTextSong(ctx context.Context, song models.SongLyrics) (verse string, err error)
	DeleteSong(ctx context.Context, song models.Song) (id int64, err error)
	UpdateSong(ctx context.Context, songDetails models.Song) (id int64, err error)
	AddNewSong(ctx context.Context, song models.Song) (id int64, err error)
}

type MusicService struct {
	music Music
	log *slog.Logger
}

func New(music Music, log *slog.Logger) *MusicService {
	return &MusicService{
		music: music,
		log: log,
	}
}


func (m *MusicService) GetSongs(ctx context.Context, song models.SongFilter) ([]models.Song, error) {
	const op = "service.music.GetSongs"

	log := m.log.With(
		slog.String("op", op),
		slog.String("group", song.BandName),
		slog.String("song", song.SongTitle),
	)

	log.Info("getting songs")

	songs, err := m.music.GetSongs(ctx, song)
	if err != nil {
		log.Error("failed to get songs")

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got songs")

	return songs, nil
}


func (m *MusicService) GetTextSong(ctx context.Context, song models.SongLyrics) (string, error) {
	const op = "service.music.GetTextSong"

	log := m.log.With(
		slog.String("op", op),
		slog.String("group", song.BandName),
		slog.String("song", song.SongTitle),
	)

	log.Info("getting text of song")

	verse, err := m.music.GetTextSong(ctx, song)
	if err != nil {
		log.Error("failed to get text of song")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got text of song")

	return verse, nil
}


func (m *MusicService) DeleteSong(ctx context.Context, song models.Song) (int64, error) {
	const op = "service.music.DeleteSong"

	log := m.log.With(
		slog.String("op", op),
		slog.String("group", song.BandName),
		slog.String("song", song.SongTitle),
	)

	log.Info("deleting song")

	id, err := m.music.DeleteSong(ctx, song)
	if err != nil {
		log.Error("failed to delete song")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("song deleted")

	return id, nil
}


func (m *MusicService) UpdateSong(ctx context.Context, songDetails models.Song)(int64, error)  {
	const op = "service.music.UpdateSong"

	log := m.log.With(
		slog.String("op", op),
		slog.String("group", songDetails.BandName),
		slog.String("song", songDetails.SongTitle),
	)

	log.Info("updating song")

	id, err := m.music.UpdateSong(ctx, songDetails)
	if err != nil {
		log.Error("failed to update song")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("song updated")

	return id, nil
}


func (m *MusicService) AddNewSong(ctx context.Context, song models.Song) (int64, error)  {
	const op = "service.music.AddNewSong"

	log := m.log.With(
		slog.String("op", op),
		slog.String("group", song.BandName),
		slog.String("song", song.SongTitle),
	)

	log.Info("adding new song")

	id, err := m.music.AddNewSong(ctx, song)
	if err != nil {
		if errors.Is(err, storage.ErrSongExists) {
			log.Warn("song already exists", sl.Err(err))

			return 0, fmt.Errorf("%s: %w", op, storage.ErrSongExists)
		}

		log.Error("failed to add song", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("song added")

	return id, nil
	
}
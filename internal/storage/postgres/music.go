package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/stepan41k/Testovoe/internal/domain/models"
	"github.com/stepan41k/Testovoe/internal/storage"
)

const (
	sizeOfPage = 10
	sizeOfVerse = 1
)


func (s *PStorage) GetSongs(ctx context.Context, song models.Song) ([]models.Song, error) {
	const op = "storage.postgres.music.GetSongs"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", op, commitErr)
		}
	}()

	rows, err := tx.Query(ctx, `
		SELECT * FROM songs
		WHERE band LIKE $1
		LIMIT $2
		OFFSET $3;
	`, song.Group, sizeOfPage, (song.Page-1)*sizeOfPage) 

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var item models.Song
		err = rows.Scan(
			&item.Group,
			&item.Song,
			&item.ReleaseDate,
			&item.Text,
			&item.Link,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		songs = append(songs, item)
	}

	return songs, err
}


func (s *PStorage) GetTextSong(ctx context.Context, song models.Song) (models.Verse, error) {
	const op = "storage.postgres.music.GetTextSong"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return models.Verse{}, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", op, commitErr)
		}
	}()

	row := tx.QueryRow(ctx, `
		WITH split_text AS (
                SELECT id, song, unnest(regexp_split_to_array(text, E'\n\n'))
        AS verse
                FROM songs
        )
        SELECT id, song, verse
        FROM split_text
		WHERE song = $1 AND band = $2
        LIMIT $3 OFFSET $4;
	`, song.Song, song.Group, sizeOfVerse, song.Page)

	var verse models.Verse

	err = row.Scan(&verse.ID, &verse.Song, &verse.Verse)
	if err != nil {
		return models.Verse{}, fmt.Errorf("%s: %w", op, err)
	}

	return verse, nil
}


func (s *PStorage) DeleteSong(ctx context.Context, song models.Song) (id int64, err error) {
	const op = "storage.postgres.music.DeleteSong"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", op, commitErr)
		}
	}()

	row := tx.QueryRow(ctx, `
		DELETE FROM songs
		WHERE song = $1 AND band = $2
	`, song.Song, song.Group)

	err = row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}


func (s *PStorage) UpdateSong(ctx context.Context, song models.Song) (id int64, err error) {
	const op = "storage.postgres.music.UpdateSong"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer func ()  {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", op, commitErr)
		}
	}()

	row := tx.QueryRow(ctx, `
		UPDATE songs (releaseDate, text, link)
		SET ($1, $2, $3)
		WHERE song = $4 AND band = $5;
	`, song.ReleaseDate, song.Text, song.Link, song.Song, song.Group)

	err = row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}


func (s *PStorage) AddNewSong(ctx context.Context, song models.Song) (id int64, err error) {
	const op = "storage.postgres.music.AddNewSong"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", op, commitErr)
		}
	}()

	row := tx.QueryRow(ctx, `
		INSERT INTO songs (band, song)
		VALUES ($1, $2)
		RETURNING id;
	`, song.Group, song.Song)

	err = row.Scan(&id)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrSongExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
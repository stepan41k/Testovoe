package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
	"github.com/stepan41k/Testovoe/internal/domain/models"
	"github.com/stepan41k/Testovoe/internal/storage"
)

const (
	sizeOfVerse = 1
)


func (s *PStorage) GetSongs(ctx context.Context, song models.SongFilter) ([]models.Song, error) {
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
	
	arguments, values, ind := []string{}, []any{}, 1
	query := `SELECT band, song, TO_CHAR(release, 'DD.MM.YYYY'), lyrics, link FROM songs`

	if song.BandName != "" || song.SongTitle != "" || song.ReleaseDate != "" {
		query += ` WHERE `
		switch {
		case song.SongTitle != "":
			arguments = append(arguments, fmt.Sprintf(`song LIKE $%d`, ind))
			values = append(values, song.SongTitle)
			ind++
		case song.BandName != "":
			arguments = append(arguments, fmt.Sprintf(`band LIKE $%d`, ind))
			values = append(values, song.BandName)
			ind++
		case song.ReleaseDate != "":
			if song.Later {
				arguments = append(arguments, fmt.Sprintf(`release > TO_DATE($%d, 'DD.MM.YYYY')`, ind))
			} else if !song.Later {
				arguments = append(arguments, fmt.Sprintf(`release <= TO_DATE($%d, 'DD.MM.YYYY')`, ind))
			} else {
				return nil, fmt.Errorf("%s: %w", op, storage.ErrNoChanges)
			}
			values = append(values, song.ReleaseDate)
			ind++
		}
		query += strings.Join(arguments, ",")
	} 

	query += fmt.Sprintf(` LIMIT $%d OFFSET $%d;`, ind, ind+1)
	values = append(values, song.PageSize, (song.Page-1)*song.PageSize)


	rows, err := tx.Query(ctx, query, values...)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrSongNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var item models.Song
		err = rows.Scan(&item)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		songs = append(songs, item)
	}
	
	return songs, err
}


func (s *PStorage) GetTextSong(ctx context.Context, song models.SongLyrics) (string, error) {
	const op = "storage.postgres.music.GetTextSong"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
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
                SELECT song, band, unnest(regexp_split_to_array(lyrics, E'\n\n'))
        AS verse
                FROM songs
        )
        SELECT verse
        FROM split_text
		WHERE song = $1 AND band = $2
        LIMIT $3 OFFSET $4;
	`, song.SongTitle, song.BandName, sizeOfVerse, song.Verse-1)

	var verse string

	err = row.Scan(&verse)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
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
		RETURNING id;
	`, song.SongTitle, song.BandName)

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

	arguments, values, ind := []string{}, []any{}, 1
	query := `UPDATE songs SET updated = NOW(), `

	if song.BandName != "" || song.SongTitle != "" || song.ReleaseDate != "" {
		switch {
		case song.Link != "":
			arguments = append(arguments, fmt.Sprintf(`link = $%d`, ind))
			values = append(values, song.Link)
			ind++
		case song.Lyrics != "":
			arguments = append(arguments, fmt.Sprintf(`lyrics = $%d`, ind))
			values = append(values, song.Lyrics)
			ind++
		case song.ReleaseDate != "":
			arguments = append(arguments, fmt.Sprintf(`release = TO_DATE($%d, 'DD.MM.YYYY')`, ind))
			values = append(values, song.ReleaseDate)
			ind++
		}
	} else {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrNoChanges)
	}

	query += strings.Join(arguments, ",")
	query += fmt.Sprintf(` WHERE song = $%d AND band = $%d RETURNING id;`, ind, ind+1)
	values = append(values, song.SongTitle, song.BandName)


	row := tx.QueryRow(ctx, query, values...)
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
		INSERT INTO songs (band, song, updated)
		VALUES ($1, $2, NOW())
		RETURNING id;
	`, song.BandName, song.SongTitle)

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
CREATE TABLE IF NOT EXISTS
    songs (
        id SERIAL PRIMARY KEY,
        band TEXT NOT NULL,
        song TEXT NOT NULL,
        release DATE,
        lyrics TEXT,
        link TEXT,
        updated TIMESTAMP
    );

CREATE UNIQUE INDEX unique_song_of_band ON songs(band, song);
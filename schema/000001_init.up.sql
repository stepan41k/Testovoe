CREATE TABLE IF NOT EXISTS
    songs (
        "id" SERIAL PRIMARY KEY,
        "band" TEXT NOT NULL,
        "song" TEXT NOT NULL,
        "releaseDate" DATE,
        "text" TEXT,
        "link" TEXT
    );

CREATE UNIQUE INDEX unique_song_of_band ON songs(band, song);
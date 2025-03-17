-- 20250317133548_create_songs_table.up.sql
CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    song_title VARCHAR(255) NOT NULL,
    release_date VARCHAR(50),
    text TEXT,
    link VARCHAR(255)
    );

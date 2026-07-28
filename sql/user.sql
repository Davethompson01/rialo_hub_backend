CREATE TYPE user_role AS ENUM ('Writer', 'Developer', 'Artist', 'Moderators', '3dArtist');

CREATE TABLE users (
    user_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    discord_username TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL,
    role user_role,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

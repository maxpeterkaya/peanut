ALTER TABLE release
    ADD COLUMN repository_id INTEGER NOT NULL REFERENCES repository (id);

ALTER TABLE asset
    ADD COLUMN release_id INTEGER NOT NULL REFERENCES release (id);

ALTER TABLE repository
    ADD COLUMN user_id INTEGER NOT NULL REFERENCES "user" (id);
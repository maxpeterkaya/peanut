CREATE TABLE repository
(
    id         serial PRIMARY KEY,
    github_id int UNIQUE NOT NULL,
    user_id   int references "user" (id),
    owner      text,
    name       text,
    token      text,
    is_private bool      default false,
    created_at timestamp default now(),
    updated_at timestamp default now()
);

CREATE TABLE release
(
    id                serial PRIMARY KEY,
    github_id int UNIQUE NOT NULL,
    name              text,
    tag_name          text,
    body              text,
    is_draft          bool default false,
    is_prerelease     bool default false,
    created_at        timestamp,
    published_at      timestamp,
    author_name       text,
    author_id         text,
    author_avatar_url text,
    repository_id     int references repository (id)
);

CREATE TABLE asset
(
    id             serial PRIMARY KEY,
    github_id int UNIQUE NOT NULL,
    api_url        text,
    url            text,
    name           text,
    content_length int,
    download_count int,
    view_count     int,
    created_at     timestamp,
    updated_at     timestamp,
    uploaded_at    timestamp,
    release_id     int references release (id)
);

CREATE TABLE "user"
(
    id           serial PRIMARY KEY,
    username     VARCHAR(40),
    display_name VARCHAR(40),
    pass_hash    text,
    created_at   timestamp,
    updated_at   timestamp
);
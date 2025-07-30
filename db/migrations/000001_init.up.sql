CREATE TABLE repository
(
    id         serial PRIMARY KEY,
    owner      text,
    name       text,
    token      text,
    is_private bool      default false,
    created_at timestamp default now(),
    updated_at timestamp default now()
);

CREATE TABLE release
(
    id                int PRIMARY KEY,
    name              text,
    tag_name          text,
    body              text,
    is_draft          bool default false,
    is_prerelease     bool default false,
    created_at        timestamp,
    published_at      timestamp,
    author_name       text,
    author_id         text,
    author_avatar_url text
);

CREATE TABLE asset
(
    id             int PRIMARY KEY,
    api_url        text,
    url            text,
    name           text,
    content_length int,
    download_count int,
    view_count     int,
    created_at     timestamp,
    updated_at     timestamp,
    uploaded_at    timestamp
);

CREATE TABLE "user"
(
    id           serial PRIMARY KEY,
    username     VARCHAR(40),
    display_name VARCHAR(40),
    pass_hash    text,
    created_at   timestamp default now(),
    updated_at   timestamp default now()
);

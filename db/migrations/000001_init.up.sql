CREATE TABLE repository
(
    id        uuid PRIMARY KEY,
    owner     text,
    name      text,
    token     text,
    isPrivate bool      default false,
    createdAt timestamp default now(),
    updatedAt timestamp default now()
);

CREATE TABLE release (
                         id int PRIMARY KEY,
                         name text,
                         tagName text,
                         body text,
                         draft bool default false,
                         prerelease bool default false,
                         createdAt timestamp,
                         publishedAt timestamp,
                         authorName text,
                         authorId text,
                         authorAvatarUrl text
);

CREATE TABLE asset
(
    id            int,
    apiURL        text,
    url           text,
    name          text,
    contentLength int,
    downloadCount int,
    viewCount     int,
    createdAt     timestamp,
    updatedAt     timestamp,
    uploadedAt    timestamp
);

CREATE TABLE "user"
(
    id          uuid DEFAULT gen_random_uuid(),
    username    VARCHAR(40),
    displayName VARCHAR(40),
    passHash    text,
    createdAt   timestamp,
    updatedAt   timestamp
);

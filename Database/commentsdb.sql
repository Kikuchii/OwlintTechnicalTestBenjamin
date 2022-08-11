BEGIN;

SET client_encoding = 'LATIN1';

CREATE TABLE comments (
    id integer PRIMARY KEY,
    textFr text,
    textEn text,
    publishedAt varchar NOT NULL,
    authorId varchar,
    tagetId varchar NOT NULL
);

COMMIT;
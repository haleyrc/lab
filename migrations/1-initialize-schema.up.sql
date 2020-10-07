CREATE TABLE IF NOT EXISTS authors (
	id           UUID        PRIMARY KEY,

	first_name   TEXT        NOT NULL DEFAULT '',
	middle_name  TEXT        NOT NULL DEFAULT '',
	last_name    TEXT        NOT NULL DEFAULT '',
	date_created TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS genres (
	id           UUID        PRIMARY KEY,

	name         TEXT        NOT NULL DEFAULT '',
	date_created TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tags (
	id           UUID        PRIMARY KEY,

	name         TEXT        NOT NULL DEFAULT '',
	description  TEXT        NOT NULL DEFAULT '',
	color        TEXT        NOT NULL DEFAULT 'dddddd',
	date_created TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS books (
	id           UUID        PRIMARY KEY,

	title        TEXT        NOT NULL DEFAULT '',
	pub_year     INTEGER     NOT NULL DEFAULT extract(year from NOW()),
	date_created TIMESTAMPTZ NOT NULL DEFAULT NOW(),

	genre_id     UUID        REFERENCES genres(id),
	author_id    UUID        REFERENCES authors(id)
);

CREATE TABLE IF NOT EXISTS book_tags (
	book_id      UUID        NOT NULL REFERENCES books(id),
	tag_id       UUID        NOT NULL REFERENCES tags(id),

	date_created TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
	id            UUID        PRIMARY KEY,

	email         TEXT        NOT NULL DEFAULT '',
	password_hash TEXT        NOT NULL DEFAULT '',
	date_created  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_books (
	user_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	book_id      UUID        NOT NULL REFERENCES books(id) ON DELETE CASCADE,

	date_created TIMESTAMPTZ NOT NULL DEFAULT NOW(),

	PRIMARY KEY (user_id, book_id)
);

INSERT INTO authors (id, first_name, middle_name, last_name) VALUES
    ('d1b68ffe-45de-4810-9c3d-4ec7da0dfce9', 'Frank', 'Patrick', 'Herbert'),
    ('997268db-725a-4862-a75a-526069c21a88', 'J.R.R.', '', 'Tolkien');


INSERT INTO genres (id, name) VALUES
    ('6c59c66b-8d17-4327-aa5b-cc77acba53f8', 'Science Fiction'),
    ('d9106b49-3dec-4cee-8999-b277bf44d2dd', 'Fantasy');

INSERT INTO tags (id, name, description, color) VALUES
    ('c0f55b16-fa7e-4a98-a286-fb3ce927eb28', 'Favorite', 'One of my favorites.', 'BADA55'),
    ('e0f8d13f-4cdf-46d2-98fd-cee52b4a277f', 'Bedtime', 'Something we could read to the kids.', '31ADC4');

INSERT INTO books (id, title, pub_year, genre_id, author_id) VALUES
    ('e8742c65-f27b-4cc4-94ce-c452cacff119', 'Dune', 1965, '6c59c66b-8d17-4327-aa5b-cc77acba53f8', 'd1b68ffe-45de-4810-9c3d-4ec7da0dfce9'),
    ('0b5a82ce-1df6-471c-a41e-c1c4524cd3a2', 'The Silmarillion', 1977, 'd9106b49-3dec-4cee-8999-b277bf44d2dd', '997268db-725a-4862-a75a-526069c21a88');

INSERT INTO book_tags (book_id, tag_id) VALUES
    ('e8742c65-f27b-4cc4-94ce-c452cacff119', 'c0f55b16-fa7e-4a98-a286-fb3ce927eb28'),
    ('0b5a82ce-1df6-471c-a41e-c1c4524cd3a2', 'c0f55b16-fa7e-4a98-a286-fb3ce927eb28'),
    ('0b5a82ce-1df6-471c-a41e-c1c4524cd3a2', 'e0f8d13f-4cdf-46d2-98fd-cee52b4a277f');

INSERT INTO users ( id, email, password_hash ) VALUES
    ('a18b97b8-cb6d-49da-8a04-59ecdee21dd5', 'ryan@haleylab.com', '$2a$10$UzTxNiR4C3qB2PYcbhgpU.1ziAOk7vsispSiBqS7gLuZizvlQcEmm');
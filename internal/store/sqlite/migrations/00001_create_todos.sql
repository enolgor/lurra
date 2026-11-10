-- +goose Up
CREATE TABLE todos (
	id    INTEGER PRIMARY KEY,
	title TEXT    NOT NULL CHECK (title <> ''),
	done  BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE todos;

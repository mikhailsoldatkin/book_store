-- +goose Up
CREATE TABLE books
(
    id      BIGSERIAL PRIMARY KEY,
    author  TEXT                     NOT NULL,
    title   TEXT                     NOT NULL,
    price   FLOAT,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS books;

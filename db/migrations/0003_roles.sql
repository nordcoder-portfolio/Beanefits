-- +goose Up
CREATE TABLE roles
(
    code TEXT PRIMARY KEY
);

-- +goose Down
DROP TABLE roles;

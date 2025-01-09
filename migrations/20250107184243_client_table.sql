-- +goose Up
-- +goose StatementBegin
CREATE TABLE clients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    client_name TEXT NOT NULL
);
-- +goose StatementEnd
-- INSERT INTO clients ( client_name) VALUES ( 'Alice');
-- INSERT INTO clients ( client_name) VALUES ( 'Bob');
-- INSERT INTO clients ( client_name) VALUES ( 'Charlie');
-- INSERT INTO clients ( client_name) VALUES ( 'Diana');

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS clients;
-- +goose StatementEnd

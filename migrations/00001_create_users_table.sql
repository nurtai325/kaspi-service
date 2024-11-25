-- +goose Up
-- +goose StatementBegin
CREATE TABLE users();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD role VARCHAR(20) NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

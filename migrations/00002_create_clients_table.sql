-- +goose Up
-- +goose StatementBegin
CREATE TABLE clients (
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	token VARCHAR(50) NOT NULL,
	phone VARCHAR(15) NOT NULL,
	expires TIMESTAMP NOT NULL,
	jid VARCHAR(50) NOT NULL,
	connected BOOLEAN NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS clients;
-- +goose StatementEnd

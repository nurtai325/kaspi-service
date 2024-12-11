package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/nurtai325/kaspi-service/internal/config"
	"github.com/nurtai325/kaspi-service/internal/db"
	"github.com/nurtai325/kaspi-service/internal/models"
)

var (
	ErrGetAllClients = errors.New("error getting all clients from the database")
)

type Client struct {
	conn *sql.DB
}

func NewClient() Client {
	return Client{
		conn: db.Conn(config.New()),
	}
}

func (c *Client) All(ctx context.Context) ([]models.Client, error) {
	return c.all(ctx, "")
}

func (c *Client) AllNotExpired(ctx context.Context) ([]models.Client, error) {
	return c.all(ctx, "WHERE expires > $1", time.Now().UTC())
}

func (c *Client) all(ctx context.Context, whereClause string, args ...any) ([]models.Client, error) {
	statement := fmt.Sprintf("SELECT id, name, token, phone, expires, jid, connected FROM clients %s ORDER BY id ASC;", whereClause)
	rows, err := c.conn.QueryContext(ctx, statement, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Join(ErrGetAllClients, err)
	}
	defer rows.Close()
	clients := make([]models.Client, 0)
	for rows.Next() {
		var client models.Client
		err := rows.Scan(&client.Id, &client.Name, &client.Token, &client.Phone, &client.Expires, &client.Jid, &client.Connected)
		if err != nil {
			return nil, errors.Join(ErrGetAllClients, err)
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func (c *Client) Create(ctx context.Context, name, phone, token string) error {
	statement := "INSERT INTO clients(name, phone, token, expires, jid, connected) VALUES($1, $2, $3, $4, $5);"
	expires := time.Now().UTC()
	connected := false
	jid := ""
	_, err := c.conn.ExecContext(ctx, statement, name, phone, token, expires, jid, connected)
	if err != nil {
		return fmt.Errorf("error creating new client: %w", err)
	}
	return nil
}

func (c *Client) Extend(ctx context.Context, id, months, days int) error {
	statement := "SELECT expires FROM clients WHERE id = $1 LIMIT 1;"
	r := c.conn.QueryRowContext(ctx, statement, id)
	var expirationDate time.Time
	err := r.Scan(&expirationDate)
	if err != nil {
		return fmt.Errorf("error extending client %d expiration date: %w", id, err)
	}

	if expirationDate.UTC().Before(time.Now().UTC()) {
		expirationDate = time.Now().UTC()
	}
	expirationDate = expirationDate.AddDate(0, months, days)

	statement = "UPDATE clients SET expires = $1 WHERE id = $2;"
	_, err = c.conn.ExecContext(ctx, statement, expirationDate.UTC(), id)
	if err != nil {
		return fmt.Errorf("error extending client %d expiration date: %w", id, err)
	}
	return nil
}

func (c *Client) Cancel(ctx context.Context, id int) error {
	statement := "UPDATE clients SET expires = $1 WHERE id = $2;"
	_, err := c.conn.ExecContext(ctx, statement, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("error cancelling client subscription: %w", err)
	}
	return nil
}

func (c *Client) ConnectWh(ctx context.Context, id int, jid string) error {
	statement := "UPDATE clients SET connected = $1, jid = $2 WHERE id = $3;"
	_, err := c.conn.ExecContext(ctx, statement, true, jid, id)
	if err != nil {
		return fmt.Errorf("error updating client connected in db client: %d jid: %s: %w", id, jid, err)
	}
	return nil
}

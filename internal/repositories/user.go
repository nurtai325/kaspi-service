package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/nurtai325/kaspi-service/internal/config"
	"github.com/nurtai325/kaspi-service/internal/db"
	"github.com/nurtai325/kaspi-service/internal/models"
)

type User struct {
	conn *sql.DB
}

func NewUser() User {
	return User{
		conn: db.Conn(config.New()),
	}
}

func (u *User) All(ctx context.Context) ([]models.User, error) {
	statement := "SELECT id, name, phone, password FROM users ORDER BY id ASC;"
	rows, err := u.conn.QueryContext(ctx, statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]models.User, 0)
	for rows.Next() {
		nextUser := models.User{}
		err = rows.Scan(&nextUser.Id, &nextUser.Name, &nextUser.Phone, &nextUser.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, nextUser)
	}
	return users, nil
}

func (u *User) OneByPhone(ctx context.Context, phone string) (models.User, error) {
	statement := "SELECT id, name, phone, password FROM users WHERE phone = $1;"
	row := u.conn.QueryRowContext(ctx, statement, phone)
	var user models.User
	err := row.Scan(&user.Id, &user.Name, &user.Phone, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("user with phone %s not found: %w", phone, err)
	}
	return user, err
}

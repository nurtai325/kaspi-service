package models

import "time"

type Client struct {
	Id        int
	Name      string
	Token     string
	Phone     string
	Expires   time.Time
	Jid       string
	Connected bool
}

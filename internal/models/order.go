package models

type Order struct {
	Id            string
	Code          string
	ProductCode   string
	Customer      string
	CustomerPhone string
	Entries       []Entry
}

type QueuedOrder struct {
	ClientName  string
	ClientPhone string
	ClientJid   string
	Token       string
	Failed      int
	Order       Order
}

type Entry struct {
	ProductName string
	Quantity    int
	Price       float32
}

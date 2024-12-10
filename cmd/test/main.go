package main

import (
	"fmt"
	"time"

	"github.com/nurtai325/kaspi-service/internal/config"
	"github.com/nurtai325/kaspi-service/internal/db"
)

func main() {
	c := db.Conn(config.New())
	rows, err := c.Query("select name, expires from clients where expires < $1;", time.Now().UTC())
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var name string
		var expires time.Time
		err = rows.Scan(&name, &expires)
		if err != nil {
			panic(err)
		}
		fmt.Println(name, expires)
	}
}

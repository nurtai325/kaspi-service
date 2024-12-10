package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/nurtai325/kaspi-service/internal/config"
	"github.com/nurtai325/kaspi-service/internal/db"
	"github.com/nurtai325/kaspi-service/internal/handlers"
	"github.com/nurtai325/kaspi-service/internal/mailing"
	"github.com/nurtai325/kaspi-service/internal/repositories"
	"github.com/nurtai325/kaspi-service/internal/whatsapp"
)

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC | log.Llongfile)
	conf := config.New()
	clientRepo := repositories.NewClient()
	messenger := whatsapp.NewMessenger(db.Conn(conf))

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		mailing.Run(clientRepo, messenger)
		wg.Done()
	}()
	go func() {
		startServer(conf)
		wg.Done()
	}()
	wg.Wait()
}

func startServer(conf config.Config) {
	handlers.Register()
	err := handlers.ParseTemplates()
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Starting http server on port: %s", conf.PORT)
	err = http.ListenAndServe(":"+conf.PORT, nil)
	if err != nil {
		err = fmt.Errorf("error http server fail: %w", err)
		log.Panic(err)
	}
}

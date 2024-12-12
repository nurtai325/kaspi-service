package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nurtai325/kaspi-service/internal/config"
	"github.com/nurtai325/kaspi-service/internal/db"
	"github.com/nurtai325/kaspi-service/internal/handlers"
	"github.com/nurtai325/kaspi-service/internal/mailing"
	"github.com/nurtai325/kaspi-service/internal/repositories"
	"github.com/nurtai325/kaspi-service/internal/whatsapp"
)

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC | log.Llongfile)
	err := mailing.RecoverState()
	if err != nil {
		log.Panic(err)
	}
	conf := config.New()
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		clientRepo := repositories.NewClient()
		messenger := whatsapp.NewMessenger(db.Conn(conf))
		mailing.Run(clientRepo, messenger)
		log.Println("Starting mailing to clients")
		wg.Done()
	}()
	go func() {
		startServer(conf)
		wg.Done()
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("application interrupted. saving state")
		err := mailing.SaveState()
		if err != nil {
			log.Println(fmt.Errorf("error saving app state: %w", err))
		}
		log.Println("application state saved")
		os.Exit(0)
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
	err = http.ListenAndServeTLS(":"+conf.PORT, "./cert/certificate.crt", "./cert/private.key", nil)
	if err != nil {
		err = fmt.Errorf("error http server fail: %w", err)
		log.Panic(err)
	}
}

package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/nurtai325/kaspi-service/internal/models"
	"github.com/nurtai325/kaspi-service/internal/repositories"
)

var (
	clientsTempl       = "clients.html"
	clientsCreateTempl = "clients_create.html"
	columns            = []string{"Аты", "Номер", "Төленді дейін", "", "Whatsapp-қа қосылған", ""}
)

const (
	kzTimeZoneOffset = 5
)

type clientsTemplData struct {
	Clients []models.Client
	Columns []string
}

func HandleClients(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		clientRepo := repositories.NewClient()
		clients, err := clientRepo.All(r.Context())
		for i := range len(clients) {
			clients[i].Expires = clients[i].Expires.Add(time.Hour * kzTimeZoneOffset)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(newErr(r, err))
			return
		}
		Execute(w, r, clientsTempl, clientsTemplData{
			Clients: clients,
			Columns: columns,
		})
		return
	}
	http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
}

type clientsCreateTemplData struct {
	Error string
}

func HandleClientsCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		clientErr := r.URL.Query().Get(errorQuery)
		Execute(w, r, clientsCreateTempl, clientsCreateTemplData{
			Error: clientErr,
		})
		return
	}

	r.ParseForm()
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	token := r.FormValue("token")

	clientsRepo := repositories.NewClient()
	err := clientsRepo.Create(r.Context(), name, phone, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(newErr(r, err))
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

const (
	monthsUnit = "months"
	daysUnit   = "days"
)

func HandleClientExtend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
	}

	r.ParseForm()
	formId := r.Form.Get("id")
	clientId, err := strconv.Atoi(formId)
	if err != nil {
		err = fmt.Errorf("client id is not a number: %w", err)
		log.Println(newErr(r, err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	formDuration := r.Form.Get("duration")
	duration, err := strconv.Atoi(formDuration)
	if err != nil {
		http.Error(w, "Күн немесе ай сан болуы керек", http.StatusBadRequest)
		return
	}

	unit := r.Form.Get("unit")
	days := 0
	months := 0
	if unit == daysUnit {
		days += duration
	} else if unit == monthsUnit {
		months += duration
	} else {
		err = fmt.Errorf("unit is not months or days")
		log.Println(newErr(r, err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	clientRepo := repositories.NewClient()
	err = clientRepo.Extend(r.Context(), clientId, months, days)
	if err != nil {
		log.Println(newErr(r, err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func HandleClientsCancel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
	}

	r.ParseForm()
	formId := r.Form.Get("id")
	clientId, err := strconv.Atoi(formId)
	if err != nil {
		err = fmt.Errorf("client id is not a number: %w", err)
		log.Println(newErr(r, err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	clientRepo := repositories.NewClient()
	err = clientRepo.Cancel(r.Context(), clientId)
	if err != nil {
		log.Println(newErr(r, err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ystv/computing_site/team"
	"github.com/ystv/computing_site/templates"
)

type Web struct {
	mux  *mux.Router
	t    *templates.Templater
	team []team.Member
}

func main() {
	web := Web{mux: mux.NewRouter(), t: templates.New()}
	log.Println("Web loaded")
	var err error
	web.team, err = team.New()
	if err != nil {
		log.Println("failed to get team: %+v", err)
	}

	web.mux.HandleFunc("/", web.indexPage).Methods("GET")
	log.Println("YSTV Computing site: 0.0.0.0:7075")
	log.Fatal(http.ListenAndServe("0.0.0.0:7075", web.mux))
}

func (web *Web) indexPage(w http.ResponseWriter, r *http.Request) {
	params := templates.DashboardParams{
		Base: templates.BaseParams{
			SystemTime: time.Now(),
		},
		Team: web.team,
	}

	err := web.t.Dashboard(w, params)
	if err != nil {
		err = fmt.Errorf("failed to render dashboard: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

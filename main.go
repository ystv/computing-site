package main

import (
	"fmt"
	"github.com/ystv/computing_site/link"
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
	link link.Link
	team []team.Member
}

func main() {
	var err error

	web := Web{mux: mux.NewRouter(), t: templates.New()}
	log.Println("Web loaded")
	web.link, err = link.New()
	web.team, err = team.New()
	if err != nil {
		log.Printf("failed to get team: %+v\n", err)
	}

	web.mux.HandleFunc("/", web.indexPage).Methods("GET")
	web.mux.HandleFunc("/ystv.ico", web.faviconHandler)
	web.mux.HandleFunc("/stylesheet.css", web.cssHandler)
	log.Println("YSTV Computing site: 0.0.0.0:7075")
	log.Fatal(http.ListenAndServe("0.0.0.0:7075", web.mux))
}

func (web *Web) indexPage(w http.ResponseWriter, _ *http.Request) {
	params := templates.DashboardParams{
		Base: templates.BaseParams{
			SystemTime: time.Now(),
		},
		Link: web.link,
		Team: web.team,
	}

	err := web.t.Dashboard(w, params)
	if err != nil {
		err = fmt.Errorf("failed to render dashboard: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (web *Web) faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/ystv.ico")
}

func (web *Web) cssHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/stylesheet.css")
}

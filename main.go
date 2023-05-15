package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ystv/computing_site/link"
	"github.com/ystv/computing_site/team"
	"github.com/ystv/computing_site/templates"
	"log"
	"net/http"
)

type Web struct {
	mux  *mux.Router
	t    *templates.Templater
	link *link.Link
	team *[]team.Member
}

func main() {
	var err error

	web := Web{
		mux: mux.NewRouter(),
	}
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
	params := &templates.DashboardParams{
		Link: web.link,
		Team: web.team,
	}

	err := web.t.RenderTemplate(w, params, "dashboard.tmpl")
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

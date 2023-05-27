package main

import (
	"embed"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ystv/computing_site/link"
	"github.com/ystv/computing_site/team"
	"github.com/ystv/computing_site/templates"
	"io/fs"
	"log"
	"net/http"
)

//go:embed public/*
var embeddedFiles embed.FS

var (
	Commit  = "unknown"
	Version = "unknown"
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
	web.link, err = link.New()
	web.team, err = team.New()
	if err != nil {
		log.Printf("failed to get team: %+v\n", err)
	}

	assetHandler := http.FileServer(getFileSystem())

	addr := "0.0.0.0:7075"

	web.mux.HandleFunc("/", web.indexPage).Methods("GET")
	//web.mux.HandleFunc("/ystv.ico", web.faviconHandler)
	//web.mux.HandleFunc("/stylesheet.css", web.cssHandler)
	web.mux.PathPrefix("/public/").Handler(http.StripPrefix("/public/", assetHandler))
	log.Printf("YSTV Computing site: %s, commit: %s, version: %s\n", addr, Commit, Version)
	log.Fatal(http.ListenAndServe(addr, web.mux))
}

func (web *Web) indexPage(w http.ResponseWriter, _ *http.Request) {
	params := &templates.DashboardParams{
		Link:    web.link,
		Team:    web.team,
		Commit:  Commit,
		Version: Version,
	}

	err := web.t.RenderTemplate(w, params, "dashboard.tmpl")
	if err != nil {
		err = fmt.Errorf("failed to render dashboard: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//func (web *Web) faviconHandler(w http.ResponseWriter, r *http.Request) {
//	http.ServeFile(w, r, "public/ystv.ico")
//}
//
//func (web *Web) cssHandler(w http.ResponseWriter, r *http.Request) {
//	http.ServeFile(w, r, "public/stylesheet.css")
//}

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(embeddedFiles, "public")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

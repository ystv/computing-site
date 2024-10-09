package main

import (
	"embed"
	"log"

	"github.com/ystv/computing_site/link"
	"github.com/ystv/computing_site/team"
)

//go:embed public/*
var embeddedFiles embed.FS

var (
	Commit  = "unknown"
	Version = "unknown"

	cert = ""
	key  = ""
)

func main() {
	if cert == "" || key == "" {
		log.Fatalf("missing required cert and key")
	}

	link1, err := link.New()
	if err != nil {
		log.Printf("failed to get link: %+v\n", err)
	}
	team1, err := team.New()
	if err != nil {
		log.Printf("failed to get team: %+v\n", err)
	}

	addr := "0.0.0.0:7075"

	router := NewRouter(&RouterConf{
		Address: addr,
		Link:    link1,
		Team:    team1,
		Commit:  Commit,
		Version: Version,
		cert:    []byte(cert),
		key:     []byte(key),
	})

	log.Printf("YSTV Computing site: %s, commit: %s, version: %s\n", addr, Commit, Version)

	err = router.Start()
	if err != nil {
		log.Fatalf("The web server couldn't be started!\n\n%s\n\nExiting!", err)
	}
}

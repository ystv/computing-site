package link

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type (
	Base struct {
		URL   string `json:"url"`
		Name  string `json:"name"`
		Extra string `json:"extra,omitempty"`
	}

	Link struct {
		Docs    []Docs    `json:"docs"`
		Comp    []Comp    `json:"comp"`
		Uni     []Uni     `json:"uni"`
		Staging []Staging `json:"staging"`
	}

	Docs struct {
		Base
	}

	Comp struct {
		Base
	}

	Uni struct {
		Base
	}

	Staging struct {
		Base
	}
)

func New() (*Link, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading .env file: %+v, continuing", err)
	}

	link := os.Getenv("LINK_JSON")
	if len(link) == 0 {
		return nil, fmt.Errorf("LINK_JSON environment variable cannot be found")
	}

	link = strings.ReplaceAll(link, "~", "\"")

	var data *Link
	err = json.Unmarshal([]byte(link), &data)
	if err != nil {
		return &Link{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return data, nil
}

package link

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed link.json
var linkJson embed.FS

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
	link, err := linkJson.ReadFile("link.json")
	if err != nil {
		fmt.Printf("error loading link.json file: %s", err)
	}

	var data *Link
	err = json.Unmarshal(link, &data)
	if err != nil {
		return &Link{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return data, nil
}

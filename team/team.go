package team

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed team.json
var teamJson embed.FS

type (
	Member struct {
		Name string `json:"name"`
		Role string `json:"role,omitempty"`
	}
)

func New() (*[]Member, error) {
	team, err := teamJson.ReadFile("team.json")
	if err != nil {
		fmt.Printf("error loading team.json file: %s", err)
	}

	var data *[]Member
	err = json.Unmarshal(team, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return data, nil
}

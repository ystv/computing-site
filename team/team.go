package team

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type (
	Member struct {
		Name string `json:"name"`
		Role string `json:"role,omitempty"`
	}
)

func New() (*[]Member, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	team := os.Getenv("TEAM_JSON")
	if len(team) == 0 {
		return nil, fmt.Errorf("TEAM_JSON environment variable cannot be found")
	}

	var data *[]Member
	err = json.Unmarshal([]byte(team), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return data, nil
}

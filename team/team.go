package team

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

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
		log.Printf("error loading .env file: %+v, continuing", err)
	}

	team := os.Getenv("TEAM_JSON")
	if len(team) == 0 {
		return nil, fmt.Errorf("TEAM_JSON environment variable cannot be found")
	}

	team = strings.ReplaceAll(team, "~", "\"")

	var data *[]Member
	err = json.Unmarshal([]byte(team), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return data, nil
}

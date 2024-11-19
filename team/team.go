package team

import (
	"encoding/json"
	"fmt"
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
		fmt.Printf("error loading .env file: %+v\n", err)
	}

	team := os.Getenv("TEAM_JSON")
	if len(team) == 0 {
		fmt.Println("TEAM_JSON environment variable not set")
		//return nil, fmt.Errorf("TEAM_JSON environment variable cannot be found")
	}

	var sec map[string]interface{}
	if err = json.Unmarshal([]byte(team), &sec); err != nil {
		fmt.Printf("error unmarshalling secret: %+v\n", err)
	}

	fmt.Println(sec)

	fmt.Println(team)
	team = strings.ReplaceAll(team, "'", "\"")
	fmt.Println(team)

	var data *[]Member
	err = json.Unmarshal([]byte(team), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return data, nil
}

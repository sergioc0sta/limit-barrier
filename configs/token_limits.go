package configs

import (
	"encoding/json"
	"errors"
	"os"
)

type TokenLimit struct {
	Token string `json:"token"`
	Limit int    `json:"limit"`
}

func LoadTokenLimits(path string) (map[string]int, error) {
	if path == "" {
		return nil, errors.New("path is required")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var items []TokenLimit
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	limits := make(map[string]int, len(items))

	return limits, nil
}

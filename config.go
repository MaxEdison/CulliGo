package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url"`
	Apikey   string `json:"api_key"`
}

type Meal struct {
	Name string `json:"name"`
}

type meals struct {
	Day       string `json:"day"`
	Breakfast []Meal `json:"breakfast"`
	Lunch     []Meal `json:"lunch"`
	Dinner    []Meal `json:"dinner"`
}

type food struct {
	WeekFlag string  `json:weekflag`
	Week     string  `json:"week"`
	Meals    []meals `json:"meals"`
}

func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

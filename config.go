package main

import (
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

	cfg := Config{
		Username: os.Getenv("username"),
		Password: os.Getenv("password"),
		URL:      os.Getenv("url"),
		Apikey:   os.Getenv("api_key"),
	}
	return &cfg, nil
}

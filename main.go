package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	solver "github.com/MaxEdison/CulliGo/Captcha_Solver"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/utils"
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
	Week  string  `json:"week"`
	Meals []meals `json:"meals"`
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

func send_data(page *rod.Page, cfg Config) error {

	page.MustElement("#txtUsernamePlain").MustInput(cfg.Username)
	page.MustElement("#txtPasswordPlain").MustInput(cfg.Password)

	img := page.MustElement("#Img1")
	_ = utils.OutputFile("Captchas/captcha.png", img.MustResource())

	code, err := solver.Solver("Captchas/captcha.png", cfg.Apikey)
	if err != nil {
		return fmt.Errorf("captcha solver failed: %w", err)
	}

	page.MustElement("#txtCaptcha").MustInput(code)
	page.MustElement("#btnEncript").MustClick()

	return nil
}

func try_login(page *rod.Page, cfg Config) error {

	if err := send_data(page, cfg); err != nil {
		return err
	}

	err := rod.Try(func() {

		text := page.Timeout(5 * time.Second).MustElement("#lblLoginError").MustText()
		if text != "" {
			panic("login-error")
		}
	})

	if errors.Is(err, context.DeadlineExceeded) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	return fmt.Errorf("invalid login or captcha")
}

func get_table_data(page *rod.Page, table_id string) ([]struct {
	day   string
	meals []Meal
}, error) {
	table, _ := page.Element("#" + table_id)

	rows, _ := table.Elements("tr:not(.HeaderStyle)")

	var data []struct {
		day   string
		meals []Meal
	}

	for _, row := range rows {
		cols, _ := row.Elements("td")

		if len(cols) < 4 {
			continue
		}

		day := strings.TrimSpace(cols[1].MustText())
		var meals []Meal

		select_elmnt, err := cols[3].Element("select")

		if err == nil {
			options, err := select_elmnt.Elements("option")
			if err == nil && len(options) > 0 {
				for _, option := range options {
					value := strings.TrimSpace(option.MustText())
					if value == "" || value == "-" {
						continue
					}

					parts := strings.Split(value, "ریال")
					if len(parts) > 0 {
						continue
					}

					name := strings.TrimSpace(parts[0])
					name = strings.ReplaceAll(name, "@", "[رستوران]")
					name = strings.ReplaceAll(name, "*ر", "[کافه کتاب]")

					last_space := strings.LastIndex(name, " ")
					name = strings.TrimSpace(name[:last_space])

					meals = append(meals, Meal{Name: name})
				}
			}
		}

		data = append(data, struct {
			day   string
			meals []Meal
		}{day: day, meals: meals})
	}

	return data, nil
}

func scraper(page *rod.Page, week string) (food, error) {

	page.MustWaitLoad()

	if week == "next" {
		page.MustElement("#lnkNextWeek").MustClick()
		page.MustWaitLoad()
		time.Sleep(3 * time.Second)
	}

	breakfast, err := get_table_data(page, "cphMain_grdReservationBreakfast")
	if err != nil {
		return food{}, err
	}

	lunch, err := get_table_data(page, "cphMain_grdReservationLunch")
	if err != nil {
		return food{}, err
	}

	dinner, err := get_table_data(page, "cphMain_grdReservationDinner")
	if err != nil {
		return food{}, err
	}

	days := make(map[string]*meals)
	for _, data := range []struct {
		meal_time string
		data      []struct {
			day   string
			meals []Meal
		}
	}{
		{"breakfast", breakfast},
		{"lunch", lunch},
		{"dinner", dinner},
	} {
		for _, entry := range data.data {

			if _, exists := days[entry.day]; !exists {
				days[entry.day] = &meals{Day: entry.day}
			}

			switch data.meal_time {
			case "breakfast":
				days[entry.day].Breakfast = entry.meals

			case "lunch":
				days[entry.day].Lunch = entry.meals

			case "dinner":
				days[entry.day].Dinner = entry.meals
			}
		}
	}

	var all_meals []meals
	for _, meals := range days {
		if meals.Breakfast == nil {
			meals.Breakfast = []Meal{}
		}

		if meals.Lunch == nil {
			meals.Lunch = []Meal{}
		}

		if meals.Dinner == nil {
			meals.Dinner = []Meal{}
		}

		all_meals = append(all_meals, *meals)
	}
}

func main() {

	cfg, err := loadConfig("Config/config.json")
	if err != nil {
		fmt.Println("[ERROR] Can't load config:", err)
		return
	}

	browser := rod.New().NoDefaultDevice().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(cfg.URL)
	page.MustWaitLoad()

	var Login_Err error

	for i := 1; i <= 3; i++ {
		Login_Err = try_login(page, *cfg)
		fmt.Println("[DEBUG] Login Attempt", i, ":", Login_Err)
		if Login_Err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if Login_Err != nil {
		fmt.Println("[ERROR] All login attempts failed:", Login_Err)
		return
	}

	page = browser.MustPage(cfg.URL + "Reservation/Reservation.aspx")

	time.Sleep(time.Hour)
}

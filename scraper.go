package main

import (
	"strings"
	"time"

	"github.com/go-rod/rod"
)

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
					if len(parts) < 2 {
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

	return food{Week: week, Meals: all_meals}, nil
}

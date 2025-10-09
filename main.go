package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-rod/rod"
)

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

	week := "this" // this | next
	food, err := scraper(page, week)

	if err != nil {
		fmt.Println("[ERROR] Failed to scrape food data:", err)
		return
	}

	final_json, err := json.MarshalIndent(food, "", "  ")

	if err != nil {
		fmt.Println("[ERROR] Failed to write JSON to file:", err)
		return
	}

	if _, err := os.Stat("data"); os.IsNotExist(err) {
		err = os.Mkdir("data", 0755)
		if err != nil {
			fmt.Println("[ERROR] Failed to create data directory:", err)
			return
		}
	}

	err = os.WriteFile("data/data.json", final_json, 0644)

	if err != nil {
		fmt.Println("[ERROR] Failed to write JSON to file:", err)
		return
	}

	fmt.Println("[INFO] DONE !")

	// time.Sleep(time.Hour)
}

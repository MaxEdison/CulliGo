package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
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

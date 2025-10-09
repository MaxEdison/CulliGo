package main

import (
	"fmt"

	solver "github.com/MaxEdison/CulliGo/Captcha_Solver"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/utils"
)

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

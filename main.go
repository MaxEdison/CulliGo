package main

import (
	solver "CulliGo/Captcha_Solver"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/utils"
)

func send_data(username, password string, page *rod.Page) error {

	page.MustElement("#txtUsernamePlain").MustInput(username)
	page.MustElement("#txtPasswordPlain").MustInput(password)

	img := page.MustElement("#Img1")
	_ = utils.OutputFile("Captchas/captcha.png", img.MustResource())

	code, err := solver.Solver("Captchas/captcha.png")
	if err != nil {
		return fmt.Errorf("captcha solver failed: %w", err)
	}

	page.MustElement("#txtCaptcha").MustInput(code)
	page.MustElement("#btnEncript").MustClick()

	return nil
}

func try_login(username, password string, page *rod.Page) error {

	if err := send_data(username, password, page); err != nil {
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
	browser := rod.New().NoDefaultDevice().MustConnect()
	page := browser.MustPage("FOOD RESERVATION WEBSITE URL")

	var Login_Err error

	for i := 1; i <= 3; i++ {
		Login_Err = try_login("USERNAME", "PASSWORD", page)

		if Login_Err == nil {
			break
		}

		time.Sleep(2 * time.Second)
	}

	if Login_Err != nil {
		return
	}

	page.MustWaitLoad()
	time.Sleep(time.Hour)
}

package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-rod/rod"
)

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

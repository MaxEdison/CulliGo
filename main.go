package main

import (
	"github.com/go-rod/rod"
)

func main() {
	browser := rod.New().NoDefaultDevice().MustConnect()
	page := browser.MustPage("FOOD RESERVATION WEBSITE URL")
	page.MustWaitLoad()
}

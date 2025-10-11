package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	cfg, err := loadConfig("Config/config.json")
	if err != nil {
		log.Fatalf("[ERROR] Can't load config: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName:               "CulliGo API",
		ServerHeader:          "CulliGo",
		DisableStartupMessage: false,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	browser := rod.New().NoDefaultDevice().MustConnect()
	defer browser.MustClose()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now(),
		})
	})

	app.Get("/scrape/:week", func(c *fiber.Ctx) error {
		week := c.Params("week", "this")

		page := browser.MustPage(cfg.URL)
		page.MustWaitLoad()

		var loginErr error
		for i := 1; i <= 3; i++ {
			loginErr = try_login(page, *cfg)
			fmt.Println("[DEBUG] Login Attempt", i, ":", loginErr)
			if loginErr == nil {
				break
			}
			time.Sleep(2 * time.Second)
		}

		if loginErr != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "login failed",
				"details": loginErr.Error(),
			})
		}

		page = browser.MustPage(cfg.URL + "Reservation/Reservation.aspx")

		food, err := scraper(page, week)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "scraping failed",
				"details": err.Error(),
			})
		}

		data, err := json.MarshalIndent(food, "", "  ")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "failed to marshal JSON",
				"details": err.Error(),
			})
		}

		if _, err := os.Stat("data"); os.IsNotExist(err) {
			err = os.Mkdir("data", 0755)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "failed to create data directory",
					"details": err.Error(),
				})
			}
		}

		err = os.WriteFile("data/data.json", data, 0644)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "failed to write data file",
				"details": err.Error(),
			})
		}

		return c.JSON(food)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("[INFO] CulliGo API running on port %s", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("[ERROR] Server failed: %v", err)
	}
}

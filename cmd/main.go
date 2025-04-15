package main

import (
	"log"

	"github.com/TnAhonk12/TechnicalTestIhsanSolusiBE/config"
	"github.com/TnAhonk12/TechnicalTestIhsanSolusiBE/db"
	"github.com/TnAhonk12/TechnicalTestIhsanSolusiBE/model"
	"github.com/TnAhonk12/TechnicalTestIhsanSolusiBE/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.LoadConfig()
	db.Init(cfg)
	db.DB.AutoMigrate(&model.Nasabah{}, &model.Transaksi{})

	app := fiber.New()

	app.Use(logger.New())
	app.Use(recover.New())

	router.Setup(app)

	log.Printf("[INFO] Server running at port %s", cfg.AppPort)
	log.Fatal(app.Listen(cfg.AppHost + ":" + cfg.AppPort))
}

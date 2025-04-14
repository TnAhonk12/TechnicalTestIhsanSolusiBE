package main

import (
	"log"

	"github.com/TnAhonk12/TechnicalTestIhsanSolusiBE/config"
	"github.com/TnAhonk12/TechnicalTestIhsanSolusiBE/db"
	"github.com/TnAhonk12/TechnicalTestIhsanSolusiBE/model"
	"github.com/TnAhonk12/TechnicalTestIhsanSolusiBE/router"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.LoadConfig()
	db.Init(cfg)
	db.DB.AutoMigrate(&model.Nasabah{})

	app := fiber.New()
	router.Setup(app)

	log.Printf("[INFO] Server running at port %s", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}

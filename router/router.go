package router

import (
	"github.com/TnAhonk12/TechnicalTestIhsanSolusiBE/handler"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/nasabah", handler.ListNasabah)
	app.Post("/daftar", handler.Daftar)
	app.Post("/tabung", handler.Tabung)
	app.Post("/tarik", handler.Tarik)
	app.Get("/saldo/:no_rekening", handler.CekSaldo)
	app.Get("/transaksi/:no_rekening", handler.ListTransaksiByRekening)
}

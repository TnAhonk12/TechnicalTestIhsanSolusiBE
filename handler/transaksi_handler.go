package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/TnAhonk12/TechnicalTestIhsanSolusiBE/db"
	"github.com/TnAhonk12/TechnicalTestIhsanSolusiBE/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Tabung(c *fiber.Ctx) error {
	var payload struct {
		NoRekening string `json:"no_rekening"`
		Nominal    int64  `json:"nominal"`
	}
	if err := c.BodyParser(&payload); err != nil {
		log.Printf("[ERROR] Gagal parsing data tabung: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"remark": "Input tidak valid"})
	}

	log.Printf("[INFO] Permintaan tabung: Rekening=%s, Nominal=%d", payload.NoRekening, payload.Nominal)

	var nasabah model.Nasabah
	err := db.DB.Where("no_rekening = ?", payload.NoRekening).First(&nasabah).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[WARNING] Rekening tidak ditemukan saat menabung: %s", payload.NoRekening)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"remark": "Rekening tidak ditemukan"})
	} else if err != nil {
		log.Printf("[CRITICAL] Gagal query rekening saat tabung: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"remark": "Gagal mengambil data rekening"})
	}

	nasabah.Saldo += payload.Nominal
	if err := db.DB.Save(&nasabah).Error; err != nil {
		log.Printf("[CRITICAL] Gagal menyimpan saldo baru ke database: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"remark": "Gagal menyimpan saldo"})
	}

	transaksi := model.Transaksi{
		NoRekening: nasabah.NoRekening,
		Tipe:       "tabung",
		Nominal:    payload.Nominal,
		Saldo:      nasabah.Saldo,
	}
	db.DB.Create(&transaksi)

	log.Printf("[INFO] Tabung berhasil: Rekening=%s, Saldo Sekarang=%d", payload.NoRekening, nasabah.Saldo)
	return c.JSON(fiber.Map{"saldo": nasabah.Saldo})
}

func Tarik(c *fiber.Ctx) error {
	var payload struct {
		NoRekening string `json:"no_rekening"`
		Nominal    int64  `json:"nominal"`
	}
	if err := c.BodyParser(&payload); err != nil {
		log.Printf("[ERROR] Gagal parsing data tarik: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"remark": "Input tidak valid"})
	}

	log.Printf("[INFO] Permintaan tarik: Rekening=%s, Nominal=%d", payload.NoRekening, payload.Nominal)

	var nasabah model.Nasabah
	err := db.DB.Where("no_rekening = ?", payload.NoRekening).First(&nasabah).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[WARNING] Rekening tidak ditemukan saat tarik: %s", payload.NoRekening)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"remark": "Rekening tidak ditemukan"})
	} else if err != nil {
		log.Printf("[CRITICAL] Gagal query rekening saat tarik: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"remark": "Gagal mengambil data rekening"})
	}

	if nasabah.Saldo < payload.Nominal {
		log.Printf("[WARNING] Gagal tarik: Saldo tidak cukup. Rekening=%s, Saldo=%d, Permintaan=%d", payload.NoRekening, nasabah.Saldo, payload.Nominal)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"remark": "Saldo tidak mencukupi"})
	}

	nasabah.Saldo -= payload.Nominal
	if err := db.DB.Save(&nasabah).Error; err != nil {
		log.Printf("[CRITICAL] Gagal memperbarui saldo setelah tarik: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"remark": "Gagal memperbarui saldo"})
	}

	transaksi := model.Transaksi{
		NoRekening: nasabah.NoRekening,
		Tipe:       "tarik",
		Nominal:    payload.Nominal,
		Saldo:      nasabah.Saldo,
	}
	db.DB.Create(&transaksi)

	log.Printf("[INFO] Tarik berhasil: Rekening=%s, Saldo Sekarang=%d", payload.NoRekening, nasabah.Saldo)
	return c.JSON(fiber.Map{"saldo": nasabah.Saldo})
}

func ListTransaksiByRekening(c *fiber.Ctx) error {
	rekening := c.Params("no_rekening")
	var transaksi []model.Transaksi

	if err := db.DB.Where("no_rekening = ?", rekening).Find(&transaksi).Error; err != nil {
		log.Printf("[CRITICAL] Gagal ambil transaksi: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"remark": "Gagal mengambil transaksi"})
	}

	log.Printf("[INFO] Ditemukan %d transaksi untuk rekening %s", len(transaksi), rekening)
	return c.JSON(transaksi)
}

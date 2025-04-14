package handler

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/TnAhonk12/TechnicalTestIhsanSolusiBE/db"
	"github.com/TnAhonk12/TechnicalTestIhsanSolusiBE/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func generateRekening() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(999999999))
}

func Daftar(c *fiber.Ctx) error {
	var nasabah model.Nasabah
	if err := c.BodyParser(&nasabah); err != nil {
		log.Printf("[ERROR] Gagal parsing request body: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"remark": "Input tidak valid"})
	}

	log.Printf("[INFO] Menerima pendaftaran nasabah: Nama=%s, NIK=%s, NoHP=%s", nasabah.Nama, nasabah.NIK, nasabah.NoHP)

	var existing model.Nasabah
	err := db.DB.Where("nik = ? OR no_hp = ?", nasabah.NIK, nasabah.NoHP).First(&existing).Error
	if err == nil {
		log.Printf("[WARNING] Duplikasi data: NIK atau NoHP sudah digunakan. NIK=%s, NoHP=%s", nasabah.NIK, nasabah.NoHP)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"remark": "NIK atau No HP sudah digunakan"})
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[CRITICAL] Gagal query ke database: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"remark": "Terjadi kesalahan saat verifikasi data"})
	}

	nasabah.NoRekening = generateRekening()
	nasabah.Saldo = 0
	if err := db.DB.Create(&nasabah).Error; err != nil {
		log.Printf("[CRITICAL] Gagal menyimpan data nasabah baru ke database: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"remark": "Gagal menyimpan data nasabah"})
	}

	transaksi := model.Transaksi{
		NoRekening: nasabah.NoRekening,
		Tipe:       "setor awal",
		Nominal:    nasabah.Saldo,
	}
	db.DB.Create(&transaksi)

	log.Printf("[INFO] Nasabah berhasil terdaftar: Nama=%s, NoRekening=%s", nasabah.Nama, nasabah.NoRekening)
	return c.JSON(fiber.Map{"no_rekening": nasabah.NoRekening})
}

func ListNasabah(c *fiber.Ctx) error {
	log.Println("[INFO] Permintaan daftar semua nasabah")

	var nasabahList []model.Nasabah
	err := db.DB.Select("nama", "no_rekening", "saldo", "nik", "no_hp").Find(&nasabahList).Error
	if err != nil {
		log.Printf("[CRITICAL] Gagal mengambil data nasabah: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"remark": "Gagal mengambil data nasabah"})
	}

	log.Printf("[INFO] Ditemukan %d nasabah", len(nasabahList))
	return c.JSON(nasabahList)
}

func CekSaldo(c *fiber.Ctx) error {
	rekening := c.Params("no_rekening")
	log.Printf("[INFO] Permintaan cek saldo: Rekening=%s", rekening)

	var nasabah model.Nasabah
	err := db.DB.Where("no_rekening = ?", rekening).First(&nasabah).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[WARNING] Rekening tidak ditemukan saat cek saldo: %s", rekening)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"remark": "Rekening tidak ditemukan"})
	} else if err != nil {
		log.Printf("[CRITICAL] Gagal query rekening saat cek saldo: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"remark": "Gagal mengambil data rekening"})
	}

	log.Printf("[INFO] Cek saldo berhasil: Rekening=%s, Saldo=%d", rekening, nasabah.Saldo)
	return c.JSON(fiber.Map{"saldo": nasabah.Saldo})
}

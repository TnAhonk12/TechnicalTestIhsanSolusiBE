package model

import "gorm.io/gorm"

type Nasabah struct {
	gorm.Model
	Nama       string `json:"nama"`
	NIK        string `json:"nik" gorm:"uniqueIndex"`
	NoHP       string `json:"no_hp" gorm:"uniqueIndex"`
	NoRekening string `json:"no_rekening" gorm:"uniqueIndex"`
	Saldo      int64  `json:"saldo"`
}

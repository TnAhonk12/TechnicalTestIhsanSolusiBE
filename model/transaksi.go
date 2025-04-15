package model

import "gorm.io/gorm"

type Transaksi struct {
	gorm.Model
	NoRekening string `json:"no_rekening"`
	Tipe       string `json:"tipe"`
	Nominal    int64  `json:"nominal"`
	Saldo      int64  `json:"saldo"`
}

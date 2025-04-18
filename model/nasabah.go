package model

type Nasabah struct {
	NoRekening string      `json:"no_rekening" gorm:"primaryKey"`
	Nama       string      `json:"nama"`
	NIK        string      `json:"nik" gorm:"uniqueIndex"`
	NoHP       string      `json:"no_hp" gorm:"uniqueIndex"`
	Saldo      int64       `json:"saldo"`
	Transaksi  []Transaksi `gorm:"foreignKey:NoRekening;references:NoRekening"`
}

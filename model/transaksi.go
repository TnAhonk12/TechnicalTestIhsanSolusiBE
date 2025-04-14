package model

type Transaksi struct {
	NoRekening string `json:"no_rekening"`
	Tipe       string `json:"tipe"`
	Nominal    int64  `json:"nominal"`
	Saldo      int64  `json:"saldo"`
}

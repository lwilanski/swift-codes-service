package models

type SwiftCode struct {
	SwiftCode     string `gorm:"primaryKey;size:11"`
	BankName      string
	Address       string
	CountryISO2   string `gorm:"size:2;index"`
	CountryName   string
	IsHeadquarter bool
}

func (s SwiftCode) Root() string { return s.SwiftCode[:8] }

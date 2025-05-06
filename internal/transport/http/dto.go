package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lwilanski/swift-codes-service/internal/models"
)

type branchDTO struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

func toBranchDTO(m models.SwiftCode) branchDTO {
	return branchDTO{
		Address:       m.Address,
		BankName:      m.BankName,
		CountryISO2:   m.CountryISO2,
		IsHeadquarter: m.IsHeadquarter,
		SwiftCode:     m.SwiftCode,
	}
}

func HQResponse(hq models.SwiftCode, branches []models.SwiftCode) gin.H {
	b := make([]branchDTO, 0, len(branches))
	for _, br := range branches {
		b = append(b, toBranchDTO(br))
	}
	return gin.H{
		"address":       hq.Address,
		"bankName":      hq.BankName,
		"countryISO2":   hq.CountryISO2,
		"countryName":   hq.CountryName,
		"isHeadquarter": true,
		"swiftCode":     hq.SwiftCode,
		"branches":      b,
	}
}

func BranchResponse(br models.SwiftCode) gin.H {
	return gin.H{
		"address":       br.Address,
		"bankName":      br.BankName,
		"countryISO2":   br.CountryISO2,
		"countryName":   br.CountryName,
		"isHeadquarter": false,
		"swiftCode":     br.SwiftCode,
	}
}

func CountryResponse(list []models.SwiftCode) gin.H {
	if len(list) == 0 {
		return gin.H{}
	}
	countryName := list[0].CountryName
	iso2 := list[0].CountryISO2
	swifts := make([]branchDTO, 0, len(list))
	for _, s := range list {
		swifts = append(swifts, toBranchDTO(s))
	}
	return gin.H{
		"countryISO2": iso2,
		"countryName": countryName,
		"swiftCodes":  swifts,
	}
}

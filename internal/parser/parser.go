package parser

import (
	"strings"

	"github.com/lwilanski/swift-codes-service/internal/models"
	"github.com/xuri/excelize/v2"
)

func ParseExcel(path string) ([]models.SwiftCode, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return nil, err
	}

	var out []models.SwiftCode
	for i, r := range rows {
		if i == 0 {
			continue
		}
		iso2 := strings.ToUpper(r[0])
		code := strings.ToUpper(r[1])
		bank := r[3]
		addr := r[4]
		countryName := strings.ToUpper(r[6])

		isHQ := strings.HasSuffix(code, "XXX")
		out = append(out, models.SwiftCode{
			SwiftCode:     code,
			BankName:      bank,
			Address:       addr,
			CountryISO2:   iso2,
			CountryName:   countryName,
			IsHeadquarter: isHQ,
		})
	}
	return out, nil
}

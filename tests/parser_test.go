package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

func TestParseExcel(t *testing.T) {
	f := excelize.NewFile()
	s := f.GetSheetName(0)
	_ = f.SetSheetRow(s, "A1",
		&[]interface{}{"COUNTRY ISO2 CODE", "SWIFT CODE", "CODE TYPE", "NAME", "ADDRESS", "TOWN", "COUNTRY NAME"})
	_ = f.SetSheetRow(s, "A2",
		&[]interface{}{"PL", "FOOOPLP1XXX", "BIC11", "Foo Bank", "Main St", "Warsaw", "POLAND"})
	_ = f.SaveAs("test.xlsx")
	defer os.Remove("test.xlsx")

	out, err := ParseExcel("test.xlsx")
	require.NoError(t, err)
	require.Len(t, out, 1)
	got := out[0]
	require.True(t, got.IsHeadquarter)
	require.Equal(t, "PL", got.CountryISO2)
	require.Equal(t, "FOOOPLP1XXX", got.SwiftCode)
}

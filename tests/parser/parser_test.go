//go:build unit
// +build unit

package parser_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"

	"github.com/lwilanski/swift-codes-service/internal/parser"
)

func TestParseExcel(t *testing.T) {
	f := excelize.NewFile()
	s := f.GetSheetName(0)
	_ = f.SetSheetRow(s, "A1", &[]any{
		"COUNTRY ISO2 CODE", "SWIFT CODE", "CODE TYPE",
		"NAME", "ADDRESS", "TOWN", "COUNTRY NAME",
	})
	_ = f.SetSheetRow(s, "A2", &[]any{
		"PL", "FOOOPLP1XXX", "BIC11",
		"Foo Bank", "Main St", "Warsaw", "POLAND",
	})
	_ = f.SaveAs("test.xlsx")
	defer os.Remove("test.xlsx")

	out, err := parser.ParseExcel("test.xlsx")
	require.NoError(t, err)
	require.Len(t, out, 1)
	require.True(t, out[0].IsHeadquarter)
	require.Equal(t, "FOOOPLP1XXX", out[0].SwiftCode)
}

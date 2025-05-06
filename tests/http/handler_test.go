//go:build unit
// +build unit

package http_test

import (
	"encoding/json"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lwilanski/swift-codes-service/internal/models"
	repo "github.com/lwilanski/swift-codes-service/internal/repository"
	th "github.com/lwilanski/swift-codes-service/internal/transport/http"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testServer(t *testing.T) *httptest.Server {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	_ = db.AutoMigrate(&models.SwiftCode{})

	r := repo.New(db)
	_ = r.UpsertMany([]models.SwiftCode{
		{SwiftCode: "BANKPLPWXXX", CountryISO2: "PL", CountryName: "POLAND", BankName: "HQ", IsHeadquarter: true},
	})

	return httptest.NewServer(th.Router(th.New(r)))
}

func TestGetByCode(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp, err := stdhttp.Get(srv.URL + "/v1/swift-codes/BANKPLPWXXX")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, 200, resp.StatusCode)

	var body map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.True(t, body["isHeadquarter"].(bool))
}

func TestCreateAndDelete(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	j := `{
		"address":"Street 1",
		"bankName":"Demo",
		"countryISO2":"PL",
		"countryName":"POLAND",
		"isHeadquarter":false,
		"swiftCode":"TESTPLPWXXX"
	}`
	resp, err := stdhttp.Post(srv.URL+"/v1/swift-codes", "application/json", strings.NewReader(j))
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)

	req, _ := stdhttp.NewRequest(stdhttp.MethodDelete, srv.URL+"/v1/swift-codes/TESTPLPWXXX", nil)
	resp, _ = stdhttp.DefaultClient.Do(req)
	require.Equal(t, 200, resp.StatusCode)
}

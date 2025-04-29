package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/your-github-name/swift-codes-service/internal/models"
	"github.com/your-github-name/swift-codes-service/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testServer(t *testing.T) *httptest.Server {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	_ = db.AutoMigrate(&models.SwiftCode{})

	repo := repository.New(db)
	_ = repo.UpsertMany([]models.SwiftCode{
		{SwiftCode: "BANKPLPWXXX", CountryISO2: "PL", CountryName: "POLAND", BankName: "HQ", IsHeadquarter: true},
	})

	return httptest.NewServer(Router(New(repo)))
}

func TestGetByCode(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/v1/swift-codes/BANKPLPWXXX")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, 200, resp.StatusCode)

	var body map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.Equal(t, true, body["isHeadquarter"])
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
	resp, err := http.Post(srv.URL+"/v1/swift-codes", "application/json", strings.NewReader(j))
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/v1/swift-codes/TESTPLPWXXX", nil)
	resp, _ = http.DefaultClient.Do(req)
	require.Equal(t, 200, resp.StatusCode)
}

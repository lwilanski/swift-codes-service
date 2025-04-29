package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/your-github-name/swift-codes-service/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func makeRepo(t *testing.T) SwiftRepo {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.SwiftCode{}))
	return New(db)
}

func seed() []models.SwiftCode {
	return []models.SwiftCode{
		{SwiftCode: "BANKPLPWXXX", CountryISO2: "PL", CountryName: "POLAND", BankName: "HQ", IsHeadquarter: true},
		{SwiftCode: "BANKPLPWA01", CountryISO2: "PL", CountryName: "POLAND", BankName: "Branch A"},
	}
}

func TestUpsertAndQuery(t *testing.T) {
	r := makeRepo(t)
	require.NoError(t, r.UpsertMany(seed()))

	hq, err := r.Get("BANKPLPWXXX")
	require.NoError(t, err)
	require.True(t, hq.IsHeadquarter)

	branches, err := r.GetBranches("BANKPLPW")
	require.NoError(t, err)
	require.Len(t, branches, 1)
}

func TestDelete(t *testing.T) {
	r := makeRepo(t)
	_ = r.UpsertMany(seed())

	require.NoError(t, r.Delete("BANKPLPWA01"))
	_, err := r.Get("BANKPLPWA01")
	require.Error(t, err)
}

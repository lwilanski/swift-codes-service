package repository

import (
	"github.com/lwilanski/swift-codes-service/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SwiftRepo struct{ db *gorm.DB }

func New(db *gorm.DB) SwiftRepo { return SwiftRepo{db} }

func (r SwiftRepo) UpsertMany(list []models.SwiftCode) error {
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&list).Error
}

func (r SwiftRepo) Get(code string) (models.SwiftCode, error) {
	var s models.SwiftCode
	return s, r.db.First(&s, "swift_code = ?", code).Error
}

func (r SwiftRepo) GetBranches(root string) ([]models.SwiftCode, error) {
	var res []models.SwiftCode
	return res, r.db.
		Where("swift_code LIKE ? AND is_headquarter = false", root+"%").
		Find(&res).Error
}

func (r SwiftRepo) CountryAll(iso2 string) ([]models.SwiftCode, error) {
	var res []models.SwiftCode
	return res, r.db.Where("country_iso2 = ?", iso2).Find(&res).Error
}

func (r SwiftRepo) Delete(code string) error {
	return r.db.Delete(&models.SwiftCode{}, "swift_code = ?", code).Error
}

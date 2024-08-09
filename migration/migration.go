package migration

import (
	"boiler-plate-clean/internal/entity"
	"github.com/RumbiaID/pkg-library/app/pkg/database"
)

func AutoMigration(CpmDB *database.Database) {
	CpmDB.MigrateDB(
		&entity.Users{},
		&entity.Wallet{},
		&entity.CategoryTransaction{},
		&entity.Transaction{},
		&entity.Example{})

	//category maker
	var category *entity.CategoryTransaction
	CpmDB.GetDB().Model(&entity.CategoryTransaction{}).Where("name = transfer").First(&category)
	if category.Name == "" {
		category.Name = "transfer"
		if err := CpmDB.GetDB().Create(&category).Error; err != nil {
			panic(err)
		}
	}
	//&entity.SMSLog{}
}

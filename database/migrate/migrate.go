package migrate

import (
	"marketapi/database"
	"marketapi/model"
)

func MigrateTablesToDB() {
	database.Connect()
	database.Database.AutoMigrate(&model.User{})
	database.Database.AutoMigrate(&model.Card{})
	database.Database.AutoMigrate(&model.Wallet{})
	database.Database.AutoMigrate(&model.WalletTransaction{})
}

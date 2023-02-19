package database

import (
	postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"jwtauth/model"
)

var Database *gorm.DB
var DATABASE_URI string = "host=localhost user=root password=root dbname=postgres port=5432 sslmode=disable"

func Connect() error {
	var err error

	Database, err = gorm.Open(postgres.Open(DATABASE_URI), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		panic(err)
	}

	Database.AutoMigrate(&model.User{})
	Database.AutoMigrate(&model.Token{})

	return nil
}

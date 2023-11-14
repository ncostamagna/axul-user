package bootstrap

import (
	"fmt"
	domain "github.com/ncostamagna/axul_domain/domain/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func DBConnection() (*gorm.DB, error) {

	dsn := os.ExpandEnv("${DATABASE_USER}:${DATABASE_PASSWORD}@(${DATABASE_HOST}:${DATABASE_PORT})/${DATABASE_NAME}?charset=utf8&parseTime=True&loc=Local")

	fmt.Println("connect: ", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if os.Getenv("DATABASE_DEBUG") == "true" {
		db = db.Debug()
	}

	if os.Getenv("DATABASE_MIGRATE") == "true" {
		if err := db.AutoMigrate(&domain.User{}); err != nil {
			return nil, err
		}

		if err := db.AutoMigrate(&domain.Role{}); err != nil {
			return nil, err
		}
	}

	return db, nil
}

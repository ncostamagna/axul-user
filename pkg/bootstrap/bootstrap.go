package bootstrap

import (
	"fmt"
	domain "github.com/ncostamagna/axul_domain/domain/user"
	"github.com/ncostamagna/go-logger-hub/loghub"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func NewLogger() loghub.Logger {
	return loghub.New(loghub.NewNativeLogger(nil, 15))
}

/*
func NewLogger() (*slog.Logger, error) {

	date := time.Now()
	file, err := os.Create(fmt.Sprintf("log/%d.log", date.UnixNano()))
	if err != nil {
		return nil, err
	}
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {

			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
			}
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				if file, ok := strings.CutPrefix(source.File, wd); ok {
					source.File = file
				}

				f := strings.Split(source.Function, "/")
				source.Function = f[len(f)-1]
			}
			return a
		},
	})), nil
}*/

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

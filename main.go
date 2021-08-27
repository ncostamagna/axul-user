package main

import (
	"github.com/digitalhouse-dev/dh-kit/logger"
	"github.com/ncostamagna/axul_user/internal/user"
	"github.com/ncostamagna/axul_user/pkg/handler"

	"github.com/joho/godotenv"

	"context"
	"flag"
	"fmt"

	"net/http"

	"os"
	"os/signal"
	"syscall"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	fmt.Println("Initial")
	var log = logger.New(logger.LogOption{Debug: true})
	_ = godotenv.Overload()

	var httpAddr = flag.String("http", ":"+os.Getenv("APP_PORT"), "http listen address")

	fmt.Println("DataBases")
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		_ = log.CatchError(err)
		os.Exit(-1)
	}
	if os.Getenv("DATABASE_DEBUG") == "true" {
		db = db.Debug()
	}

	if os.Getenv("DATABASE_MIGRATE") == "true" {
		err := db.AutoMigrate(&user.User{})
		_ = log.CatchError(err)
	}

	flag.Parse()
	ctx := context.Background()

	var srv user.Service
	{
		repository := user.NewRepository(db, log)
		srv = user.NewService(repository, log)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	mux := http.NewServeMux()

	mux.Handle("/", handler.NewHTTPServer(ctx, user.MakeEndpoints(srv)))

	http.Handle("/", accessControl(mux))

	go func() {
		fmt.Println("listening on port", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()

	err = <-errs
	if err != nil {
		_ = log.CatchError(err)
	}

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

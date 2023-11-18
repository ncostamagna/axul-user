package main

import (
	"github.com/joho/godotenv"
	authentication "github.com/ncostamagna/axul_auth/auth"
	"github.com/ncostamagna/axul_user/internal/user"
	"github.com/ncostamagna/axul_user/internal/user/role"
	"github.com/ncostamagna/axul_user/pkg/bootstrap"
	"github.com/ncostamagna/axul_user/pkg/handler"
	"time"

	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {

	slog, err := bootstrap.NewLogger()
	if err != nil {
		panic(err)
	}

	_ = godotenv.Load()

	slog.Info("DataBases")
	db, err := bootstrap.DBConnection()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(-1)
	}

	flag.Parse()
	ctx := context.Background()

	token := os.Getenv("TOKEN")
	auth, err := authentication.New(token)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(-1)
	}

	var service user.Service
	{
		repository := user.NewRepository(db, slog)
		service = user.NewService(repository, auth, slog)
	}

	var roleService role.Service
	{
		repository := role.NewRepository(db, slog)
		roleService = role.NewService(repository, service, slog)
	}
	h := handler.NewHTTPServer(ctx, user.MakeEndpoints(service))
	h = handler.NewHTTPRolesServer(ctx, h, role.MakeEndpoints(roleService))

	url := os.Getenv("APP_URL")
	srv := &http.Server{
		Handler:      accessControl(h),
		Addr:         url,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  4 * time.Second,
	}

	errs := make(chan error)

	go func() {
		slog.Info(fmt.Sprintf("listening on %s", url))
		errs <- srv.ListenAndServe()
	}()

	err = <-errs
	if err != nil {
		slog.Error(err.Error())
	}

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

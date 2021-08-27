package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/digitalhouse-dev/dh-kit/response"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/ncostamagna/axul_user/internal/user"
)

//NewHTTPServer is a server handler
func NewHTTPServer(ctx context.Context, endpoints user.Endpoints) http.Handler {

	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllHandler,
		encodeResponse,
	)).Methods("GET")

	r.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Store),
		decodeStoreHandler,
		encodeResponse,
	)).Methods("POST")

	return r

}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func decodeStoreHandler(_ context.Context, r *http.Request) (interface{}, error) {
	var req user.StoreReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetAllHandler(_ context.Context, r *http.Request) (interface{}, error) {
	v := r.URL.Query()
	d, _ := strconv.ParseInt(v.Get("days"), 0, 64)

	req := user.GetAllReq{
		Days: d,
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.WriteHeader(200)
	return json.NewEncoder(w).Encode(r)
}

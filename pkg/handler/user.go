package handler

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/digitalhouse-dev/dh-kit/response"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/ncostamagna/axul_user/internal/user"
)

// NewHTTPServer is a server handler
func NewHTTPServer(_ context.Context, endpoints user.Endpoints) http.Handler {

	r := gin.Default()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Use(ginDecode())

	r.GET("/users/{id}/token/{token}", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Token),
		decodeTokenHandler,
		encodeResponse,
		opts...,
	)))

	r.GET("/users/{id}", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetHandler,
		encodeResponse,
		opts...,
	)))

	r.GET("/users", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllHandler,
		encodeResponse,
		opts...,
	)))

	r.POST("/users/login", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Login),
		decodeLoginHandler,
		encodeResponse,
		opts...,
	)))

	r.POST("/users", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Store),
		decodeStoreHandler,
		encodeResponse,
		opts...,
	)))

	return r

}

func ginDecode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "params", c.Params)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func decodeStoreHandler(_ context.Context, r *http.Request) (interface{}, error) {
	var req user.StoreReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	pp := ctx.Value("params").(gin.Params)
	req := user.GetReq{
		ID: pp.ByName("id"),
	}

	return req, nil
}

func decodeGetAllHandler(_ context.Context, r *http.Request) (interface{}, error) {
	req := user.GetAllReq{}

	return req, nil
}

func decodeLoginHandler(_ context.Context, r *http.Request) (interface{}, error) {
	req := user.LoginReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeTokenHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	pp := ctx.Value("params").(gin.Params)
	req := user.TokenReq{
		Token: pp.ByName("token"),
		ID:    pp.ByName("id"),
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.WriteHeader(200)
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/ncostamagna/axul-user/internal/user"
	"github.com/ncostamagna/go-http-utils/response"
	"net/http"
	"strconv"
)

// NewHTTPServer is a server handler
func NewHTTPServer(_ context.Context, endpoints user.Endpoints) http.Handler {

	r := gin.Default()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Use(ginDecode())


	//Deprecated
	r.GET("/users/:id/token/:token", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Token),
		decodeTokenHandler,
		encodeResponse,
		opts...,
	)))

	//Deprecated
	r.GET("/users/:id", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetHandler,
		encodeResponse,
		opts...,
	)))

	r.GET("/users", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetHandler,
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

	r.PATCH("/users/:id", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdate,
		encodeResponse,
		opts...,
	)))

	r.PUT("/users/:id/password", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.UpdatePassword),
		decodeUpdatePassword,
		encodeResponse,
		opts...,
	)))

	return r

}

func ginDecode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "params", c.Params)
		ctx = context.WithValue(ctx, "header", c.Request.Header)
		fmt.Println(c.Header)
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
	pp := ctx.Value("header").(http.Header)
	
	if len(pp["Authorization"]) < 1 {
		return nil, response.BadRequest("invalid authentication")
	}

	req := user.GetReq{
		Authorization: pp["Authorization"][0],
	}

	return req, nil
}

func decodeGetAllHandler(_ context.Context, r *http.Request) (interface{}, error) {
	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := user.GetAllReq{
		UserName: v.Get("username"),
		Limit:    limit,
		Page:     page,
	}

	return req, nil
}

func decodeUpdate(ctx context.Context, r *http.Request) (interface{}, error) {

	var req user.UpdateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	params := ctx.Value("params").(gin.Params)
	req.ID = params.ByName("id")

	return req, nil
}

func decodeUpdatePassword(ctx context.Context, r *http.Request) (interface{}, error) {

	var req user.UpdatePasswordReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	params := ctx.Value("params").(gin.Params)
	req.ID = params.ByName("id")

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

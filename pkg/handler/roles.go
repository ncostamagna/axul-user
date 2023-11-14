package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/digitalhouse-dev/dh-kit/response"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/ncostamagna/axul_user/internal/user/role"
	"net/http"
)

func NewHTTPRolesServer(_ context.Context, r http.Handler, endpoints role.Endpoints) http.Handler {

	var router *gin.Engine
	if r == nil {
		router = gin.Default()
	} else {
		router = r.(*gin.Engine)
	}

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	router.POST("/users/:id/apps", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeAppStoreHandler,
		encodeResponse,
		opts...,
	)))

	router.PUT("/users/:id/apps/:app", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.AddRoles),
		decodeAddRoleHandler,
		encodeResponse,
		opts...,
	)))

	router.GET("/users/:id/apps/:app", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetRole),
		decodeGetRoleHandler,
		encodeResponse,
		opts...,
	)))

	return router

}

func decodeAppStoreHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	var req role.AppReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	pp := ctx.Value("params").(gin.Params)
	req.ID = pp.ByName("id")

	return req, nil
}

func decodeAddRoleHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	var req role.AddRoles
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	pp := ctx.Value("params").(gin.Params)
	req.ID = pp.ByName("id")
	req.App = pp.ByName("app")
	return req, nil
}

func decodeGetRoleHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	var req role.AppReq

	pp := ctx.Value("params").(gin.Params)
	req.ID = pp.ByName("id")
	req.App = pp.ByName("app")
	return req, nil
}

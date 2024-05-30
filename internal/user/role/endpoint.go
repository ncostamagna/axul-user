package role

import (
	"context"
	//domain "github.com/ncostamagna/axul_domain/domain/user"
	"errors"
	"github.com/ncostamagna/go-http-utils/response"
	// auth "github.com/ncostamagna/axul_auth/auth"
)

type (
	AppReq struct {
		ID  string `json:"id"`
		App string `json:"app"`
	}

	AddRoles struct {
		ID    string   `json:"id"`
		App   string   `json:"app"`
		Roles []string `json:"roles"`
	}

	CreateRole struct {
		ID    string   `json:"id"`
		Apps  []string `json:"apps"`
		Roles []string `json:"roles"`
	}
)

type Controller func(ctx context.Context, request interface{}) (interface{}, error)

// Endpoints struct
type Endpoints struct {
	Create   Controller
	AddRoles Controller
	GetRole  Controller
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create:   makeCreateEndpoint(s),
		AddRoles: makeAddRolesEndpoint(s),
		GetRole:  makeGetRolesEndpoint(s),
	}
}

func makeCreateEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AppReq)

		if req.App == "" || req.ID == "" {
			return nil, response.BadRequest(ErrUserIDAndAppAreRequired.Error())
		}

		role, err := service.Create(ctx, req.ID, req.App)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("", role, nil), nil
	}
}

func makeAddRolesEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRoles)
		if req.App == "" || req.ID == "" {
			return nil, response.BadRequest(ErrUserIDAndAppAreRequired.Error())
		}

		if err := service.AddRole(ctx, req.ID, req.App, req.Roles); err != nil {
			if errors.As(err, &InvalidRole{}) {
				return nil, response.BadRequest(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("", req, nil), nil
	}
}

func makeGetRolesEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AppReq)

		if req.App == "" || req.ID == "" {
			return nil, response.BadRequest(ErrUserIDAndAppAreRequired.Error())
		}

		f := Filters{
			UserID: []string{req.ID},
			App:    []string{req.App},
		}

		roles, err := service.GetAll(ctx, f, 0, 0, "")
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		if len(roles) < 1 {
			return nil, response.NotFound(ErrUserAppNotFound{
				req.ID, req.App,
			}.Error())
		}

		return response.OK("", roles[0], nil), nil
	}
}

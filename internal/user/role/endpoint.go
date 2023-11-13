package role

import (
	"context"
	//domain "github.com/ncostamagna/axul_domain/domain/user"
	//"errors"
	"github.com/digitalhouse-dev/dh-kit/response"
	// auth "github.com/ncostamagna/axul_auth/auth"
)

type (
	CreateReq struct {
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
	Create Controller
	AddRoles  Controller
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		AddRoles:  makeAddRolesEndpoint(s),
	}
}

func makeCreateEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateReq)


		if req.App == "" || req.ID == ""{
			return nil, response.BadRequest("app and user id are required")
		}

		role, err := service.Create(ctx, req.ID, req.App)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Success("success", role, nil, nil), nil
	}
}

func makeAddRolesEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRoles)

		/*
			if req.UserName == "" || req.FirstName == "" || req.LastName == "" || req.Password == "" || req.Email == "" {
				return nil, response.BadRequest("fields required")
			}

			user, err := service.Create(ctx, req.UserName, req.FirstName, req.LastName, req.Password, req.Email, req.Phone, req.ClientID, req.ClientSecret, req.Token, req.Language)
			if err != nil {
				return nil, response.InternalServerError(err.Error())
			}*/

		return response.Success("success", req, nil, nil), nil
	}
}

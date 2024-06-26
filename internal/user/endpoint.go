package user

import (
	"context"
	"fmt"
	"errors"
	auth "github.com/ncostamagna/axul_auth/auth"
	domain "github.com/ncostamagna/axul_domain/domain/user"
	"github.com/ncostamagna/go-http-utils/meta"
	"github.com/ncostamagna/go-http-utils/response"
)

type (
	StoreReq struct {
		UserName     string `json:"username"`
		FirstName    string `json:"firstname"`
		LastName     string `json:"lastname"`
		Password     string `json:"password"`
		Email        string `json:"email"`
		Language     string `json:"language"`
		Phone        string `json:"phone"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Token        string `json:"token"`
	}

	GetAllReq struct {
		ID       []string `json:"id"`
		UserName string   `json:"username"`
		Limit    int      `json:"limit"`
		Page     int      `json:"page"`
	}

	GetReq struct {
		Authorization string `json:"Authorization"`
	}

	UpdateReq struct {
		ID        string  `json:"id"`
		FirstName *string `json:"firstname"`
		LastName  *string `json:"lastname"`
		Email     *string `json:"email"`
		Language  *string `json:"language"`
		Phone     *string `json:"phone"`
		Photo     *string `json:"photo"`
	}

	UpdatePasswordReq struct {
		ID          string `json:"id"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	LoginReq struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	LoginRes struct {
		User  *domain.User `json:"user"`
		Token string       `json:"token"`
	}

	TokenReq struct {
		ID    string `json:"id"`
		Token string `json:"token"`
	}

	AuthRes struct {
		Authorization int32        `json:"authorization"`
		User          *domain.User `json:"user"`
	}

	Config struct {
		LimPageDef string
	}
)

type Controller func(ctx context.Context, request interface{}) (interface{}, error)

// Endpoints struct
type Endpoints struct {
	Get            Controller
	GetAll         Controller
	Store          Controller
	Login          Controller
	Token          Controller
	Update         Controller
	UpdatePassword Controller
	Delete         Controller
}

func MakeEndpoints(s Service, config Config) Endpoints {
	return Endpoints{
		Get:            makeGetEndpoint(s),
		GetAll:         makeGetAllEndpoint(s, config),
		Store:          makeStoreEndpoint(s),
		Login:          makeLoginEndpoint(s),
		Token:          makeTokenEndpoint(s),
		Update:         makeUpdateEndpoint(s),
		UpdatePassword: makeUpdatePasswordEndpoint(s),
		Delete:         makeDeleteEndpoint(s),
	}
}

func makeGetEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetReq)
		fmt.Println(req)
		user, err := service.GetByToken(ctx, req.Authorization)
		if err != nil {
			if err == NotFound {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("", user, nil), nil
	}
}

func makeGetAllEndpoint(service Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAllReq)
		filters := Filters{ID: req.ID, UserName: req.UserName}

		count, err := service.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, config.LimPageDef)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		users, err := service.GetAll(ctx, filters, meta.Offset(), meta.Limit(), "")
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("", users, meta), nil
	}
}

func makeStoreEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StoreReq)

		if req.UserName == "" || req.FirstName == "" || req.LastName == "" || req.Password == "" || req.Email == "" {
			return nil, response.BadRequest("fields required")
		}

		user, err := service.Create(ctx, req.UserName, req.FirstName, req.LastName, req.Password, req.Email, req.Phone, req.ClientID, req.ClientSecret, req.Token, req.Language)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("", user, nil), nil
	}
}

func makeLoginEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginReq)

		users, err := service.GetAll(ctx, Filters{UserName: req.UserName}, 0, 0, "")
		if len(users) != 1 {
			return nil, response.Unauthorized(InvalidAuthentication.Error())
		}
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		token, err := service.Login(ctx, &users[0], req.Password)
		if err != nil {
			if err == InvalidAuthentication {
				return nil, response.Unauthorized(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("", LoginRes{&users[0], token}, nil), nil
	}
}

func makeTokenEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(TokenReq)
		user, err := service.TokenAccess(ctx, req.ID, req.Token)

		if err != nil {
			if err == NotFound {
				return nil, response.NotFound(err.Error())
			}

			if err == InvalidAuthentication || err == auth.ErrInvalidAuthentication {
				return nil, response.Unauthorized(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("", AuthRes{Authorization: 1, User: user}, nil), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateReq)

		if req.FirstName != nil && *req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}

		if req.LastName != nil && *req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}

		if req.Email != nil && *req.Email == "" {
			return nil, response.BadRequest(ErrEmailRequired.Error())
		}

		if err := s.Update(ctx, req.ID, req.FirstName, req.LastName, req.Email, req.Phone, req.Photo, req.Language); err != nil {

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("", nil, nil), nil
	}
}

func makeUpdatePasswordEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdatePasswordReq)

		if req.NewPassword == "" {
			return nil, response.BadRequest(ErrNewPasswordRequired.Error())
		}

		if req.OldPassword == "" {
			return nil, response.BadRequest(ErrOldPasswordRequired.Error())
		}

		if err := s.UpdatePassword(ctx, req.ID, req.NewPassword, req.OldPassword); err != nil {

			if err == InvalidPassword {
				return nil, response.BadRequest(err.Error())
			}

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("", nil, nil), nil
	}
}

func makeDeleteEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		return response.OK("", nil, nil), nil
	}
}

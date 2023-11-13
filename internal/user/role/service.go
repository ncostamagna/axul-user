package role

import (
	"context"
	"fmt"
	domain "github.com/ncostamagna/axul_domain/domain/user"
	//authentication "github.com/ncostamagna/axul_auth/auth"
	"github.com/ncostamagna/axul_user/internal/user"
	"github.com/digitalhouse-dev/dh-kit/logger"
)

type Service interface {
	Create(ctx context.Context, userId, app string) (*domain.Role, error)
	//Store(ctx context.Context, id string, application, roles []string) (*domain.User, error)
}

type service struct {
	repo    Repository
	userSrv user.Service
	//auth    authentication.Auth
	logger  logger.Logger
}

// NewService is a service handler
func NewService(repo Repository, userSrv user.Service, logger logger.Logger) Service {
	return &service{
		repo:    repo,
		userSrv: userSrv,
		logger:  logger,
	}
}

func (s *service) Create(ctx context.Context, userId, app string) (*domain.Role, error) {

	role := domain.Role{
		UserID: userId,
		App:    app,
	}

	if err := s.repo.Create(ctx, &role); err != nil {
		return nil, s.logger.CatchError(err)
	}
	s.logger.DebugMessage(fmt.Sprintf("Create %s Role", role.ID))

	return &role, nil

}

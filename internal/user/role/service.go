package role

import (
	"context"
	"fmt"
	"github.com/digitalhouse-dev/dh-kit/logger"
	domain "github.com/ncostamagna/axul_domain/domain/user"
	"github.com/ncostamagna/axul_user/internal/user"
)

type Filters struct {
	UserID []string
	App    []string
}

type Service interface {
	Create(ctx context.Context, userId, app string) (*domain.Role, error)
	AddRole(ctx context.Context, userId, app string, roles []string) error
	GetAll(ctx context.Context, filters Filters, offset, limit int, pload string) ([]domain.Role, error)
	Count(ctx context.Context, filters Filters) (int, error)
}

type service struct {
	repo    Repository
	userSrv user.Service
	//auth    authentication.Auth
	logger logger.Logger
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

func (s *service) AddRole(ctx context.Context, userId, app string, roles []string) error {

	role := domain.Role{
		UserID: userId,
		App:    app,
	}

	for _, r := range roles {
		if err := role.AddRole(r); err != nil {
			return InvalidRole{r}
		}
	}

	if err := s.repo.Update(ctx, role.UserID, role.App, &role.Role); err != nil {
		return err
	}

	return nil

}

func (s *service) GetAll(ctx context.Context, filters Filters, offset, limit int, pload string) ([]domain.Role, error) {
	roles, err := s.repo.GetAll(ctx, filters, offset, limit)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (s service) Count(ctx context.Context, filters Filters) (int, error) {
	return s.repo.Count(ctx, filters)
}

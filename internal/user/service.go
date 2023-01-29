package user

import (
	"context"
	"fmt"
	"github.com/ncostamagna/axul_domain/domain"

	authentication "github.com/ncostamagna/axul_auth/auth"

	"github.com/digitalhouse-dev/dh-kit/logger"
	"golang.org/x/crypto/bcrypt"
)

type Filters struct {
	ID       []string
	UserName string
}

type Service interface {
	Get(ctx context.Context, id, pload string) (*domain.User, error)
	//GetByUserName(ctx context.Context, username string) (*domain.User, error)
	GetAll(ctx context.Context, filters Filters, offset, limit int, pload string) ([]domain.User, error)
	Create(ctx context.Context, userName, firstName, lastName, password, email, phone, clientID, clientSecret, token string) (*domain.User, error)
	Update(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	Login(ctx context.Context, user *domain.User, password string) (string, error)
	TokenAccess(ctx context.Context, id, token string) (*domain.User, error)
	Count(ctx context.Context, filters Filters) (int, error)
}

type service struct {
	repo   Repository
	auth   authentication.Auth
	logger logger.Logger
}

// NewService is a service handler
func NewService(repo Repository, auth authentication.Auth, logger logger.Logger) Service {
	return &service{
		repo:   repo,
		auth:   auth,
		logger: logger,
	}
}

func (s *service) Get(ctx context.Context, id, pload string) (*domain.User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		_ = s.logger.CatchError(err)
		return nil, NotFound
	}

	s.logger.DebugMessage(fmt.Sprintf("Get %s User", user.ID))
	return user, nil
}

/*
func (s *service) GetByUserName(ctx context.Context, username string) (*domain.User, error) {
	user, err := s.repo.GetByUserName(ctx, username)
	if err != nil {
		_ = s.logger.CatchError(err)
		return nil, NotFound
	}

	s.logger.DebugMessage(fmt.Sprintf("Get %s User", user.ID))
	return user, nil
}*/

func (s *service) GetAll(ctx context.Context, filters Filters, offset, limit int, pload string) ([]domain.User, error) {
	users, err := s.repo.GetAll(ctx, filters, offset, limit)
	if err != nil {
		return nil, s.logger.CatchError(err)
	}

	s.logger.DebugMessage(fmt.Sprintf("Get %d Users", len(users)))
	return users, nil
}

func (s *service) Create(ctx context.Context, userName, firstName, lastName, password, email, phone, clientID, clientSecret, token string) (*domain.User, error) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, s.logger.CatchError(err)
	}

	user := domain.User{
		UserName:     userName,
		FirstName:    firstName,
		LastName:     lastName,
		Password:     string(hashPassword),
		Email:        email,
		Phone:        phone,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Token:        token,
	}

	if err := s.repo.Create(ctx, &user); err != nil {
		return nil, s.logger.CatchError(err)
	}
	s.logger.DebugMessage(fmt.Sprintf("Create %s User", user.ID))

	return &user, nil

}

func (s *service) Update(ctx context.Context, id string) error {
	return nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *service) Login(ctx context.Context, user *domain.User, password string) (string, error) {

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		_ = s.logger.CatchError(err)
		return "", InvalidAuthentication
	}

	token, err := s.auth.Create(user.ID, user.UserName, 0)
	if err != nil {
		_ = s.logger.CatchError(err)
		return "", InvalidAuthentication
	}
	return token, nil
}

func (s *service) TokenAccess(ctx context.Context, id, token string) (*domain.User, error) {

	if err := s.auth.Access(id, token); err != nil {
		return nil, err
	}

	user, err := s.repo.Get(ctx, id)
	if err != nil {
		_ = s.logger.CatchError(err)
		return nil, NotFound
	}

	s.logger.DebugMessage(fmt.Sprintf("Get %s User with token %s", id, token))

	return user, nil

}
func (s service) Count(ctx context.Context, filters Filters) (int, error) {
	return s.repo.Count(ctx, filters)
}

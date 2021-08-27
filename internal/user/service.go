package user

import (
	"context"
	"fmt"

	"github.com/digitalhouse-dev/dh-kit/logger"
)

type Filters struct {
	ID   []string
	Days int64
}

type Service interface {
	Get(ctx context.Context, id, pload string) (*User, error)
	GetAll(ctx context.Context, filters Filters, offset, limit int, pload string) (*[]User, error)
	Create(ctx context.Context, userName, firstName, lastName, password, email, phone, clientID, clientSecret, token string) (*User, error)
	Update(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo   Repository
	logger logger.Logger
}

//NewService is a service handler
func NewService(repo Repository, logger logger.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) Get(ctx context.Context, id, pload string) (*User, error) {
	return nil, nil
}

func (s *service) GetAll(ctx context.Context, filters Filters, offset, limit int, pload string) (*[]User, error) {
	users, err := s.repo.GetAll(ctx, filters, offset, limit)
	if err != nil {
		return nil, s.logger.CatchError(err)
	}

	s.logger.DebugMessage(fmt.Sprintf("Get %d Users", len(*users)))
	return users, nil
}

func (s *service) Create(ctx context.Context, userName, firstName, lastName, password, email, phone, clientID, clientSecret, token string) (*User, error) {

	user := User{
		UserName:     userName,
		FirstName:    firstName,
		LastName:     lastName,
		Password:     password,
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

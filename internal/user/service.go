package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/digitalhouse-dev/dh-kit/logger"
	"golang.org/x/crypto/bcrypt"
)

type Filters struct {
	ID   []string
	Days int64
}

var NotFound = errors.New("Record not found")
var FieldIsRequired = errors.New("Required values")
var InvalidAuthentication = errors.New("Invalid authentication")

type Service interface {
	Get(ctx context.Context, id, pload string) (*User, error)
	GetByUserName(ctx context.Context, username string) (*User, error)
	GetAll(ctx context.Context, filters Filters, offset, limit int, pload string) (*[]User, error)
	Create(ctx context.Context, userName, firstName, lastName, password, email, phone, clientID, clientSecret, token string) (*User, error)
	Update(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	Login(ctx context.Context, user *User, password string) (string, error)
	TokenAccess(ctx context.Context, id, token string) error
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
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		_ = s.logger.CatchError(err)
		return nil, NotFound
	}

	s.logger.DebugMessage(fmt.Sprintf("Get %s User", user.ID))
	return user, nil
}

func (s *service) GetByUserName(ctx context.Context, username string) (*User, error) {
	user, err := s.repo.GetByUserName(ctx, username)
	if err != nil {
		_ = s.logger.CatchError(err)
		return nil, NotFound
	}

	s.logger.DebugMessage(fmt.Sprintf("Get %s User", user.ID))
	return user, nil
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

	if userName == "" || firstName == "" || lastName == "" || password == "" || email == "" {
		return nil, s.logger.CatchError(FieldIsRequired)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, s.logger.CatchError(err)
	}

	user := User{
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

func (s *service) Login(ctx context.Context, user *User, password string) (string, error) {

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		_ = s.logger.CatchError(err)
		return "", InvalidAuthentication
	}

	token, err := CreateJWT(user.ID, user.UserName, 0)
	encToken := encrypt(token, "6470fc52afd689ca17df8667729b2c0460ce90b781a01b0010d2c4c31c85cb21")
	if err != nil {
		_ = s.logger.CatchError(err)
		return "", InvalidAuthentication
	}
	return encToken, nil
}

func (s *service) TokenAccess(ctx context.Context, id, token string) error {

	decToken := decrypt(token, "6470fc52afd689ca17df8667729b2c0460ce90b781a01b0010d2c4c31c85cb21")
	user, err := AccessJWT(decToken)

	if err != nil || user.ID != id {
		return InvalidAuthentication
	}

	return nil
}
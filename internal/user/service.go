package user

import (
	"context"
	"fmt"
	domain "github.com/ncostamagna/axul_domain/domain/user"
	"github.com/ncostamagna/go-logger-hub/loghub"

	authentication "github.com/ncostamagna/axul_auth/auth"

	"golang.org/x/crypto/bcrypt"
)

type Filters struct {
	ID       []string
	UserName string
}

type Service interface {
	Get(ctx context.Context, id, pload string) (*domain.User, error)
	GetByToken(ctx context.Context, token string) (*domain.User, error)
	GetAll(ctx context.Context, filters Filters, offset, limit int, pload string) ([]domain.User, error)
	Create(ctx context.Context, userName, firstName, lastName, password, email, phone, clientID, clientSecret, token, language string) (*domain.User, error)
	Update(ctx context.Context, id string, firstname, lastname, email, phone, photo, language *string) error
	UpdatePassword(ctx context.Context, id, newPassword, oldPassword string) error
	Delete(ctx context.Context, id string) error
	Login(ctx context.Context, user *domain.User, password string) (string, error)
	TokenAccess(ctx context.Context, id, token string) (*domain.User, error)
	Count(ctx context.Context, filters Filters) (int, error)
}

type service struct {
	repo   Repository
	auth   authentication.Auth
	logger loghub.Logger
}

// NewService is a service handler
func NewService(repo Repository, auth authentication.Auth, logger loghub.Logger) Service {
	return &service{
		repo:   repo,
		auth:   auth,
		logger: logger,
	}
}

func (s *service) Get(ctx context.Context, id, pload string) (*domain.User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Warn(err)
		return nil, NotFound
	}

	s.logger.Info(fmt.Sprintf("Get %s User", user.ID))
	return user, nil
}

func (s *service) GetByToken(ctx context.Context, token string) (*domain.User, error) {
	u, err := s.auth.Check(token)
	fmt.Println(u)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, u.ID, "")
}

func (s *service) GetAll(ctx context.Context, filters Filters, offset, limit int, pload string) ([]domain.User, error) {
	users, err := s.repo.GetAll(ctx, filters, offset, limit)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	s.logger.Info(fmt.Sprintf("Get %d Users", len(users)))
	return users, nil
}

func (s *service) Create(ctx context.Context, userName, firstName, lastName, password, email, phone, clientID, clientSecret, token, language string) (*domain.User, error) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	var lang domain.Language

	switch domain.Language(language) {
	case domain.English, domain.Spanish:
		lang = domain.Language(language)
	default:
		lang = domain.English
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
		Language:     lang,
	}

	if err := s.repo.Create(ctx, &user); err != nil {
		s.logger.Error(err)
		return nil, err
	}
	s.logger.Info(fmt.Sprintf("Create %s User", user.ID))

	return &user, nil

}

func (s *service) Update(ctx context.Context, id string, firstname, lastname, email, phone, photo, language *string) error {

	var lang *string
	if language != nil {
		switch domain.Language(*language) {
		case domain.English, domain.Spanish:
			lang = language
		default:
			l := string(domain.English)
			lang = &l
		}
	}

	if err := s.repo.Update(ctx, id, firstname, lastname, email, phone, photo, lang, nil); err != nil {
		return err
	}

	return nil
}

func (s *service) UpdatePassword(ctx context.Context, id, newPassword, oldPassword string) error {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Error(err)
		return ErrNotFound{id}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		s.logger.Error(err)
		return InvalidPassword
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	hashNewPassword := string(hashPassword)
	if err := s.repo.Update(ctx, id, nil, nil, nil, nil, nil, nil, &hashNewPassword); err != nil {
		return err
	}

	return nil
}
func (s *service) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *service) Login(ctx context.Context, user *domain.User, password string) (string, error) {

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.logger.Error(err)
		return "", InvalidAuthentication
	}

	token, err := s.auth.Create(user.ID, user.UserName, "", true, 0)
	if err != nil {
		s.logger.Error(err)
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
		s.logger.Warn(err)
		return nil, NotFound
	}

	s.logger.Info(fmt.Sprintf("Get %s User with token %s", id, token))

	return user, nil

}
func (s service) Count(ctx context.Context, filters Filters) (int, error) {
	return s.repo.Count(ctx, filters)
}

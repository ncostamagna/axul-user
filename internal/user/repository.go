package user

import (
	"context"
	"github.com/digitalhouse-dev/dh-kit/logger"
	"github.com/google/uuid"
	"github.com/ncostamagna/axul_domain/domain"
	"gorm.io/gorm"
)

type Repository interface {
	GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error)
	Get(ctx context.Context, id string) (*domain.User, error)
	GetByUserName(ctx context.Context, username string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

type repo struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewRepository(db *gorm.DB, log logger.Logger) Repository {
	return &repo{db, log}
}

func (r *repo) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error) {
	var tx *gorm.DB
	var user []domain.User
	tx = r.db.WithContext(ctx).Model(&user)

	result := tx.Order("created_at desc").Find(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *repo) Get(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	tx := r.db.WithContext(ctx).Model(&user)

	result := tx.Where("id = ?", id).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (r *repo) GetByUserName(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	tx := r.db.WithContext(ctx).Model(&user)

	result := tx.Where("user_name = ?", username).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (r *repo) Create(ctx context.Context, user *domain.User) error {
	user.ID = uuid.New().String()
	return r.db.Create(&user).Error
}

func (r *repo) Update(ctx context.Context, id string) error {
	return nil
}

func (r *repo) Delete(ctx context.Context, id string) error {
	return nil
}

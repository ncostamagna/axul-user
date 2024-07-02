package user

import (
	"context"
	"github.com/google/uuid"
	domain "github.com/ncostamagna/axul_domain/domain/user"
	"github.com/ncostamagna/go-logger-hub/loghub"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error)
	Get(ctx context.Context, id string) (*domain.User, error)
	//GetByUserName(ctx context.Context, username string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, id string, firstname, lastname, email, phone, photo, language, password *string) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filters Filters) (int, error)
}

type repo struct {
	db     *gorm.DB
	logger loghub.Logger
}

func NewRepository(db *gorm.DB, logger loghub.Logger) Repository {
	return &repo{db, logger}
}

func (r *repo) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error) {
	var user []domain.User

	tx := r.db.WithContext(ctx).Model(&user)
	applyFilters(tx, filters)
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

/*
func (r *repo) GetByUserName(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	tx := r.db.WithContext(ctx).Model(&user)

	result := tx.Where("user_name = ?", username).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}*/

func (r *repo) Create(ctx context.Context, user *domain.User) error {
	user.ID = uuid.New().String()
	return r.db.Create(&user).Error
}

func (r *repo) Update(ctx context.Context, id string, firstname, lastname, email, phone, photo, language, password *string) error {

	values := make(map[string]interface{})

	if firstname != nil {
		values["first_name"] = *firstname
	}

	if lastname != nil {
		values["last_name"] = *lastname
	}

	if email != nil {
		values["email"] = *email
	}

	if phone != nil {
		values["phone"] = *phone
	}

	if photo != nil {
		values["photo"] = *photo
	}

	if language != nil {
		values["language"] = *language
	}

	if password != nil {
		values["password"] = *password
	}

	result := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		r.logger.Error(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound{id}
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *repo) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	tx := r.db.WithContext(ctx).Model(domain.User{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		r.logger.Error(err)
		return 0, err
	}

	return int(count), nil
}

func applyFilters(tx *gorm.DB, f Filters) *gorm.DB {

	if f.UserName != "" {
		tx = tx.Where("lower(user_name) = ?", strings.ToLower(f.UserName))
	}

	return tx
}

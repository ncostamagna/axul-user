package role

import (
	"context"
	domain "github.com/ncostamagna/axul_domain/domain/user"
	"github.com/ncostamagna/go-logger-hub/loghub"
	"gorm.io/gorm"
)

type Repository interface {
	GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Role, error)
	//Get(ctx context.Context, id string) (*domain.User, error)
	Create(ctx context.Context, role *domain.Role) error
	Update(ctx context.Context, userID, app string, role *uint64) error
	//Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filters Filters) (int, error)
}

type repo struct {
	db     *gorm.DB
	logger loghub.Logger
}

func NewRepository(db *gorm.DB, log loghub.Logger) Repository {
	return &repo{db, log}
}

func (r *repo) Create(ctx context.Context, role *domain.Role) error {
	return r.db.Create(&role).Error
}

func (r *repo) Update(ctx context.Context, userID, app string, role *uint64) error {
	values := make(map[string]interface{})

	if role != nil {
		values["role"] = *role
	}

	result := r.db.WithContext(ctx).Model(&domain.Role{}).Where("user_id = ? and app = ?", userID, app).Updates(values)
	if err := result.Error; err != nil {
		r.logger.Error(err)
		return err
	}

	if result.RowsAffected == 0 {
		return ErrUserAppNotFound{userID, app}
	}

	return nil
}

func (r *repo) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Role, error) {
	var role []domain.Role

	tx := r.db.WithContext(ctx).Model(&role)
	applyFilters(tx, filters)
	result := tx.Order("created_at desc").Find(&role)

	if err := result.Error; err != nil {
		r.logger.Error(err)
		return nil, err
	}

	return role, nil
}

func (r *repo) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	tx := r.db.WithContext(ctx).Model(domain.Role{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		r.logger.Error(err)
		return 0, err
	}

	return int(count), nil
}

/*
func (r *repo) Get(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	tx := r.db.WithContext(ctx).Model(&user)

	result := tx.Where("id = ?", id).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}


func (r *repo) Update(ctx context.Context, id string, firstname, lastname, email, phone,photo, language, password *string) error {

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
		return r.logger.CatchError(result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound{id}
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id string) error {
	return nil
}
*/

func applyFilters(tx *gorm.DB, f Filters) *gorm.DB {

	if f.UserID != nil {
		tx = tx.Where("user_id in (?)", f.UserID)
	}

	if f.App != nil {
		tx = tx.Where("app in (?)", f.App)
	}

	return tx
}

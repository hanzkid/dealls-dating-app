package repository

import (
	"main/entity"
	"main/helpers"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	FindByID(id int) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	Save(user *entity.User) (*entity.User, error)
	Update(user *entity.User) (*entity.User, error)
	Subscribe(c echo.Context) (*entity.User, error)
	CheckSubscription(c echo.Context) (bool, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindByID(id int) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Profile").Preload("Subscription").First(&user, id).Error
	if err != nil {
		return &entity.User{}, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).Preload("Profile").First(&user).Error
	if err != nil {
		return &entity.User{}, err
	}
	return &user, nil
}

func (r *UserRepository) Save(user *entity.User) (*entity.User, error) {
	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		return &entity.User{}, err
	}
	user.Password = *hashedPassword
	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if err := tx.Create(&entity.Profile{
			UserID: user.ID,
		}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &entity.User{}, err
	}
	return user, nil
}

func (r *UserRepository) Update(user *entity.User) (*entity.User, error) {
	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		return &entity.User{}, err
	}
	user.Password = *hashedPassword
	err = r.db.Save(&user).Error
	if err != nil {
		return &entity.User{}, err
	}
	return user, nil
}

func (r *UserRepository) Subscribe(ctx echo.Context) (*entity.User, error) {
	userId := uint(ctx.Get("user_id").(int))
	var user entity.User
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&entity.Subscription{
			UserID:     userId,
			ValidUntil: time.Now().AddDate(0, 1, 0),
		}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &entity.User{}, err
	}

	// Load the created subscription into the user object
	err = r.db.Preload("Subscription").First(&user, userId).Error
	if err != nil {
		return &entity.User{}, err
	}

	return &user, nil
}

func (r *UserRepository) CheckSubscription(ctx echo.Context) (bool, error) {
	userId := ctx.Get("user_id")
	var subscription entity.Subscription
	err := r.db.Where("user_id = ?", userId).First(&subscription).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	if time.Now().After(subscription.ValidUntil) {
		return false, nil
	}
	return true, nil
}

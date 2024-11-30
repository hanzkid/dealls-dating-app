package repository

import (
	"main/entity"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ProfileRepositoryInterface interface {
	FindByUserID(userId int) (*entity.Profile, error)
	FindByID(id int) (*entity.Profile, error)
	Save(profile *entity.Profile) (*entity.Profile, error)
	GetRandomProfile(c echo.Context) (*entity.Profile, error)
	SaveViewLog(c echo.Context, profileId int) error
}

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepositoryInterface {
	return &ProfileRepository{
		db: db,
	}
}

func (r *ProfileRepository) FindByID(id int) (*entity.Profile, error) {
	var profile entity.Profile
	if err := r.db.First(&profile, id).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) Save(profile *entity.Profile) (*entity.Profile, error) {
	if err := r.db.Save(profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *ProfileRepository) GetRandomProfile(ctx echo.Context) (*entity.Profile, error) {
	viewerId := ctx.Get("user_id").(int)
	var profile entity.Profile
	if err := r.db.
		Where("id NOT IN (?)", r.db.Table("profile_view_logs").Select("profile_id").Where("viewer_id = ? AND DATE(created_at) = DATE(NOW())", viewerId)).
		Order("RANDOM()").
		First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) SaveViewLog(ctx echo.Context, profileId int) error {
	viewerId := ctx.Get("profile_id").(int)
	viewLog := entity.ProfileViewLog{
		ViewerID:  uint(viewerId),
		ProfileID: uint(profileId),
	}
	if err := r.db.Save(&viewLog).Error; err != nil {
		return err
	}
	return nil
}

func (r *ProfileRepository) FindByUserID(userId int) (*entity.Profile, error) {
	var profile entity.Profile
	if err := r.db.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

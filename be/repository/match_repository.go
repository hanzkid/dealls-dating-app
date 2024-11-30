package repository

import (
	"errors"
	"main/entity"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type MatchRepositoryInterface interface {
	FindMatchByProfileID(profileID int) ([]*entity.Profile, error)
	CheckMatch(profileID, partnerID int) (*entity.Match, error)
	CheckPendingMatch(profileID, partnerID int) (*entity.Match, error)
	AcceptMatch(profileID, partnerID int) error
	RejectMatch(profileID, partnerID int) error
	CreateMatch(profileID, partnerID int) error
	CheckDailyLimit(ctx echo.Context) (bool, error)
}

type MatchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepositoryInterface {
	return &MatchRepository{
		db: db,
	}
}

func (r *MatchRepository) FindMatchByProfileID(profileID int) ([]*entity.Profile, error) {
	var matches []*entity.Profile
	var matchesAsInitiator []entity.Match
	if err := r.db.Preload("Partner").Where("profile_id = ? AND status = ?", profileID, entity.StatusAccepted).Find(&matchesAsInitiator).Error; err != nil {
		return nil, err
	}
	for _, match := range matchesAsInitiator {
		matches = append(matches, &match.Partner)
	}

	var matchesAsPartner []entity.Match
	if err := r.db.Preload("Profile").Where("partner_id = ? AND status = ?", profileID, entity.StatusAccepted).Find(&matchesAsPartner).Error; err != nil {
		return nil, err
	}
	for _, match := range matchesAsPartner {
		matches = append(matches, &match.Profile)
	}

	return matches, nil
}

func (r *MatchRepository) CheckMatch(profileID, partnerID int) (*entity.Match, error) {

	var match entity.Match
	if err := r.db.Where("profile_id = ? AND partner_id = ? AND status != ?", profileID, partnerID, entity.StatusRejected).First(&match).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &match, nil
}

func (r *MatchRepository) AcceptMatch(profileID, partnerID int) error {
	match, err := r.CheckMatch(partnerID, profileID)
	if err != nil {
		return err
	}
	if match != nil {
		if err := r.db.Model(&entity.Match{}).Where("profile_id = ? AND partner_id = ?", partnerID, profileID).Update("status", entity.StatusAccepted).Error; err != nil {
			return err
		}
		return nil
	}
	return nil
}

// CheckPendingMatch checks if there is a pending match between two profiles
func (r *MatchRepository) CheckPendingMatch(profileID, partnerID int) (*entity.Match, error) {
	var match entity.Match
	if err := r.db.Where("profile_id = ? AND partner_id = ? AND status = ?", profileID, partnerID, entity.StatusPending).First(&match).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &match, nil
}

func (r *MatchRepository) RejectMatch(profileID, partnerID int) error {
	match, err := r.CheckMatch(partnerID, profileID)
	if err != nil {
		return err
	}
	if match != nil {
		if err := r.db.Model(&entity.Match{}).Where("profile_id = ? AND partner_id = ?", partnerID, profileID).Update("status", entity.StatusRejected).Error; err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (r *MatchRepository) CreateMatch(profileID, partnerID int) error {
	match := &entity.Match{
		ProfileID: profileID,
		PartnerID: partnerID,
		Status:    entity.StatusPending,
	}
	if err := r.db.Create(match).Error; err != nil {
		return err
	}
	return nil
}

func (r *MatchRepository) CheckDailyLimit(ctx echo.Context) (bool, error) {
	var isPremium bool
	query := "SELECT EXISTS(SELECT 1 FROM subscriptions WHERE user_id = ? AND valid_until > NOW()) AS isPremium"
	if err := r.db.Raw(query, ctx.Get("user_id").(int)).Scan(&isPremium).Error; err != nil {
		return false, err
	}
	if isPremium {
		return true, nil
	}
	viewerId := ctx.Get("profile_id").(int)
	var count int64
	if err := r.db.Model(&entity.ProfileViewLog{}).
		Where("viewer_id = ? AND DATE(created_at) = DATE(NOW())", viewerId).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count < 10, nil
}

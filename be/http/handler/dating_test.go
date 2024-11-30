package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"main/entity"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProfileRepository struct {
	mock.Mock
}

type MockMatchRepository struct {
	mock.Mock
}

func (m *MockProfileRepository) FindByUserID(userId int) (*entity.Profile, error) {
	args := m.Called(userId)
	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockProfileRepository) Save(profile *entity.Profile) (*entity.Profile, error) {
	args := m.Called(profile)
	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockProfileRepository) GetRandomProfile(c echo.Context) (*entity.Profile, error) {
	args := m.Called(c)
	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockProfileRepository) SaveViewLog(c echo.Context, profileID int) error {
	args := m.Called(c, profileID)
	return args.Error(0)
}

func (m *MockProfileRepository) FindByID(profileID int) (*entity.Profile, error) {
	args := m.Called(profileID)
	return args.Get(0).(*entity.Profile), args.Error(1)
}

func (m *MockMatchRepository) CheckDailyLimit(c echo.Context) (bool, error) {
	args := m.Called(c)
	return args.Bool(0), args.Error(1)
}

func (m *MockMatchRepository) CheckPendingMatch(partnerID, profileID int) (*entity.Match, error) {
	args := m.Called(partnerID, profileID)
	return args.Get(0).(*entity.Match), args.Error(1)
}

func (m *MockMatchRepository) RejectMatch(profileID, partnerID int) error {
	args := m.Called(profileID, partnerID)
	return args.Error(0)
}

func (m *MockMatchRepository) CheckMatch(profileID, partnerID int) (*entity.Match, error) {
	args := m.Called(profileID, partnerID)
	return args.Get(0).(*entity.Match), args.Error(1)
}

func (m *MockMatchRepository) AcceptMatch(profileID, partnerID int) error {
	args := m.Called(profileID, partnerID)
	return args.Error(0)
}

func (m *MockMatchRepository) CreateMatch(profileID, partnerID int) error {
	args := m.Called(profileID, partnerID)
	return args.Error(0)
}

func (m *MockMatchRepository) FindMatchByProfileID(profileID int) ([]*entity.Profile, error) {
	args := m.Called(profileID)
	return args.Get(0).([]*entity.Profile), args.Error(1)
}

func TestProfile(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockProfileRepo := new(MockProfileRepository)
	mockMatchRepo := new(MockMatchRepository)

	handler := NewDatingHandler(mockProfileRepo, mockMatchRepo)

	mockMatchRepo.On("CheckDailyLimit", c).Return(true, nil)
	mockProfile := &entity.Profile{ID: 1}
	mockProfileRepo.On("GetRandomProfile", c).Return(mockProfile, nil)
	mockProfileRepo.On("SaveViewLog", c, 1).Return(nil)

	err := handler.Profile(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestProfileDailyLimit(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockProfileRepo := new(MockProfileRepository)
	mockMatchRepo := new(MockMatchRepository)

	handler := NewDatingHandler(mockProfileRepo, mockMatchRepo)

	mockMatchRepo.On("CheckDailyLimit", c).Return(false, nil)

	err := handler.Profile(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestSwipedProfile(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("profile_id", 1)

	mockProfileRepo := new(MockProfileRepository)
	mockMatchRepo := new(MockMatchRepository)

	handler := NewDatingHandler(mockProfileRepo, mockMatchRepo)
	reqBody := `{"profile_id": 2, "swipe": true}`
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)
	c.Set("profile_id", 1)

	mockMatchRepo.On("CheckDailyLimit", c).Return(true, nil)
	mockProfileRepo.On("FindByID", 1).Return(&entity.Profile{ID: 1}, nil)
	mockMatchRepo.On("CheckMatch", 1, 2).Return((*entity.Match)(nil), nil)
	mockMatchRepo.On("CheckMatch", 2, 1).Return((*entity.Match)(nil), nil)
	mockMatchRepo.On("CreateMatch", 1, 2).Return(nil)
	mockProfileRepo.On("GetRandomProfile", c).Return(&entity.Profile{ID: 2}, nil)
	mockProfileRepo.On("SaveViewLog", c, 2).Return(nil)

	err := handler.SwipedProfile(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestSwipedProfileDailyLimit(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("profile_id", 1)

	mockProfileRepo := new(MockProfileRepository)
	mockMatchRepo := new(MockMatchRepository)

	handler := NewDatingHandler(mockProfileRepo, mockMatchRepo)
	reqBody := `{"profile_id": 2, "swipe": true}`
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)
	c.Set("profile_id", 1)

	mockMatchRepo.On("CheckDailyLimit", c).Return(false, nil)
	mockProfileRepo.On("FindByID", 1).Return(&entity.Profile{ID: 1}, nil)
	mockMatchRepo.On("CheckMatch", 1, 2).Return((*entity.Match)(nil), nil)
	mockMatchRepo.On("CheckMatch", 2, 1).Return((*entity.Match)(nil), nil)
	mockMatchRepo.On("CreateMatch", 1, 2).Return(nil)
	mockProfileRepo.On("GetRandomProfile", c).Return(&entity.Profile{ID: 2}, nil)
	mockProfileRepo.On("SaveViewLog", c, 2).Return(nil)

	err := handler.SwipedProfile(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestMatchList(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("profile_id", 1)

	mockProfileRepo := new(MockProfileRepository)
	mockMatchRepo := new(MockMatchRepository)

	handler := NewDatingHandler(mockProfileRepo, mockMatchRepo)

	mockMatches := []*entity.Profile{{ID: 1}, {ID: 2}}
	mockMatchRepo.On("FindMatchByProfileID", 1).Return(mockMatches, nil)

	err := handler.MatchList(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

package helpers

import (
	"main/config"
	"main/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotNil(t, hashedPassword)

	err = bcrypt.CompareHashAndPassword([]byte(*hashedPassword), []byte(password))
	assert.NoError(t, err)
}

func TestComparePassword(t *testing.T) {
	password := "testpassword"
	hashedPassword, _ := HashPassword(password)

	isValid := ComparePassword(*hashedPassword, password)
	assert.True(t, isValid)

	isValid = ComparePassword(*hashedPassword, "wrongpassword")
	assert.False(t, isValid)
}

func TestGenerateAccessToken(t *testing.T) {
	user := &entity.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
		Profile: entity.Profile{
			ID: 1,
		},
	}
	cfg := &config.JWT{
		Secret: "testsecret",
		Expiry: 3600,
	}

	token, err := GenerateAccessToken(user, cfg)
	assert.NoError(t, err)
	assert.NotNil(t, token)
}

func TestValidateToken(t *testing.T) {
	user := &entity.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
		Profile: entity.Profile{
			ID: 1,
		},
	}
	cfg := &config.JWT{
		Secret: "testsecret",
		Expiry: 3600,
	}

	tokenString, _ := GenerateAccessToken(user, cfg)
	token, err := ValidateToken(*tokenString, cfg.Secret)
	assert.NoError(t, err)
	assert.NotNil(t, token)

	invalidToken := "invalidtoken"
	token, err = ValidateToken(invalidToken, cfg.Secret)
	assert.Error(t, err)
	assert.Nil(t, token)
}

func TestConvertStringToInt(t *testing.T) {
	str := "123"
	num := ConvertStringToInt(str)
	assert.Equal(t, 123, num)

	str = "invalid"
	num = ConvertStringToInt(str)
	assert.Equal(t, 0, num)
}

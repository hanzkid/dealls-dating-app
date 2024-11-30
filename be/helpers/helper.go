package helpers

import (
	"log"
	"main/config"
	"main/entity"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (*string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return nil, err
	}
	hashedPasswordStr := string(hashedPassword)
	return &hashedPasswordStr, nil
}

func ComparePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenerateAccessToken(user *entity.User, cfg *config.JWT) (*string, error) {
	secretKey := []byte(cfg.Secret)

	claims := &jwt.MapClaims{
		"user_id":    user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"profile_id": user.Profile.ID,
		"iss":        "dating-app",
		"exp":        time.Now().Add(time.Second * time.Duration(cfg.Expiry)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return nil, err
	}

	return &tokenString, nil
}

func ValidateToken(jwtString string, secret string) (*jwt.Token, error) {
	secretKey := []byte(secret)

	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		log.Printf("Failed to validate token: %v", err)
		return nil, err
	}

	return token, nil
}

func ConvertStringToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("Failed to convert string to int: %v", err)
		return 0
	}
	return num
}

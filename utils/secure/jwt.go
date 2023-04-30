package secure

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/johnyeocx/usual/server2/utils/enums"
)

var (
	accessTokenExpiry = time.Minute * 200
	refreshTokenExpiry = time.Minute * 200
)

func GenerateAccessToken(userID string, userType string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	secretKey := os.Getenv("JWT_ACCESS_SECRET")
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["type"] = userType

	claims["exp"] = time.Now().Add(accessTokenExpiry).Unix()
	tokenString, err := token.SignedString([]byte(secretKey))


	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(userID string, userType string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	secretKey := os.Getenv("JWT_REFRESH_SECRET")

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID

	claims["type"] = userType


	claims["exp"] = time.Now().Add(refreshTokenExpiry).Unix()
	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseAccessToken(tokenStr string) (string, string, error) {
	
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("invalid token: %v", token.Header["alg"])
		}
		secretKey := os.Getenv("JWT_ACCESS_SECRET")
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Println("Invalid jwt token:", err)
		return "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := claims["user_id"].(string)
		userType := claims["type"].(string)
		return userId, userType, nil
		
	} else {
		log.Print("Err", err)
		return "", "", err
	}
}

func ParseRefreshToken(tokenStr string) (string, enums.TokenType, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("invalid token: %v", token.Header["alg"])
		}
		secretKey := os.Getenv("JWT_REFRESH_SECRET")
		return []byte(secretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(string)
		tokenType := claims["type"].(enums.TokenType)
		return userID, tokenType, nil
	} else {
		log.Print("Err", err)
		return "", "", err
	}
}

func GenerateTokensFromId(id int, acctType enums.AccountType) (*string, *string, error) {
	accessToken, err := GenerateAccessToken(strconv.Itoa(id), string(acctType))
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := GenerateRefreshToken(strconv.Itoa(id), string(acctType))
	if err != nil {
		return nil, nil, err
	}

	return &accessToken, &refreshToken, nil
}

func GenerateTokenPair(userId int, tokenType enums.TokenType) (string, string, error) {
	accessToken, err := GenerateAccessToken(strconv.Itoa(userId), string(tokenType))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken(strconv.Itoa(userId), string(tokenType))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
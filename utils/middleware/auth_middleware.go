package middleware

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/johnyeocx/usual/server2/db/models"
	"github.com/johnyeocx/usual/server2/db/user_db"
	"github.com/johnyeocx/usual/server2/utils/enums/TokenType"
	"github.com/johnyeocx/usual/server2/utils/secure"
)

type contextKey struct {
	key string
}
var UserCtxKey = contextKey{
	key: "user_id"}

var UserTypeCtxKey = contextKey{
	key: "user_type"}
	
func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		var accessToken string
		accessToken, err := c.Cookie("access_token")
		if err != nil || len(accessToken) == 0 {
			const BEARER_SCHEMA = "Bearer "
			authHeader := c.GetHeader("Authorization")
		
			if authHeader == "" || len(authHeader) < len("Bearer  "){
				c.Next()
				return
			}

			accessToken = authHeader[len(BEARER_SCHEMA):]
		}

		userId, userType, err := secure.ParseAccessToken(accessToken)
		if err != nil {
			c.Next();
			return
		}

		c.Set(UserCtxKey.key, userId)
		c.Set(UserTypeCtxKey.key, userType)
		c.Next()
	}
}

func UserCtx(c *gin.Context) (interface{}, interface{}, error) {
	userID, exists := c.Get(UserCtxKey.key)

	if !exists {
		err := models.RequestError{
			StatusCode: http.StatusUnauthorized, 
			Err: errors.New("user id not found in context"),
		}

		return nil, nil, err.Err
	}

	userType, exists := c.Get(UserTypeCtxKey.key)

	if !exists {
		err := models.RequestError{
			StatusCode: http.StatusUnauthorized, 
			Err: errors.New("user type not found in context"),
		}
		return nil, nil, err.Err
	}
	return userID, userType,  nil
}

func AuthenticateUser(c *gin.Context, sqlDB *sql.DB) (int, error) {

	userId, tokenType, err := UserCtx(c)

	if err != nil {
		log.Println("Failed to get user id from context")
		c.JSON(http.StatusUnauthorized, err)
		return -1, err
	}
	
	if tokenType != string(TokenType.User) {
		log.Println("Invalid token type")
		err := errors.New("wrong type token")
		c.JSON(http.StatusUnauthorized, err)
		return -1, err
	}

	userIdInt, err := strconv.Atoi(userId.(string))
	if err != nil {
		log.Println("Invalid user id: str to int err")
		c.JSON(http.StatusUnauthorized, err)
		return -1, err
	}
	
	if ok := user_db.ValidateUser(sqlDB, userIdInt); !ok {
		log.Println("Invalid user id: DB err")
		err := errors.New("invalid customer id")
		c.JSON(http.StatusUnauthorized, err)
		return -1, err
	}

	return userIdInt, nil
}

// func AuthenticateBId(c *gin.Context, sqlDB *sql.DB) (*int, error) {

// 	businessId, userType, err := UserCtx(c)
	
// 	if err != nil {
// 		return nil, err
// 	}

// 	if userType != enums.BusUserType {
// 		return nil, errors.New("wrong type token")
// 	}

// 	businessIdInt, err := strconv.Atoi(businessId.(string))
// 	if err != nil {
// 		return nil, err
// 	}
	
// 	// if ok := db.ValidateBusinessId(sqlDB, businessIdInt); !ok {
// 	// 	return nil, fmt.Errorf("invalid business id")
// 	// }

// 	return &businessIdInt, nil
// }
package middlewares

import (
	"github.com/Deeksharma/taskmanager/internal/config"
	"github.com/Deeksharma/taskmanager/internal/errors"
	"github.com/Deeksharma/taskmanager/internal/log"
	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Claims struct {
	UserName string `json:"user_name"`
	UserType string `json:"user_type"`
	UserId   string `json:"user_id"`
	jwt.StandardClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var header = c.Request.Header.Get("Authorization")
		authToken := strings.Split(header, "Bearer ")
		if header == "" || len(authToken) < 2 {
			err := errors.ErrAuthHeaderNotPresent
			log.InfoWithFields(c, map[string]interface{}{
				"error": err.Error(),
			}, "unauthorized")
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
			return
		}
		var tokenString = authToken[1]

		var keyfunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetString("server.secret_key")), nil
		}

		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tokenString, claims, keyfunc)

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				log.InfoWithFields(c, map[string]interface{}{
					"error": err.Error(),
				}, "unauthorized")
				c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
				return
			}
			log.InfoWithFields(c, map[string]interface{}{
				"error": err.Error(),
			}, "something went wrong!")
			c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if !tkn.Valid {
			err = errors.ErrInvalidToken
			log.InfoWithFields(c, map[string]interface{}{
				"error": err.Error(),
			}, "invalid token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
			return
		}
		log.InfoWithFields(c, map[string]interface{}{
			"user_name": claims.UserName,
			"user_id":   claims.UserId,
			"user_type": claims.UserType,
		}, "JWT claims")

		c.Set("user_name", &claims.UserName)
		c.Set("user_id", &claims.UserId)
		c.Set("user_type", &claims.UserType)

		c.Next()
	}
}

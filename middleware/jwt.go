package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTAuth gin middleware，check jwt token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "请求头未携带token，无权限访问",
			})
			c.Abort()
			return
		}

		j := NewJWT()
		// parseToken
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == ErrTokenExpired {
				c.JSON(http.StatusOK, gin.H{
					"status": -1,
					"msg":    "授权已过期",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)
	}
}

// JWT  the jwt main struct and includ sign key
type JWT struct {
	SigningKey []byte
}

var (
	//ErrTokenExpired Token is expired
	ErrTokenExpired error = errors.New("Token is expired")
	//ErrTokenNotValidYet Token not Valid yet
	ErrTokenNotValidYet error = errors.New("Token not Valid yet")
	//ErrTokenMalformed Token format is incorrect
	ErrTokenMalformed error = errors.New("Token format is incorrect")
	//ErrTokenInvalid   Token is invalid
	ErrTokenInvalid error = errors.New("Token is invalid")
	//SignKey secret string
	SignKey string = "newtrekWang"
)

//CustomClaims jwt playload
type CustomClaims struct {
	ID    string `json:"userId"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	jwt.StandardClaims
}

//NewJWT  New a JWT instance
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

// GetSignKey get sign key for  anothe pkg
func GetSignKey() string {
	return SignKey
}

// SetSignKey set sign key for anothe pkg
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

// CreateToken  Create Token with claims
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

//ParseToken Parse Token from tokenString
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}

// RefreshToken Refresh Token by a jwt token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", ErrTokenInvalid
}

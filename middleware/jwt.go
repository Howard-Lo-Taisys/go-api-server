package middleware

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	ErrTokenExpired     error  = errors.New("TokenIsExpired")
	ErrTokenNotValidYet error  = errors.New("TokenNotActiveYet")
	ErrTokenMalformed   error  = errors.New("NotAToken")
	ErrTokenInvalid     error  = errors.New("TokenInvelid")
	ErrVerificationFail error  = errors.New("VerificationFailed")
	SignKey             string = "Howard.Lo"
)

type JWT struct {
	SigningKey []byte
}

type CustomClaims struct {
	UserName string `gorm:"size:255;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"password"`
	jwt.StandardClaims
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(SignKey),
	}
}

func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
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

	return nil, ErrVerificationFail

}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "Token is null in Header",
				"data":   nil,
			})
			c.Abort()
			return
		}

		j := NewJWT()

		claims, err := j.ParseToken(token)

		if err != nil {
			if err == ErrTokenExpired {
				c.JSON(http.StatusOK, gin.H{
					"status": -1,
					"msg":    "Token is expired",
					"data":   nil,
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
				"data":   nil,
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)
	}
}

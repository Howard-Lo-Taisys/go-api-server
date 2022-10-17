package controller

import (
	"fmt"
	"go-api-server/config"
	"go-api-server/middleware"
	"go-api-server/models"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Register(c *gin.Context) {
	input := models.RegReq{}
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user = models.User{
		Id:       uuid.NewString(),
		UserName: input.UserName,
		Password: input.Password,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Id": &user.Id, "UserName": &user.UserName})
}

func LoginCheck(login models.LoginReq) (bool, models.User, error) {
	userData := models.User{}
	userExist := false

	var user models.User
	dbErr := config.DB.Where("user_name = ?", login.UserName).Find(&user).Error

	if dbErr != nil {
		return userExist, userData, dbErr
	}
	if login.UserName == user.UserName && login.Password == user.Password {
		userExist = true
		userData.UserName = user.UserName
		userData.Password = user.Password
	}

	if !userExist {
		return userExist, userData, fmt.Errorf("login failed")
	}
	return userExist, userData, nil
}

func Login(c *gin.Context) {
	var loginReq models.LoginReq
	if c.ShouldBindJSON(&loginReq) == nil {
		isPass, user, err := LoginCheck(loginReq)
		if isPass {
			generateToken(c, user)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "Verification failed:" + err.Error(),
				"data":   nil,
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "User data parsing failed.",
			"data":   nil,
		})
	}
}

func generateToken(c *gin.Context, user models.User) {

	j := middleware.NewJWT()

	claims := middleware.CustomClaims{
		UserName: user.UserName,
		Password: user.Password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Unix() + 3600),
			Issuer:    "Howard.Lo",
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
			"data":   nil,
		})
	}

	data := models.LoginResult{
		UserName: user.UserName,
		Token:    token,
	}

	log.Println(token)

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "login success.",
		"data":   data,
	})
}

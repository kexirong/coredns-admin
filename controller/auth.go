package controller

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/middleware"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	var conf = config.Get()

	username := c.FormValue("username")
	password := c.FormValue("password")

	if username != conf.Username || password != conf.Password {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": "username or password is incorrect"})
	}

	claims := &middleware.JWTCustomClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(middleware.SigningKey))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"reason": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token": t,
		"token_type":   "Bearer",
	})

}

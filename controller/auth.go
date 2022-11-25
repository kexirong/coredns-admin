package controller

import (
	"net/http"

	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/middleware"
	"github.com/labstack/echo/v4"
)

var jwt = middleware.NewJWT()

func Login(c echo.Context) error {
	var conf = config.Get()

	username := c.FormValue("username")
	password := c.FormValue("password")

	if username != conf.Username || password != conf.Password {
		return c.JSON(http.StatusBadRequest, echo.Map{"reason": "username or password is incorrect"})
	}

	token, err := jwt.CreateToken(middleware.CustomClaims{Username: username})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"reason": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

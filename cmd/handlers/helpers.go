package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(ctx echo.Context, status int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(status)
	err := t.Render(context.Background(), ctx.Response().Writer)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "failed to render response template")

	}
	return nil
}

func readCookie(c echo.Context, ck string) (string, error) {
	cookie, err := c.Cookie(ck)
	if err != nil {
		return "", nil
	}
	return cookie.Value, nil
}

func setCookie(c echo.Context, name, value string) echo.Context {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	c.SetCookie(cookie)

	return c

}

func getIdCookie(c echo.Context) (string, error) {
	id, err := readCookie(c, "id")
	if err != nil {
		log.Println("Error reading id cookie:", err)
	}
	return id, err
}

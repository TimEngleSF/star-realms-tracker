package main

import (
	"html/template"
	"io"

	"github.com/TimEngleSF/star-realms-score-keeper/cmd/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

// Add a Renderer method that satisfies Echo's Renderer interface
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/public", "public")

	e.GET("/", handlers.HandleIndexPage)

	/* ADD PLAYERS */
	e.POST("players", handlers.HandleAddPlayers)

	/* CURRENT PLAYER */
	e.POST("current", handlers.HandleSelectFirstPlayer)

	e.PUT("current", handlers.HandleUpdateCurrPlayer)

	/* RESET GAME ENDPOINT*/
	e.PUT("reset", handlers.HandleResetGame)

	/* SCORE ENDPOINTS */
	e.PUT("score", handlers.HandleUpdateScore)

	/* LAUNCH SERVER */
	e.Logger.Fatal(e.Start(":8080"))

}

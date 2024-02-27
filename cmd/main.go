package main

import (
	"context"
	"html/template"
	"io"
	"net/http"

	"github.com/TimEngleSF/star-realms-score-keeper/cmd/game"
	"github.com/TimEngleSF/star-realms-score-keeper/cmd/handlers"
	"github.com/TimEngleSF/star-realms-score-keeper/views"
	"github.com/a-h/templ"
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

// func newTemplate() *Template {
// 	return &Template{
// 		templates: template.Must(template.ParseGlob("views/*.html")),
// 	}
// }

var tempID int

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/public", "public")

	// Set up templates for the templates in views folder
	// t := newTemplate()
	// e.Renderer = t

	/* INITIALIZE INSTANCE */
	instance := game.InstanceState{Errors: make(map[string]string)}

	/* INITIALIZE GAME */

	var Game game.Game

	// Game := game.Game{
	// 	Players: []game.Player{
	// 		{Id: 0, Name: "Lily", Authority: 1, IsCurrent: true},
	// 		{Id: 1, Name: "Kara", Authority: 50, IsCurrent: false},
	// 	},
	// }
	// Game.Current = &Game.Players[0]
	// instance.Game = &Game

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Clear all errors at the start of each request
			instance.Errors = make(map[string]string)
			return next(c)
		}
	})

	e.GET("/", handlers.HandleIndexPage)

	/* ADD PLAYERS */
	e.POST("players", handlers.HandleAddPlayers)

	/* CURRENT PLAYER */
	e.POST("current", handlers.HandleSelectFirstPlayer)

	e.PUT("current", handlers.HandleUpdateCurrPlayer)

	/* RESET GAME ENDPOINT*/
	e.PUT("reset", func(c echo.Context) error {
		Game.Restart()
		return render(c, http.StatusContinue, views.NewGameForm())
	})

	/* SCORE ENDPOINTS */
	e.PUT("score", handlers.HandleUpdateScore(&instance))

	/* LAUNCH SERVER */
	e.Logger.Fatal(e.Start(":8080"))

}
func render(ctx echo.Context, status int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(status)
	err := t.Render(context.Background(), ctx.Response().Writer)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "failed to render response template")

	}
	return nil
}

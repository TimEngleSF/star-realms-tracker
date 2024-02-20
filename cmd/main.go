package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

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

func newTemplate() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/css", "css")

	// Set up templates for the templates in views folder
	t := newTemplate()
	e.Renderer = t

	/* INITIALIZE GAME */
	var Game Game

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world!")
	})

	e.GET("game", func(c echo.Context) error {
		return c.Render(http.StatusOK, "base", Game)
	})

	e.POST("players", func(c echo.Context) error {
		if len(Game.Players) < 2 {
			for i := 0; i < 2; i++ {
				var player Player
				player.Id = i
				player.Authority = 50
				player.Name = c.FormValue(fmt.Sprintf("player%d-name", i))
				Game.Players.AddPlayer(player)
			}
		}
		return c.Render(http.StatusAccepted, "select-current", Game)
	})

	e.POST("current", func(c echo.Context) error {
		sp, err := strconv.Atoi(c.FormValue("player-radio"))

		if err != nil {
			log.Println("Error parsing int from player-radio value")
			Game.Current = &Game.Players[0]
		} else {
			Game.Current = &Game.Players[sp]
		}
		// TODO: Display scoreboard
		return c.Render(201, "base", Game)
		// Game.Current = &Game.Players[]
	})
	e.PUT("reset", func(c echo.Context) error {
		Game.Restart()
		return c.Render(http.StatusContinue, "new-game-form", Game)
	})
	e.Logger.Fatal(e.Start(":8080"))
}

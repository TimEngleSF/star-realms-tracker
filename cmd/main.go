package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

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

	// Game := Game{
	// 	Players: []Player{
	// 		{Id: 0, Name: "Lily", Authority: 1, IsCurrent: true},
	// 		{Id: 1, Name: "Kara", Authority: 50, IsCurrent: false},
	// 	},
	// }
	// Game.Current = &Game.Players[0]
	instance.Game = &Game
	instance.Id = tempID

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Clear all errors at the start of each request
			instance.Errors = make(map[string]string)
			return next(c)
		}
	})

	e.GET("/", handlers.HandleIndexPage(&instance))

	e.POST("players", func(c echo.Context) error {
		if len(Game.Players) < 2 {
			for i := 0; i < 2; i++ {
				var player game.Player
				player.Id = i
				player.Authority = 50
				player.Name = c.FormValue(fmt.Sprintf("player%d-name", i))
				Game.Players.AddPlayer(player)
			}
		}
		return render(c, http.StatusAccepted, views.SelectCurrentPlayerTemplate(Game.Players))
		// return c.Render(http.StatusAccepted, "select-current", instance)
	})

	/* CURRENT PLAYER */
	e.POST("current", handlers.HandleSelectFirstPlayer(&instance))

	e.PUT("current", func(c echo.Context) error {

		for i := range Game.Players {
			p := &Game.Players[i]
			p.IsCurrent = !p.IsCurrent
			if p.IsCurrent {
				Game.Current = p
			}
		}
		return render(c, http.StatusContinue, views.ScoreboardTemplate(Game))
		// return c.Render(http.StatusContinue, "scoreboard", instance)
	})

	/* RESET GAME ENDPOINT*/
	e.PUT("reset", func(c echo.Context) error {
		Game.Restart()
		return render(c, http.StatusContinue, views.NewGameForm())
		// return c.Render(http.StatusContinue, "new-game-form", instance)
	})

	/* SCORE ENDPOINTS */
	e.PUT("score", func(c echo.Context) error {
		// Query ID and action
		id, err := strconv.Atoi(c.QueryParam("player"))
		scoreAction := c.QueryParam("action")

		if err != nil {
			log.Printf("Error parsing id for PUT '/score': %v\n", err)
			c.Response().Header().Set(
				"HX-Trigger",
				`{"error": {"id": "scoreboard-error-msg", "message":  "Error updating score: Invalid player ID"}}`,
			)
			return render(c, http.StatusInternalServerError, views.ScoreboardTemplate(Game))
			// return c.Render(500, "scoreboard", instance)
		}

		// Get queried player and their score
		player := &Game.Players[id]
		score := &player.Authority
		// Query score action

		// Edit player's score
		if scoreAction == "add" {
			player.IncrementAuthority()
		} else if scoreAction == "subtract" {
			player.DecrementAuthority()
		}
		// When user score is zero
		if *score == 0 {
			instance.Game.Loser = player
			winnerId := 0
			if id == 0 {
				winnerId = 1
			}
			instance.Game.Winner = &Game.Players[winnerId]
			instance.Game.Complete = true
			return render(c, http.StatusOK, views.WinnerTemplate(Game))
			// return c.Render(201, "winner-display", instance)
		}
		// Render updated scores
		return render(c, http.StatusOK, views.ScoreboardTemplate(Game))
		// return c.Render(200, "scoreboard", instance)
	})
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

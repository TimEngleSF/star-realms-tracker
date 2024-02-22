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

type InstanceState struct {
	Game   *Game
	Errors map[string]string
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/public", "public")

	// Set up templates for the templates in views folder
	t := newTemplate()
	e.Renderer = t

	/* INITIALIZE INSTANCE */
	instance := InstanceState{Errors: make(map[string]string)}

	/* INITIALIZE GAME */

	var Game Game

	//Game := Game{
	//	Players: []Player{
	//		{Id: 0, Name: "Lily", Authority: 1},
	//		{Id: 1, Name: "Kara", Authority: 50},
	//	},
	//}
	//Game.Current = &Game.Players[0]
	instance.Game = &Game

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Clear all errors at the start of each request
			instance.Errors = make(map[string]string)
			return next(c)
		}
	})

	e.GET("/", func(c echo.Context) error {

		for _, p := range Game.Players {
			fmt.Println(p.Authority)
		}

		return c.Render(http.StatusOK, "base", instance)
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
		return c.Render(http.StatusAccepted, "select-current", instance)
	})

	/* CURRENT PLAYER */
	// TODO: Set isCurrent on player to true
	e.POST("current", func(c echo.Context) error {
		sp, err := strconv.Atoi(c.FormValue("player-radio"))
		if err != nil {
			log.Println("Error parsing int from player-radio value")
			Game.Current = &Game.Players[0]
		} else {
			Game.Current = &Game.Players[sp]
		}
		Game.Current.isCurrent = true
		return c.Render(201, "scoreboard", instance)
	})

	e.PUT("current", func(c echo.Context) error {

		for i := range Game.Players {
			p := &Game.Players[i]
			p.isCurrent = !p.isCurrent
			if p.isCurrent {
				Game.Current = p
			}
		}

		return c.Render(http.StatusContinue, "scoreboard", instance)
	})

	/* RESET GAME ENDPOINT*/
	e.PUT("reset", func(c echo.Context) error {
		Game.Restart()
		return c.Render(http.StatusContinue, "new-game-form", instance)
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
			return c.Render(500, "scoreboard", instance)
		}

		// Get queried player and their score
		player := &Game.Players[id]
		score := &player.Authority
		// Query score action

		// Edit player's score
		if scoreAction == "add" {
			player.incrementAuthority()
		} else if scoreAction == "subtract" {
			player.decrementAuthority()
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
			return c.Render(201, "winner-display", instance)
		}
		// Render updated scores
		return c.Render(200, "scoreboard", instance)
	})
	e.Logger.Fatal(e.Start(":8081"))
}

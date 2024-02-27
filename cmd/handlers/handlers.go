package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/TimEngleSF/star-realms-score-keeper/cmd/game"
	"github.com/TimEngleSF/star-realms-score-keeper/views"
	"github.com/labstack/echo/v4"
)

/* INDEX SCREEN */
func HandleIndexPage(c echo.Context) error {
	// TODO: Set up in memory store for gamess
	// TODO: Install uuid
	// TODO: Setup supabase
	var i game.InstanceState
	id, _ := getIdCookie(c)

	if id == "" {
		i = game.NewInstance()
		game.InstancesInMemory = append(game.InstancesInMemory, i)
		c = setCookie(c, "id", i.Id)
	} else {
		var err error
		i, err = game.InstancesInMemory.GetInstanceById(id)

		if err != nil {
			fmt.Println("There was an error getting InstanceById", err)
			// FIXME: Is there a better way to handle this error than just creating a new Instance?
			i = game.NewInstance()
			game.InstancesInMemory = append(game.InstancesInMemory, i)
			c = setCookie(c, "id", i.Id)
		}
	}

	return render(c, http.StatusOK, views.Index(*i.Game))
}

/* ADD PLAYERS SCREEN */
func HandleAddPlayers(ins *game.InstanceState) echo.HandlerFunc {
	return func(c echo.Context) error {
		g := ins.Game
		if len(g.Players) < 2 {
			for i := 0; i < 2; i++ {
				var player game.Player
				player.Id = i
				player.Authority = 50
				player.Name = c.FormValue(fmt.Sprintf("player%d-name", i))
				g.Players.AddPlayer(player)
			}
		}
		return render(c, http.StatusAccepted, views.SelectCurrentPlayerTemplate(g.Players))
	}
}

/* SELECT PLAYER SCREEN */
func HandleSelectFirstPlayer(ins *game.InstanceState) echo.HandlerFunc {
	g := ins.Game

	return func(c echo.Context) error {
		sp, err := strconv.Atoi(c.FormValue("player-radio"))
		if err != nil {
			log.Println("Error parsing int from player-radio value")
			g.Current = &g.Players[0]
		} else {
			g.Current = &g.Players[sp]
		}
		g.Current.IsCurrent = true
		return render(c, 201, views.ScoreboardTemplate(*g))
	}
}

/* SCOREBOARD SCREEN */

func HandleUpdateCurrPlayer(ins *game.InstanceState) echo.HandlerFunc {
	g := ins.Game
	return func(c echo.Context) error {
		g.Players.ResetAuthorityDifference()
		for i := range g.Players {
			p := &g.Players[i]
			p.IsCurrent = !p.IsCurrent
			if p.IsCurrent {
				g.Current = p
			}
		}
		return render(c, http.StatusContinue, views.ScoreboardTemplate(*g))
	}
}

func HandleUpdateScore(ins *game.InstanceState) echo.HandlerFunc {
	g := ins.Game

	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.QueryParam("player"))
		scoreAction := c.QueryParam("action")

		if err != nil {
			log.Printf("Error parsing id for PUT '/score': %v\n", err)
			c.Response().Header().Set(
				"HX-Trigger",
				`{"error": {"id": "scoreboard-error-msg", "message":  "Error updating score: Invalid player ID"}}`,
			)
			return render(c, http.StatusInternalServerError, views.ScoreboardTemplate(*g))
		}

		player := &g.Players[id]
		score := &player.Authority

		if scoreAction == "add" {
			player.IncrementAuthority()
		} else if scoreAction == "subtract" {
			player.DecrementAuthority()
		}

		if *score == 0 {
			g.Loser = player
			winnerId := 0
			if id == 0 {
				winnerId = 1
			}
			g.Winner = &g.Players[winnerId]
			g.Complete = true
			return render(c, http.StatusOK, views.WinnerTemplate(*g))
		}

		return render(c, http.StatusOK, views.ScoreboardTemplate(*g))
	}
}

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
func HandleIndexPage(ins *game.InstanceState) echo.HandlerFunc {
	g := ins.Game
	return func(c echo.Context) error {
		// TODO: Set up in memory store for gamess
		// TODO: Install uuid
		// TODO: Setup supabase
		hasCookie := false
		id, err := readCookie(c, "id")
		if err != nil {
			log.Println("Error reading id cookie:", err)
		}
		if id != "" {
			hasCookie = true
		}

		if !hasCookie {
			// TODO: SSSetup a new Instance with new Game
			// TODO: add ID to Instance type
			ins.Id++
			c = setCookie(c, "id", strconv.Itoa(ins.Id))
		}

		return render(c, http.StatusOK, views.Index(*g))
	}
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

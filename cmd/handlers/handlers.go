package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/TimEngleSF/star-realms-score-keeper/cmd/game"
	"github.com/TimEngleSF/star-realms-score-keeper/views"
	"github.com/labstack/echo/v4"
)

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

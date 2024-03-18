package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/TimEngleSF/star-realms-score-keeper/cmd/game"
	"github.com/TimEngleSF/star-realms-score-keeper/views"
	"github.com/google/uuid"
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
		newId := uuid.NewString()
		c = setCookie(c, "id", newId)

		i = game.NewInstance(newId)
		game.InstancesInMemory = append(game.InstancesInMemory, i)
	} else {
		var err error
		i, err = game.InstancesInMemory.GetInstanceById(id)

		if err != nil {
			fmt.Println("There was an error getting InstanceById", err)
			// FIXME: Is there a better way to handle this error than just creating a new Instance?
			newId := uuid.NewString()
			i = game.NewInstance(newId)
			game.InstancesInMemory = append(game.InstancesInMemory, i)
			c = setCookie(c, "id", i.Id)
		}
	}

	return render(c, http.StatusOK, views.Index(*i.Game))
}

/* RESET GAME */
func HandleResetGame(c echo.Context) error {
	c, i, err := getInstance(c)
	if err != nil {
		c.Response().Header().Set(
			"HX-Trigger",
			`{"error": {"id": "scoreboard-error-msg", "message":  "Error getting user instance"}}`,
		)
	}
	g := i.Game
	g.Restart()
	i.Game = g
	return render(c, http.StatusAccepted, views.NewGameForm())

}

/* ADD PLAYERS SCREEN */
func HandleAddPlayers(c echo.Context) error {
	c, i, err := getInstance(c)
	if err != nil {
		c.Response().Header().Set(
			"HX-Trigger",
			`{"error": {"id": "scoreboard-error-msg", "message":  "Error getting user instance"}}`,
		)
	}
	g := i.Game
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

/* SELECT PLAYER SCREEN */
func HandleSelectFirstPlayer(c echo.Context) error {
	c, i, err := getInstance(c)
	if err != nil {
		c.Response().Header().Set(
			"HX-Trigger",
			`{"error": {"id": "scoreboard-error-msg", "message":  "Error getting user instance"}}`,
		)
	}
	g := i.Game

	// Set start time of first turn
	g.GameDuration.TurnStartTime = time.Now()

	// Get first player id
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

/* SCOREBOARD SCREEN */

func HandleUpdateCurrPlayer(c echo.Context) error {
	c, i, err := getInstance(c)

	if err != nil {
		c.Response().Header().Set(
			"HX-Trigger",
			`{"error": {"id": "scoreboard-error-msg", "message":  "Error getting user instance"}}`,
		)
	}
	g := i.Game

	g.Players.ResetAuthorityDifference()

	// Calculate completed turn duration
	gDuration := g.GameDuration
	dur := time.Since(gDuration.TurnStartTime)

	// Update GameDuration fields
	gDuration.PrevTurnDuration = dur
	gDuration.TotalDuration += dur

	// Iterate over players
	for n := range g.Players {
		// Select first player
		p := &g.Players[n]

		// Add Turn duration for completed turn's current player
		if p.IsCurrent {
			t := p.TurnsDuration + dur
			p.TurnsDuration = t
		}
		// Reverse p.IsCurrent boolean
		p.IsCurrent = !p.IsCurrent
		if p.IsCurrent {
			g.Current = p
		}
	}
	newTurnTime := time.Now()
	gDuration.TurnStartTime = newTurnTime

	return render(c, http.StatusContinue, views.ScoreboardTemplate(*g))
}

func HandleUpdateScore(c echo.Context) error {
	c, i, err := getInstance(c)

	if err != nil {
		c.Response().Header().Set(
			"HX-Trigger",
			`{"error": {"id": "scoreboard-error-msg", "message":  "Error getting user instance"}}`,
		)
	}
	g := i.Game

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

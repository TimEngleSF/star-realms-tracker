package game

import (
	"errors"
	"time"
)

var InstancesInMemory = Instances{}

type Instances []InstanceState

type InstanceState struct {
	Game   *Game
	Errors map[string]string
	Id     string
}

type GameDuration struct {
	TurnStartTime    time.Time
	PrevTurnDuration time.Duration
	TotalDuration    time.Duration
}

type Game struct {
	Players       Players
	Current       *Player
	Winner        *Player
	Loser         *Player
	Complete      bool
	Date          *time.Time
	CurrTurnScore *CurrTurnScore
	GameDuration  *GameDuration
}

type CurrTurnScore struct {
	currPlayer  int
	otherPlayer int
}

type Player struct {
	Id                  int
	Name                string
	Authority           int
	IsCurrent           bool
	AuthorityDifference int
	TurnsDuration       time.Duration
}

type Players []Player

/* PLAYER METHODS */
func (p *Player) IncrementAuthority() {
	p.Authority++
	p.AuthorityDifference++
}

func (p *Player) DecrementAuthority() {
	p.Authority--
	p.AuthorityDifference--
}

/* PLAYERS METHODS */
func (ps *Players) AddPlayer(p Player) {
	*ps = append(*ps, p)
}

func (ps *Players) ResetAuthorityDifference() {
	for i := range *ps {
		(*ps)[i].AuthorityDifference = 0
	}
}

/* GAME METHODS */
func (g *Game) SwitchCurrentPlayer() {
	if g.Current == &g.Players[0] {
		g.Current = &g.Players[1]
	} else {
		g.Current = &g.Players[0]
	}
}

func (g *Game) Restart() {
	g.Players = Players{}
	g.Current = nil
	g.Winner = nil
	g.Loser = nil
	g.Complete = false
	g.Date = nil
	g.GameDuration = &GameDuration{}
}

// TODO: In the future a db will be used and the pointer will no longer be necessary
// The Games type and GamesInMemory below will also be unnesecary
func NewGame() Game {
	var g Game
	g.Players = Players{}
	g.Complete = false
	g.GameDuration = &GameDuration{}
	return g
}
func NewInstance(id string) InstanceState {
	g := NewGame()
	var i = InstanceState{}
	i.Game = &g
	i.Id = id
	return i
}

func (is Instances) GetInstanceById(id string) (InstanceState, error) {
	var ti InstanceState
	for _, i := range is {
		if i.Id == id {
			ti = i
			break
		}
	}
	if ti.Id == "" {
		err := errors.New("Could not locate Client Instance")
		return ti, err
	}
	return ti, nil
}

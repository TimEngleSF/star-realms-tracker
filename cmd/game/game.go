package game

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type InstanceState struct {
	Game   *Game
	Errors map[string]string
	Id     string
}

type Game struct {
	Players       Players
	Current       *Player
	Winner        *Player
	Loser         *Player
	Complete      bool
	Date          *time.Time
	CurrTurnScore *CurrTurnScore
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
}

// TODO: In the future a db will be used and the pointer will no longer be necessary
// The Games type and GamesInMemory below will also be unnesecary
func NewGame() Game {
	var g Game
	g.Players = Players{}
	g.Complete = false

	return g
}
func NewInstance() InstanceState {
	g := NewGame()
	var i = InstanceState{}
	i.Game = &g
	i.Id = uuid.NewString()
	return i
}

type Instances []InstanceState

var InstancesInMemory = Instances{}

func (is *Instances) GetInstanceById(id string) (InstanceState, error) {
	var ti InstanceState
	for _, i := range *is {
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

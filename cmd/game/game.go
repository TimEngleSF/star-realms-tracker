package game

import "time"

type InstanceState struct {
	Game   *Game
	Errors map[string]string
	Id     int
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

func (p *Player) IncrementAuthority() {
	p.Authority++
	p.AuthorityDifference++
}

func (p *Player) DecrementAuthority() {
	p.Authority--
	p.AuthorityDifference--
}

func (ps *Players) AddPlayer(p Player) {
	*ps = append(*ps, p)
}

func (g *Game) SwitchCurrentPlayer() {
	if g.Current == &g.Players[0] {
		g.Current = &g.Players[1]
	} else {
		g.Current = &g.Players[0]
	}
}

func (ps *Players) ResetAuthorityDifference() {
	for i, _ := range *ps {
		(*ps)[i].AuthorityDifference = 0
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

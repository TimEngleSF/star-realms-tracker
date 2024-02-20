package main

import "time"

type Game struct {
	Players  Players
	Current  *Player
	Winner   *Player
	Loser    *Player
	Complete bool
	Date     *time.Time
}

type Player struct {
	Id        int
	Name      string
	Authority int
}

type Players []Player

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

func (g *Game) Restart() {
	g.Players = Players{}
	g.Current = nil
	g.Winner = nil
	g.Loser = nil
	g.Complete = false
	g.Date = nil
}

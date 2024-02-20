package main

type Player struct {
	Id        int
	Name      string
	Authority int
}

type Players []Player

func (ps *Players) AddPlayer(p Player) {
	*ps = append(*ps, p)
}

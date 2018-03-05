package game

import "fmt"

type Monster struct {
	Pos
	Rune      rune
	Name      string
	Hitpoints int
	Strength  int
	Speed     float64
}

func NewRat(p Pos) *Monster {
	return &Monster{p, 'R', "Rat", 5, 5, 2.0}
}

func NewSpider(p Pos) *Monster {
	return &Monster{p, 'S', "Spider", 10, 10, 1.0}
}

func (m *Monster) Update(level *Level) {
	playerPos := level.Player.Pos
	positions := level.astar(m.Pos, playerPos)
	// Most be > 1 because the 1st position is the monsters current
	if len(positions) > 1 {
		fmt.Println("Move")
		m.Move(positions[1], level)
	}
}

func (m *Monster) Move(to Pos, level *Level) {
	_, exists := level.Monsters[to]
	// TODO check if the tile being moved to is valid
	// TODO if player is in the way, attack player
	if !exists && to != level.Player.Pos {
		delete(level.Monsters, m.Pos)
		level.Monsters[to] = m
		m.Pos = to
	}
}

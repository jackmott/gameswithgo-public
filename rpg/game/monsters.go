package game

import "fmt"

type Monster struct {
	Character
}

func NewRat(p Pos) *Monster {
	//TODO learn why our syntax was no good here
	//return &Monster{X: p.X, Y: p.Y, Rune: 'R', Name: "Rat", Hitpoints: 5, Strength: 5, Speed: 1.5, ActionPoints: 0.0}
	monster := &Monster{}
	monster.Pos = p
	monster.Rune = 'R'
	monster.Name = "Rat"
	monster.Hitpoints = 500
	monster.Strength = 0
	monster.Speed = 2.0
	monster.ActionPoints = 0.0
	return monster
}

func NewSpider(p Pos) *Monster {
	monster := &Monster{}
	monster.Pos = p
	monster.Rune = 'S'
	monster.Name = "Spider"
	monster.Hitpoints = 1000
	monster.Strength = 0
	monster.Speed = 1.0
	monster.ActionPoints = 0.0
	return monster
}

func (m *Monster) Update(level *Level) {
	m.ActionPoints += m.Speed
	playerPos := level.Player.Pos
	apInt := int(m.ActionPoints)
	positions := level.astar(m.Pos, playerPos)
	moveIndex := 1
	fmt.Println("apInt:", apInt)
	for i := 0; i < apInt; i++ {
		fmt.Println("moveIndeX", moveIndex, "pos len:", len(positions))
		if moveIndex < len(positions) {
			m.Move(positions[moveIndex], level)
			moveIndex++
			m.ActionPoints--
		}
	}
}

func (m *Monster) Move(to Pos, level *Level) {
	fmt.Println("Moving")
	_, exists := level.Monsters[to]
	// TODO check if the tile being moved to is valid
	if !exists && to != level.Player.Pos {
		delete(level.Monsters, m.Pos)
		level.Monsters[to] = m
		m.Pos = to
	} else {
		fmt.Println("attacking")
		Attack(m, level.Player)
		if m.Hitpoints <= 0 {
			delete(level.Monsters, m.Pos)
		}
		if level.Player.Hitpoints <= 0 {
			fmt.Println("you died!")
			panic("ded")
		}
	}
}

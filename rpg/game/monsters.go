package game

type Monster struct {
	Character
}

func NewRat(p Pos) *Monster {
	monster := &Monster{}
	monster.Pos = p
	monster.Rune = 'R'
	monster.Name = "Rat"
	monster.Hitpoints = 5
	monster.Strength = 0
	monster.Speed = 2.0
	monster.ActionPoints = 0.0
	monster.SightRange = 10
	return monster
}

func NewSpider(p Pos) *Monster {
	monster := &Monster{}
	monster.Pos = p
	monster.Rune = 'S'
	monster.Name = "Spider"
	monster.Hitpoints = 1000
	monster.Strength = 5
	monster.Speed = 1.0
	monster.ActionPoints = 0.0
	monster.SightRange = 10
	return monster
}

func (m *Monster) Update(level *Level) {
	m.ActionPoints += m.Speed
	playerPos := level.Player.Pos
	apInt := int(m.ActionPoints)
	positions := level.astar(m.Pos, playerPos)

	//Do we have any path to the goal?
	if len(positions) == 0 {
		m.Pass()
		return
	}

	moveIndex := 1
	for i := 0; i < apInt; i++ {
		if moveIndex < len(positions) {
			m.Move(positions[moveIndex], level)
			moveIndex++
			m.ActionPoints--
		}
	}
}

func (m *Monster) Pass() {
	m.ActionPoints -= m.Speed
}

func (m *Monster) Move(to Pos, level *Level) {
	_, exists := level.Monsters[to]
	// TODO check if the tile being moved to is valid
	if !exists && to != level.Player.Pos {
		delete(level.Monsters, m.Pos)
		level.Monsters[to] = m
		m.Pos = to
		return
	}

	if to == level.Player.Pos {
		level.Attack(&m.Character, &level.Player.Character)
		if m.Hitpoints <= 0 {
			delete(level.Monsters, m.Pos)
		}
		if level.Player.Hitpoints <= 0 {
			panic("ded")
		}
	}

}

package game

type ItemType int

const (
	Weapon ItemType = iota
	Helmet
	Other
)

type Item struct {
	Typ ItemType
	Entity
	power float64
	// TODO
	// weapon - attack bonus
	// armor - armor class
}

func NewSword(p Pos) *Item {
	return &Item{Weapon, Entity{p, "Sword", 's'}, 2.0}
}

func NewHelmet(p Pos) *Item {
	return &Item{Helmet, Entity{p, "Helmet", 'h'}, .5}
}

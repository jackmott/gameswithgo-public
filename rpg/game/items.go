package game

type Item struct {
	Entity
	// TODO
	// weapon - attack bonus
	// armor - armor class
}

func NewSword(p Pos) *Item {
	return &Item{Entity{p, "Sword", 's'}}
}

func NewHelmet(p Pos) *Item {
	return &Item{Entity{p, "Helmet", 'h'}}
}

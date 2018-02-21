package main

import (
	"github.com/jackmott/rpg/game"
	"github.com/jackmott/rpg/ui2d"
)

func main() {
	ui := &ui2d.UI2d{}
	game.Run(ui)
}

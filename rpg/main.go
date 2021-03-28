package main

import (
	"github.com/LucasK1/gameswithgo/rpg/game"
	"github.com/LucasK1/gameswithgo/rpg/ui2d"
)

func main() {
	ui := &ui2d.UI2d{}
	game.Run(ui)
}

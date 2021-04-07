package main

import (
	"github.com/LucasK1/gameswithgo/rpg/game"
	"github.com/LucasK1/gameswithgo/rpg/ui2d"
)

func main() {
	game := game.NewGame(1)

	go func() {
		ui := ui2d.NewUI(game.InputChan, game.LevelChans[0])
		ui.Run()
	}()
	game.Run()

}

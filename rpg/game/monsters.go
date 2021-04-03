package game

import "fmt"

type Monster struct {
	Character
}

func NewRat(pos Pos) *Monster {
	return &Monster{Character: Character{Entity: Entity{Pos: pos, Name: "Rat", Rune: 'R'}, HP: 50, Strength: 5, Speed: 2.0, AP: 0.0}}
}

func NewSpider(pos Pos) *Monster {
	return &Monster{Character: Character{Entity: Entity{Pos: pos, Name: "Spider", Rune: 'S'}, HP: 100, Strength: 10, Speed: 1.0, AP: 0.0}}
}

func (m *Monster) Update(level *Level) {
	m.AP += m.Speed
	playerPos := level.Player.Pos

	apInt := int(m.AP)

	positions := level.astar(m.Pos, playerPos)

	moveIndex := 1
	for i := 0; i < apInt; i++ {
		if moveIndex < len(positions) {
			fmt.Println(positions)
			m.Move(positions[moveIndex], level)
			moveIndex++
			m.AP--
		}
	}
}

func (m *Monster) Move(to Pos, level *Level) {
	_, exists := level.Monsters[to]
	if !exists && to != level.Player.Pos {
		delete(level.Monsters, m.Pos)
		level.Monsters[to] = m
		m.Pos = to
	} else {
		Attack(&m.Character, &level.Player.Character)
		fmt.Println("Player HP:", level.Player.HP)
		fmt.Println("Monster HP:", m.HP)

		if m.HP <= 0 {
			delete(level.Monsters, m.Pos)
		}

		if level.Player.HP <= 0 {
			fmt.Println("You Died")
			panic("You Died")
		}

	}
}

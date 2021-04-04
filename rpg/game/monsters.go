package game

import "fmt"

type Monster struct {
	Character
}

func NewRat(pos Pos) *Monster {
	return &Monster{Character: Character{Entity: Entity{Pos: pos, Name: "Rat", Rune: 'R'}, HP: 500, Strength: 0, Speed: 2.0, AP: 0.0}}
}

func NewSpider(pos Pos) *Monster {
	return &Monster{Character: Character{Entity: Entity{Pos: pos, Name: "Spider", Rune: 'S'}, HP: 1000, Strength: 0, Speed: 1.0, AP: 0.0}}
}

func (m *Monster) Update(level *Level) {
	m.AP += m.Speed
	playerPos := level.Player.Pos

	apInt := int(m.AP)

	positions := level.astar(m.Pos, playerPos)

	if len(positions) == 0 {
		m.Pass()
		return
	}

	moveIndex := 1
	for i := 0; i < apInt; i++ {
		if moveIndex < len(positions) {
			m.Move(positions[moveIndex], level)
			moveIndex++
			m.AP--
		}
	}
}

func (m *Monster) Pass() {
	m.AP -= m.Speed
}

func (m *Monster) Move(to Pos, level *Level) {
	_, exists := level.Monsters[to]
	if !exists && to != level.Player.Pos {
		delete(level.Monsters, m.Pos)
		level.Monsters[to] = m
		m.Pos = to
		return
	}
	if to == level.Player.Pos {

		level.AddEvent(m.Name + " attacks " + level.Player.Name + "!")
		Attack(&m.Character, &level.Player.Character)
		fmt.Println("Player HP:", level.Player.HP)
		fmt.Println("Monster HP:", m.HP)

		if m.HP <= 0 {
			level.AddEvent(level.Player.Name + " killed the " + m.Name)
			delete(level.Monsters, m.Pos)
		}

		if level.Player.HP <= 0 {
			fmt.Println("You Died")
			panic("You Died")
		}
	}

}

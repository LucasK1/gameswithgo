package game

import "fmt"

type Monster struct {
	Pos
	Rune     rune
	Name     string
	Hp       int
	Strength int
	Speed    float64
	AP       float64
}

func NewRat(pos Pos) *Monster {
	return &Monster{pos, 'R', "Rat", 5, 5, 2.0, 0.0}
}

func NewSpider(pos Pos) *Monster {
	return &Monster{pos, 'S', "Spider", 10, 10, 1.0, 0.0}
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
	}
}

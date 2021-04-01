package game

import "fmt"

type Monster struct {
	Pos
	Rune     rune
	Name     string
	Hp       int
	Strength int
	Speed    float64
}

func NewRat(pos Pos) *Monster {
	return &Monster{pos, 'R', "Rat", 5, 5, 2}
}

func NewSpider(pos Pos) *Monster {
	return &Monster{pos, 'S', "Spider", 10, 10, 1}
}

func (m *Monster) Update(level *Level) {
	playerPos := level.Player.Pos
	positions := level.astar(m.Pos, playerPos)
	if len(positions) > 1 {
		fmt.Println(positions)
		m.Move(positions[1], level)
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

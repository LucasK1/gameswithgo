package game

type Monster struct {
	Character
}

func NewRat(pos Pos) *Monster {
	return &Monster{Character: Character{Entity: Entity{Pos: pos, Name: "Rat", Rune: 'R'}, HP: 200, Strength: 0, Speed: 2.0, AP: 0.0, SightRange: 10}}
}

func NewSpider(pos Pos) *Monster {
	return &Monster{Character: Character{Entity: Entity{Pos: pos, Name: "Spider", Rune: 'S'}, HP: 100, Strength: 0, Speed: 1.0, AP: 0.0, SightRange: 10}}
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
		level.Attack(&m.Character, &level.Player.Character)

		if m.HP <= 0 {
			delete(level.Monsters, m.Pos)
		}
		if level.Player.HP <= 0 {
			panic("DEAD")
		}
	}

}

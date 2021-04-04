package game

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Game struct {
	LevelChans []chan *Level
	InputChan  chan *Input
	Level      *Level
}

func NewGame(numWindows int, path string) *Game {
	levelChans := make([]chan *Level, numWindows)
	for i := range levelChans {
		levelChans[i] = make(chan *Level)
	}
	inputChan := make(chan *Input)

	return &Game{LevelChans: levelChans, InputChan: inputChan, Level: loadLevelFromFile(path)}
}

type InputType int

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	QuitGame
	CloseWindow
	Search
)

type Input struct {
	Type         InputType
	LevelChannel chan *Level
}

type Tile rune

const (
	StoneWall  Tile = '#'
	DirtFloor  Tile = '.'
	ClosedDoor Tile = '|'
	OpenDoor   Tile = '/'
	Blank      Tile = 0
	Pending    Tile = -1
)

type Pos struct {
	X, Y int
}

type Entity struct {
	Pos
	Name string
	Rune rune
}

type Character struct {
	Entity
	HP       int
	Strength int
	Speed    float64
	AP       float64
}

type Player struct {
	Character
}

type Level struct {
	Map      [][]Tile
	Player   Player
	Monsters map[Pos]*Monster
	Events   []string
	EventPos int
	Debug    map[Pos]bool
}

type Attackable interface {
	GetAP() float64
	SetAP(float64)
	GetHP() int
	SetHP(int)
	GetAttackPower() int
}

func (c *Character) GetAP() float64 {
	return c.AP
}

func (c *Character) SetAP(ap float64) {
	c.AP = ap
}

func (c *Character) GetHP() int {
	return c.HP
}

func (c *Character) SetHP(hp int) {
	c.HP = hp
}

func (c *Character) GetAttackPower() int {
	return c.Strength
}

func Attack(a1, a2 Attackable) {
	a1.SetAP(a1.GetAP() - 1)
	a2.SetHP(a2.GetHP() - a1.GetAttackPower())
}

func (level *Level) AddEvent(event string) {
	level.Events[level.EventPos] = event
	level.EventPos++
	if level.EventPos == len(level.Events) {
		level.EventPos = 0
	}

}

func loadLevelFromFile(filename string) *Level {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	levelLines := make([]string, 0)
	longestRow := 0
	index := 0
	for scanner.Scan() {
		levelLines = append(levelLines, scanner.Text())
		if len(levelLines[index]) > longestRow {
			longestRow = len(levelLines[index])
		}
		index++
	}
	level := &Level{}

	level.Events = make([]string, 10)

	level.Player.Name = "Dralanor"
	level.Player.Rune = '@'
	level.Player.HP = 20
	level.Player.Strength = 20
	level.Player.Speed = 1
	level.Player.AP = 0

	level.Map = make([][]Tile, len(levelLines))
	level.Monsters = make(map[Pos]*Monster)

	for i := range level.Map {
		level.Map[i] = make([]Tile, longestRow)
	}

	for y, line := range levelLines {
		for x, character := range line {
			var t Tile
			switch character {
			case ' ', '\t', '\n', '\r':
				t = Blank
			case '#':
				t = StoneWall
			case '|':
				t = ClosedDoor
			case '/':
				t = OpenDoor
			case '.':
				t = DirtFloor
			case '@':
				level.Player.X = x
				level.Player.Y = y
				t = Pending
			case 'R':
				level.Monsters[Pos{x, y}] = NewRat(Pos{x, y})
				t = Pending
			case 'S':
				level.Monsters[Pos{x, y}] = NewSpider(Pos{x, y})
				t = Pending
			default:
				panic("Invalid character in map")
			}
			level.Map[y][x] = t
		}
	}

	for y, row := range level.Map {
		for x, tile := range row {
			if tile == Pending {

				level.Map[y][x] = level.bfsFloor(Pos{x, y})
			}
		}
	}

	return level
}

func inRange(level *Level, pos Pos) bool {
	return pos.X < len(level.Map[0]) && pos.Y < len(level.Map) && pos.X >= 0 && pos.Y >= 0
}

func canWalk(level *Level, pos Pos) bool {
	if inRange(level, pos) {
		t := level.Map[pos.Y][pos.X]
		switch t {
		case StoneWall, ClosedDoor, Blank:
			return false
		default:
			return true
		}
	}
	return false
}

func checkDoor(level *Level, pos Pos) {
	t := level.Map[pos.Y][pos.X]

	if t == ClosedDoor {
		level.Map[pos.Y][pos.X] = OpenDoor
	}
}

func (player *Player) Move(to Pos, level *Level) {
	monster, exists := level.Monsters[to]
	if !exists {
		player.Pos = to
	} else {
		Attack(&player.Character, &monster.Character)
		level.AddEvent(player.Name + " attacked monster")

		if monster.HP <= 0 {
			level.AddEvent(player.Name + " killed the " + monster.Name)
			delete(level.Monsters, monster.Pos)
		}

		if level.Player.HP <= 0 {
			fmt.Println("You Died")
			panic("You Died")
		}

	}
}

func (gameStruct *Game) handleInput(input *Input) {
	level := gameStruct.Level
	p := level.Player
	switch input.Type {
	case Up:
		newPos := Pos{X: p.X, Y: p.Y - 1}
		if canWalk(level, newPos) {
			level.Player.Move(newPos, level)
		} else {
			checkDoor(level, newPos)
		}

	case Down:
		newPos := Pos{X: p.X, Y: p.Y + 1}
		if canWalk(level, newPos) {
			level.Player.Move(newPos, level)

		} else {
			checkDoor(level, newPos)
		}

	case Left:
		newPos := Pos{X: p.X - 1, Y: p.Y}
		if canWalk(level, newPos) {
			level.Player.Move(newPos, level)

		} else {
			checkDoor(level, newPos)
		}

	case Right:
		newPos := Pos{p.X + 1, p.Y}
		if canWalk(level, newPos) {
			level.Player.Move(newPos, level)
		} else {
			checkDoor(level, newPos)
		}

	case Search:
		level.astar(level.Player.Pos, Pos{X: 3, Y: 3})

	case CloseWindow:
		close(input.LevelChannel)
		chanIndex := 0
		for i, c := range gameStruct.LevelChans {
			if c == input.LevelChannel {
				chanIndex = i
				break
			}
		}
		gameStruct.LevelChans = append(gameStruct.LevelChans[:chanIndex], gameStruct.LevelChans[chanIndex+1:]...)
	}
}

func getNeighbors(level *Level, pos Pos) []Pos {
	neighbors := make([]Pos, 0, 4)

	left := Pos{X: pos.X - 1, Y: pos.Y}
	right := Pos{X: pos.X + 1, Y: pos.Y}
	up := Pos{X: pos.X, Y: pos.Y - 1}
	down := Pos{X: pos.X, Y: pos.Y + 1}

	if canWalk(level, right) {
		neighbors = append(neighbors, right)
	}
	if canWalk(level, left) {
		neighbors = append(neighbors, left)
	}
	if canWalk(level, up) {
		neighbors = append(neighbors, up)
	}
	if canWalk(level, down) {
		neighbors = append(neighbors, down)
	}

	return neighbors
}

func (level *Level) bfsFloor(start Pos) Tile {
	frontier := make([]Pos, 0, 8)

	frontier = append(frontier, start)

	visited := make(map[Pos]bool)
	visited[start] = true
	level.Debug = visited

	for len(frontier) > 0 {
		current := frontier[0]

		currentTile := level.Map[current.Y][current.X]
		switch currentTile {
		case DirtFloor:
			return DirtFloor
		default:
		}

		frontier = frontier[1:]
		for _, next := range getNeighbors(level, current) {
			if !visited[next] {
				frontier = append(frontier, next)
				visited[next] = true
			}
		}
	}
	return DirtFloor
}

func (level *Level) astar(start Pos, goal Pos) []Pos {
	frontier := make(pqueue, 0, 8)
	frontier = frontier.push(start, 1)

	cameFrom := make(map[Pos]Pos)
	cameFrom[start] = start

	costSoFar := make(map[Pos]int)
	costSoFar[start] = 0

	var current Pos
	for len(frontier) > 0 {
		frontier, current = frontier.pop()

		if current == goal {
			path := make([]Pos, 0)
			p := current
			for p != start {
				path = append(path, p)
				p = cameFrom[p]
			}
			path = append(path, p)

			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}

			return path
		}

		for _, next := range getNeighbors(level, current) {
			newCost := costSoFar[current] + 1
			_, exists := costSoFar[next]
			if !exists || newCost < costSoFar[next] {
				costSoFar[next] = newCost

				xDist := int(math.Abs(float64(goal.X - next.X)))
				yDist := int(math.Abs(float64(goal.Y - next.Y)))

				priority := newCost + xDist + yDist
				frontier = frontier.push(next, priority)

				cameFrom[next] = current
			}
		}
	}
	return nil
}

func (gameStruct *Game) Run() {

	for _, lchan := range gameStruct.LevelChans {
		lchan <- gameStruct.Level
	}

	for input := range gameStruct.InputChan {
		if input.Type == QuitGame {
			return
		}
		gameStruct.handleInput(input)
		for _, monster := range gameStruct.Level.Monsters {
			monster.Update(gameStruct.Level)
		}

		if len(gameStruct.LevelChans) == 0 {
			return
		}

		for _, lchan := range gameStruct.LevelChans {
			lchan <- gameStruct.Level
		}
	}
}

package game

import (
	"bufio"
	"math"
	"os"
	"strconv"
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

type Tile struct {
	Rune    rune
	Visible bool
}

const (
	StoneWall  rune = '#'
	DirtFloor  rune = '.'
	ClosedDoor rune = '|'
	OpenDoor   rune = '/'
	Blank      rune = 0
	Pending    rune = -1
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
	HP         int
	Strength   int
	Speed      float64
	AP         float64
	SightRange int
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

func (level *Level) Attack(c1, c2 *Character) {
	c1.AP--
	c1AttackPower := c1.Strength
	c2.HP -= c1AttackPower

	if c2.HP > 0 {
		level.AddEvent(c1.Name + " attacked " + c2.Name + " for " + strconv.Itoa(c1AttackPower))
	} else {
		level.AddEvent(c1.Name + " killed " + c2.Name)
	}
}

func (level *Level) AddEvent(event string) {
	level.Events[level.EventPos] = event
	level.EventPos++
	if level.EventPos == len(level.Events) {
		level.EventPos = 0
	}
}

func (level *Level) lineOfSight() {
	pos := level.Player.Pos
	dist := level.Player.SightRange

	for y := pos.Y - dist; y <= pos.Y+dist; y++ {
		for x := pos.X - dist; x <= pos.X+dist; x++ {
			xDelta := pos.X - x
			yDelta := pos.Y - y
			d := math.Sqrt(float64(xDelta*xDelta + yDelta*yDelta))
			if d <= float64(dist) {
				level.bresenham(pos, Pos{x, y})
			}
		}
	}
}

func (level *Level) bresenham(start, end Pos) {
	steep := math.Abs(float64(end.Y-start.Y)) > math.Abs(float64(end.X-start.X))

	if steep {
		start.X, start.Y = start.Y, start.X
		end.X, end.Y = end.Y, end.X
	}

	deltaY := math.Abs(float64(end.Y - start.Y))
	err := 0
	y := start.Y
	ystep := 1
	if start.Y >= end.Y {
		ystep = -1
	}

	if start.X > end.X {

		deltaX := start.X - end.X

		for x := start.X; x > end.X; x-- {
			var pos Pos
			if steep {
				pos = Pos{X: y, Y: x}
			} else {
				pos = Pos{X: x, Y: y}

			}
			level.Map[pos.Y][pos.X].Visible = true
			if !canSeeThrough(level, pos) {
				return
			}
			err += int(deltaY)
			if 2*err >= deltaX {
				y += ystep
				err -= deltaX
			}
		}
	} else {

		deltaX := end.X - start.X

		for x := start.X; x < end.X; x++ {
			var pos Pos
			if steep {
				pos = Pos{X: y, Y: x}
			} else {
				pos = Pos{X: x, Y: y}

			}
			level.Map[pos.Y][pos.X].Visible = true
			if !canSeeThrough(level, pos) {
				return
			}
			err += int(deltaY)
			if 2*err >= deltaX {
				y += ystep
				err -= deltaX
			}
		}
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

	// level.Debug = make(map[Pos]bool, 0)
	level.Events = make([]string, 10)
	level.Player.Name = "Dralanor"
	level.Player.Rune = '@'
	level.Player.HP = 20
	level.Player.Strength = 20
	level.Player.Speed = 1
	level.Player.AP = 0
	level.Player.SightRange = 7

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
				t.Rune = Blank
			case '#':
				t.Rune = StoneWall
			case '|':
				t.Rune = ClosedDoor
			case '/':
				t.Rune = OpenDoor
			case '.':
				t.Rune = DirtFloor
			case '@':
				level.Player.X = x
				level.Player.Y = y
				t.Rune = Pending
			case 'R':
				level.Monsters[Pos{x, y}] = NewRat(Pos{x, y})
				t.Rune = Pending
			case 'S':
				level.Monsters[Pos{x, y}] = NewSpider(Pos{x, y})
				t.Rune = Pending
			default:
				panic("Invalid character in map")
			}
			level.Map[y][x] = t
		}
	}

	for y, row := range level.Map {
		for x, tile := range row {
			if tile.Rune == Pending {

				level.Map[y][x] = level.bfsFloor(Pos{x, y})
			}
		}
	}
	level.lineOfSight()
	return level
}

func inRange(level *Level, pos Pos) bool {
	return pos.X < len(level.Map[0]) && pos.Y < len(level.Map) && pos.X >= 0 && pos.Y >= 0
}

func canWalk(level *Level, pos Pos) bool {
	if inRange(level, pos) {
		t := level.Map[pos.Y][pos.X]
		switch t.Rune {
		case StoneWall, ClosedDoor, Blank:
			return false
		}
		_, exists := level.Monsters[pos]
		return !exists
	}
	return false
}

func canSeeThrough(level *Level, pos Pos) bool {
	if inRange(level, pos) {
		t := level.Map[pos.Y][pos.X]
		switch t.Rune {
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

	if t.Rune == ClosedDoor {
		level.Map[pos.Y][pos.X].Rune = OpenDoor
		level.lineOfSight()
	}
}

func (player *Player) Move(to Pos, level *Level) {
	player.Pos = to
	for y, row := range level.Map {
		for x := range row {
			level.Map[y][x].Visible = false
		}
	}

	level.lineOfSight()
}

func (level *Level) resolveMovement(pos Pos) {
	monster, exists := level.Monsters[pos]
	if exists {
		level.Attack(&level.Player.Character, &monster.Character)
		if monster.HP <= 0 {
			delete(level.Monsters, monster.Pos)
		}
		if level.Player.HP <= 0 {
			panic("DEAD")
		}
	} else if canWalk(level, pos) {
		level.Player.Move(pos, level)
	} else {
		checkDoor(level, pos)
	}
}

func (gameStruct *Game) handleInput(input *Input) {
	level := gameStruct.Level
	p := level.Player
	switch input.Type {
	case Up:
		newPos := Pos{X: p.X, Y: p.Y - 1}
		level.resolveMovement(newPos)

	case Down:
		newPos := Pos{X: p.X, Y: p.Y + 1}
		level.resolveMovement(newPos)

	case Left:
		newPos := Pos{X: p.X - 1, Y: p.Y}
		level.resolveMovement(newPos)

	case Right:
		newPos := Pos{p.X + 1, p.Y}
		level.resolveMovement(newPos)

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

	for len(frontier) > 0 {
		current := frontier[0]

		currentTile := level.Map[current.Y][current.X]
		switch currentTile.Rune {
		case DirtFloor:
			return Tile{DirtFloor, false}
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
	return Tile{DirtFloor, false}
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

		// p := gameStruct.Level.Player.Pos
		// line := bresenham(p, Pos{X: p.X + 5, Y: p.Y + 5})
		// for _, pos := range line {
		// 	gameStruct.Level.Debug[pos] = true
		// }

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

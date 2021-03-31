package game

import (
	"bufio"
	"math"
	"os"
	"sort"
	"time"
)

type GameUI interface {
	Draw(*Level)
	GetInput() *Input
}

type InputType int

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	Quit
	Search
)

type Input struct {
	Type InputType
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
}

type Player struct {
	Entity
}

type Level struct {
	Map    [][]Tile
	Player Player
	Debug  map[Pos]bool
}

type priorityPos struct {
	Pos
	priority int
}

type priorityArray []priorityPos

func (p priorityArray) Len() int {
	return len(p)
}
func (p priorityArray) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p priorityArray) Less(i, j int) bool {
	return p[i].priority < p[j].priority
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
	level.Map = make([][]Tile, len(levelLines))

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
			case 'P':
				level.Player.X = x
				level.Player.Y = y
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

			SearchLoop:
				for searchX := x - 1; searchX <= x+1; searchX++ {
					for searchY := y - 1; searchY <= y+1; searchY++ {
						searchTile := level.Map[searchY][searchX]
						switch searchTile {
						case DirtFloor:
							level.Map[y][x] = DirtFloor
							break SearchLoop
						}
					}
				}
			}
		}
	}

	return level
}

func canWalk(level *Level, pos Pos) bool {
	t := level.Map[pos.Y][pos.X]
	switch t {
	case StoneWall, ClosedDoor, Blank:
		return false
	default:
		return true
	}
}

func checkDoor(level *Level, pos Pos) {
	t := level.Map[pos.Y][pos.X]

	if t == ClosedDoor {
		level.Map[pos.Y][pos.X] = OpenDoor
	}
}

func handleInput(ui GameUI, level *Level, input *Input) {
	p := level.Player
	switch input.Type {
	case Up:
		if canWalk(level, Pos{X: p.X, Y: p.Y - 1}) {
			level.Player.Y--
		} else {
			checkDoor(level, Pos{X: p.X, Y: p.Y - 1})
		}

	case Down:
		if canWalk(level, Pos{X: p.X, Y: p.Y + 1}) {
			level.Player.Y++
		} else {
			checkDoor(level, Pos{X: p.X, Y: p.Y + 1})
		}

	case Left:
		if canWalk(level, Pos{X: p.X - 1, Y: p.Y}) {
			level.Player.X--
		} else {
			checkDoor(level, Pos{X: p.X - 1, Y: p.Y})
		}

	case Right:
		if canWalk(level, Pos{p.X + 1, p.Y}) {
			level.Player.X++
		} else {
			checkDoor(level, Pos{p.X + 1, p.Y})
		}
	case Search:
		bfs(ui, level, level.Player.Pos)
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

func bfs(ui GameUI, level *Level, start Pos) {
	frontier := make([]Pos, 0, 8)

	frontier = append(frontier, start)

	visited := make(map[Pos]bool)
	visited[start] = true
	level.Debug = visited

	for len(frontier) > 0 {
		current := frontier[0]
		frontier = frontier[1:]
		for _, next := range getNeighbors(level, current) {
			if !visited[next] {
				frontier = append(frontier, next)
				visited[next] = true
				ui.Draw(level)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

}

func astar(ui GameUI, level *Level, start Pos, goal Pos) {
	frontier := make(priorityArray, 0, 8)
	frontier = append(frontier, priorityPos{start, 1})

	cameFrom := make(map[Pos]Pos)
	cameForm[start] = start

	costSoFar := make(map[Pos]int)
	costSoFar[start] = 0

	for len(frontier) > 0 {
		sort.Stable(frontier)
		current := frontier[0]

		if current.Pos == goal {
			break
		}

		frontier = frontier[1:]

		for _, next := range getNeighbors(level, current.Pos) {
			newCost := costSoFar[current.Pos] + 1
			_, exists := costSoFar[next]
			if !exists || newCost < costSoFar[next] {
			}
			costSoFar[next] = newCost
			xDist := int(math.Abs(float64(goal.X - next.X)))
			yDist := int(math.Abs(float64(goal.Y - next.Y)))
		}
	}
}

func Run(ui GameUI) {
	level := loadLevelFromFile("game/maps/level1.map")

	for {

		ui.Draw(level)

		input := ui.GetInput()

		if input != nil && input.Type == Quit {
			return
		}
		handleInput(ui, level, input)
	}
}

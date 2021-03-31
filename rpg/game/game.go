package game

import (
	"bufio"
	"os"
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

func canWalk(level *Level, x, y int) bool {
	t := level.Map[y][x]
	switch t {
	case StoneWall, ClosedDoor, Blank:
		return false
	default:
		return true
	}
}

func checkDoor(level *Level, x, y int) {
	t := level.Map[y][x]

	if t == ClosedDoor {
		level.Map[y][x] = OpenDoor
	}
}

func handleInput(level *Level, input *Input) {
	p := level.Player
	switch input.Type {
	case Up:
		if canWalk(level, p.X, p.Y-1) {
			level.Player.Y--
		} else {
			checkDoor(level, p.X, p.Y-1)
		}

	case Down:
		if canWalk(level, p.X, p.Y+1) {
			level.Player.Y++
		} else {
			checkDoor(level, p.X, p.Y+1)
		}

	case Left:
		if canWalk(level, p.X-1, p.Y) {
			level.Player.X--
		} else {
			checkDoor(level, p.X-1, p.Y)
		}

	case Right:
		if canWalk(level, p.X+1, p.Y) {
			level.Player.X++
		} else {
			checkDoor(level, p.X+1, p.Y)
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
		handleInput(level, input)
	}
}

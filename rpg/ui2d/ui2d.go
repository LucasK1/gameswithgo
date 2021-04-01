package ui2d

import (
	"bufio"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/LucasK1/gameswithgo/rpg/game"
	"github.com/veandco/go-sdl2/sdl"
)

type ui struct {
	winWidth          int
	winHeight         int
	renderer          *sdl.Renderer
	window            *sdl.Window
	textureAtlas      *sdl.Texture
	textureIndex      map[game.Tile][]sdl.Rect
	keyboardState     []uint8
	prevKeyboardState []uint8
	centerX           int
	centerY           int
	r                 *rand.Rand
	levelChan         chan *game.Level
	inputChan         chan *game.Input
}

func NewUI(inputChan chan *game.Input, levelChan chan *game.Level) *ui {

	ui := &ui{}
	ui.inputChan = inputChan
	ui.levelChan = levelChan
	ui.r = rand.New(rand.NewSource(1))
	ui.winHeight = 720
	ui.winWidth = 1280

	window, err := sdl.CreateWindow("RPG", 200, 200, int32(ui.winWidth), int32(ui.winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	ui.window = window

	ui.renderer, err = sdl.CreateRenderer(ui.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	ui.textureAtlas = ui.imgFileToTexture("ui2d/assets/tiles.png")
	ui.loadTextureIndex()

	ui.keyboardState = sdl.GetKeyboardState()
	ui.prevKeyboardState = make([]uint8, len(ui.keyboardState))
	copy(ui.prevKeyboardState, ui.keyboardState)

	ui.centerX = -1
	ui.centerY = -1

	return ui
}

func (ui *ui) loadTextureIndex() {
	ui.textureIndex = make(map[game.Tile][]sdl.Rect)

	infile, err := os.Open("ui2d/assets/atlas-index.txt")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		tileRune := game.Tile(line[0])
		xy := line[1:]

		splitXYC := strings.Split(xy, ",")

		x, err := strconv.ParseInt(strings.TrimSpace(splitXYC[0]), 10, 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(strings.TrimSpace(splitXYC[1]), 10, 64)
		if err != nil {
			panic(err)
		}
		variationCount, err := strconv.ParseInt(strings.TrimSpace(splitXYC[2]), 10, 64)
		if err != nil {
			panic(err)
		}

		var rects []sdl.Rect
		for i := 0; i < int(variationCount); i++ {
			rects = append(rects, sdl.Rect{X: int32(x * 32), Y: int32(y * 32), W: 32, H: 32})
			x++
			if x > 62 {
				x = 0
				y++
			}
		}
		ui.textureIndex[tileRune] = rects
	}
}

func (ui *ui) imgFileToTexture(filename string) *sdl.Texture {
	infile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer infile.Close()

	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	i := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {

			r, g, b, a := img.At(x, y).RGBA()
			pixels[i] = byte(r / 256)
			i++
			pixels[i] = byte(g / 256)
			i++
			pixels[i] = byte(b / 256)
			i++
			pixels[i] = byte(a / 256)
			i++
		}
	}

	tex, err := ui.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, int32(w), int32(h))
	if err != nil {
		panic(err)
	}
	tex.Update(nil, pixels, w*4)

	err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}

	return tex
}

func init() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

}

func (ui *ui) Draw(level *game.Level) {

	p := level.Player

	if ui.centerX == -1 && ui.centerY == -1 {
		ui.centerX = p.X
		ui.centerY = p.Y
	}

	limit := 5

	if p.X > ui.centerX+limit {
		ui.centerX++
	} else if p.X < ui.centerX-limit {
		ui.centerX--
	} else if p.Y > ui.centerY+limit {
		ui.centerY++
	} else if p.Y < ui.centerY-limit {
		ui.centerY--
	}

	offsetX := int32((ui.winWidth / 2) - ui.centerX*32)
	offsetY := int32((ui.winHeight / 2) - ui.centerY*32)

	ui.renderer.Clear()
	ui.r.Seed(1)

	for y, row := range level.Map {
		for x, tile := range row {
			if tile != game.Blank {
				srcRects := ui.textureIndex[tile]
				srcRect := srcRects[ui.r.Intn(len(srcRects))]
				dstRect := sdl.Rect{X: int32(x)*32 + offsetX, Y: int32(y)*32 + offsetY, W: 32, H: 32}

				pos := game.Pos{X: x, Y: y}
				if level.Debug[pos] {
					ui.textureAtlas.SetColorMod(128, 0, 0)
				} else {
					ui.textureAtlas.SetColorMod(255, 255, 255)
				}

				ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)
			}
		}
	}

	ui.renderer.Copy(ui.textureAtlas, &sdl.Rect{X: 21 * 32, Y: 59 * 32, W: 32, H: 32}, &sdl.Rect{X: int32(p.X)*32 + offsetX, Y: int32(p.Y)*32 + offsetY, W: 32, H: 32})

	ui.renderer.Present()

}

func (ui *ui) Run() {

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				ui.inputChan <- &game.Input{Type: game.QuitGame}
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					ui.inputChan <- &game.Input{Type: game.CloseWindow, LevelChannel: ui.levelChan}
				}
			}
		}

		select {
		case newLevel, ok := <-ui.levelChan:
			if ok {
				ui.Draw(newLevel)
			}
		default:
		}

		if sdl.GetKeyboardFocus() == ui.window || sdl.GetMouseFocus() == ui.window {
			var input game.Input
			if ui.keyboardState[sdl.SCANCODE_UP] == 0 && ui.prevKeyboardState[sdl.SCANCODE_UP] != 0 {
				input.Type = game.Up
			}
			if ui.keyboardState[sdl.SCANCODE_DOWN] == 0 && ui.prevKeyboardState[sdl.SCANCODE_DOWN] != 0 {
				input.Type = game.Down
			}
			if ui.keyboardState[sdl.SCANCODE_LEFT] == 0 && ui.prevKeyboardState[sdl.SCANCODE_LEFT] != 0 {
				input.Type = game.Left
			}
			if ui.keyboardState[sdl.SCANCODE_RIGHT] == 0 && ui.prevKeyboardState[sdl.SCANCODE_RIGHT] != 0 {
				input.Type = game.Right
			}
			if ui.keyboardState[sdl.SCANCODE_S] == 0 && ui.prevKeyboardState[sdl.SCANCODE_S] != 0 {
				input.Type = game.Search
			}

			copy(ui.prevKeyboardState, ui.keyboardState)

			if input.Type != game.None {
				ui.inputChan <- &input
			}
		}
		sdl.Delay(10)
	}
}

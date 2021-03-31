package ui2d

import (
	"bufio"
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/LucasK1/gameswithgo/rpg/game"
	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight = 1280, 720

var renderer *sdl.Renderer
var textureAtlas *sdl.Texture
var textureIndex map[game.Tile][]sdl.Rect
var keyboardState []uint8
var prevKeyboardState []uint8

var centerX int
var centerY int

func loadTextureIndex() {
	textureIndex = make(map[game.Tile][]sdl.Rect)

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
		textureIndex[tileRune] = rects
	}
}

func imgFileToTexture(filename string) *sdl.Texture {
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

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, int32(w), int32(h))
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
	sdl.LogSetAllPriority(sdl.LOG_PRIORITY_VERBOSE)
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}

	window, err := sdl.CreateWindow("RPG", 200, 200, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)

	if err != nil {
		fmt.Println(err)
		return
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	textureAtlas = imgFileToTexture("ui2d/assets/tiles.png")
	loadTextureIndex()

	keyboardState = sdl.GetKeyboardState()
	prevKeyboardState = make([]uint8, len(keyboardState))
	copy(prevKeyboardState, keyboardState)

	centerX = -1
	centerY = -1
}

type UI2d struct {
}

func (ui *UI2d) Draw(level *game.Level) {

	p := level.Player

	if centerX == -1 && centerY == -1 {
		centerX = p.X
		centerY = p.Y
	}

	limit := 5

	if p.X > centerX+limit {
		centerX++
	} else if p.X < centerX-limit {
		centerX--
	} else if p.Y > centerY+limit {
		centerY++
	} else if p.Y < centerY-limit {
		centerY--
	}

	offsetX := int32((winWidth / 2) - centerX*32)
	offsetY := int32((winHeight / 2) - centerY*32)

	rand.Seed(1)

	renderer.Clear()
	for y, row := range level.Map {
		for x, tile := range row {
			if tile != game.Blank {
				srcRects := textureIndex[tile]
				srcRect := srcRects[rand.Intn(len(srcRects))]
				dstRect := sdl.Rect{X: int32(x)*32 + offsetX, Y: int32(y)*32 + offsetY, W: 32, H: 32}

				pos := game.Pos{X: x, Y: y}
				if level.Debug[pos] {
					textureAtlas.SetColorMod(128, 0, 0)
				} else {
					textureAtlas.SetColorMod(255, 255, 255)
				}

				renderer.Copy(textureAtlas, &srcRect, &dstRect)
			}
		}
	}

	renderer.Copy(textureAtlas, &sdl.Rect{X: 21 * 32, Y: 59 * 32, W: 32, H: 32}, &sdl.Rect{X: int32(p.X)*32 + offsetX, Y: int32(p.Y)*32 + offsetY, W: 32, H: 32})

	renderer.Present()

}

func (ui *UI2d) GetInput() *game.Input {

	for {

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return &game.Input{Type: game.Quit}
			}
		}

		var input game.Input
		if keyboardState[sdl.SCANCODE_UP] == 0 && prevKeyboardState[sdl.SCANCODE_UP] != 0 {
			input.Type = game.Up
		}
		if keyboardState[sdl.SCANCODE_DOWN] == 0 && prevKeyboardState[sdl.SCANCODE_DOWN] != 0 {
			input.Type = game.Down
		}
		if keyboardState[sdl.SCANCODE_LEFT] == 0 && prevKeyboardState[sdl.SCANCODE_LEFT] != 0 {
			input.Type = game.Left
		}
		if keyboardState[sdl.SCANCODE_RIGHT] == 0 && prevKeyboardState[sdl.SCANCODE_RIGHT] != 0 {
			input.Type = game.Right
		}

		copy(prevKeyboardState, keyboardState)

		if input.Type != game.None {
			return &input
		}

	}
}

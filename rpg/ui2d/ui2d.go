package ui2d

import (
	"bufio"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/LucasK1/gameswithgo/rpg/game"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type sounds struct {
	doorOpens []*mix.Chunk
	footsteps []*mix.Chunk
}

func playRandomSound(chunks []*mix.Chunk, volume int) {
	chunkIndex := rand.Intn(len(chunks))
	chunks[chunkIndex].Volume(volume)
	chunks[chunkIndex].Play(-1, 0)
}

type ui struct {
	winWidth          int
	winHeight         int
	renderer          *sdl.Renderer
	window            *sdl.Window
	textureAtlas      *sdl.Texture
	textureIndex      map[rune][]sdl.Rect
	keyboardState     []uint8
	prevKeyboardState []uint8
	centerX           int
	centerY           int
	r                 *rand.Rand
	levelChan         chan *game.Level
	inputChan         chan *game.Input
	fontSmall         *ttf.Font
	fontMedium        *ttf.Font
	fontLarge         *ttf.Font
	strToTexSm        map[string]*sdl.Texture
	strToTexMd        map[string]*sdl.Texture
	strToTexLg        map[string]*sdl.Texture
	eventBackground   *sdl.Texture
	sounds            sounds
}

func NewUI(inputChan chan *game.Input, levelChan chan *game.Level) *ui {

	ui := &ui{}
	ui.strToTexSm = make(map[string]*sdl.Texture)
	ui.strToTexMd = make(map[string]*sdl.Texture)
	ui.strToTexLg = make(map[string]*sdl.Texture)
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

	ui.fontSmall, err = ttf.OpenFont("ui2d/assets/font.ttf", int(float64(ui.winHeight)*0.025))
	if err != nil {
		panic(err)
	}
	ui.fontMedium, err = ttf.OpenFont("ui2d/assets/font.ttf", 32)
	if err != nil {
		panic(err)
	}
	ui.fontLarge, err = ttf.OpenFont("ui2d/assets/font.ttf", 64)
	if err != nil {
		panic(err)
	}

	ui.eventBackground = ui.GetSinglePixelTex(sdl.Color{R: 0, G: 0, B: 0, A: 156})
	ui.eventBackground.SetBlendMode(sdl.BLENDMODE_BLEND)

	err = mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096)
	if err != nil {
		panic(err)
	}

	music, err := mix.LoadMUS("ui2d/assets/ambient.ogg")
	if err != nil {
		panic(err)
	}

	err = music.Play(-1)
	if err != nil {
		panic(err)
	}

	footstepBase := "ui2d/assets/footstep0"
	for i := 0; i < 10; i++ {
		footstepFile := footstepBase + strconv.Itoa(i) + ".ogg"
		footstep, err := mix.LoadWAV(footstepFile)
		if err != nil {
			panic(err)
		}
		ui.sounds.footsteps = append(ui.sounds.footsteps, footstep)
	}

	doorOpenBase := "ui2d/assets/doorOpen_"
	for i := 1; i < 3; i++ {
		doorOpenFile := doorOpenBase + strconv.Itoa(i) + ".ogg"
		doorOpen, err := mix.LoadWAV(doorOpenFile)
		if err != nil {
			panic(err)
		}
		ui.sounds.doorOpens = append(ui.sounds.doorOpens, doorOpen)
	}

	return ui
}

type FontSize int

const (
	FontSmall FontSize = iota
	FontMedium
	FontLarge
)

func (ui *ui) stringToTexture(s string, color sdl.Color, size FontSize) *sdl.Texture {

	var font *ttf.Font
	switch size {
	case FontSmall:
		font = ui.fontSmall
		tex, exists := ui.strToTexSm[s]
		if exists {
			return tex
		}
	case FontMedium:
		font = ui.fontMedium
		tex, exists := ui.strToTexMd[s]
		if exists {
			return tex
		}
	case FontLarge:
		font = ui.fontLarge
		tex, exists := ui.strToTexLg[s]
		if exists {
			return tex
		}
	}

	fontSurface, err := font.RenderUTF8Blended(s, color)
	if err != nil {
		panic(err)
	}

	tex, err := ui.renderer.CreateTextureFromSurface(fontSurface)
	if err != nil {
		panic(err)
	}

	switch size {
	case FontSmall:
		ui.strToTexSm[s] = tex

	case FontMedium:
		ui.strToTexMd[s] = tex

	case FontLarge:
		ui.strToTexLg[s] = tex
	}

	return tex
}

func (ui *ui) loadTextureIndex() {
	ui.textureIndex = make(map[rune][]sdl.Rect)

	infile, err := os.Open("ui2d/assets/atlas-index.txt")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		tileRune := rune(line[0])
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
	err = ttf.Init()
	if err != nil {
		panic(err)
	}

	err = mix.Init(mix.INIT_OGG)
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
		diff := level.Player.X - (ui.centerX + limit)
		ui.centerX += diff
	} else if p.X < ui.centerX-limit {
		diff := (ui.centerX - limit) - level.Player.X
		ui.centerX -= diff
	} else if p.Y > ui.centerY+limit {
		diff := level.Player.Y - (ui.centerY + limit)
		ui.centerY += diff
	} else if p.Y < ui.centerY-limit {
		diff := (ui.centerY - limit) - level.Player.Y
		ui.centerY -= diff
	}

	offsetX := int32((ui.winWidth / 2) - ui.centerX*32)
	offsetY := int32((ui.winHeight / 2) - ui.centerY*32)

	ui.renderer.Clear()
	ui.r.Seed(1)

	for y, row := range level.Map {
		for x, tile := range row {
			if tile.Rune != game.Blank {
				srcRects := ui.textureIndex[tile.Rune]
				srcRect := srcRects[ui.r.Intn(len(srcRects))]
				if tile.Visible || tile.Seen {
					dstRect := sdl.Rect{X: int32(x)*32 + offsetX, Y: int32(y)*32 + offsetY, W: 32, H: 32}

					pos := game.Pos{X: x, Y: y}
					if level.Debug[pos] {
						ui.textureAtlas.SetColorMod(128, 0, 0)
					} else if tile.Seen && !tile.Visible {
						ui.textureAtlas.SetColorMod(128, 128, 128)
					} else {
						ui.textureAtlas.SetColorMod(255, 255, 255)
					}

					ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)

					if tile.OverlayRune != game.Blank {
						srcRect := ui.textureIndex[tile.OverlayRune][0]
						ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)

					}
				}
			}
		}
	}
	ui.textureAtlas.SetColorMod(255, 255, 255)

	for pos, monster := range level.Monsters {

		if level.Map[pos.Y][pos.X].Visible {
			monsterSrcRect := ui.textureIndex[monster.Rune][0]

			ui.renderer.Copy(ui.textureAtlas, &monsterSrcRect, &sdl.Rect{X: int32(pos.X)*32 + offsetX, Y: int32(pos.Y)*32 + offsetY, W: 32, H: 32})
		}
	}

	playerSrcRect := ui.textureIndex['@'][0]

	ui.renderer.Copy(ui.textureAtlas, &playerSrcRect, &sdl.Rect{X: int32(p.X)*32 + offsetX, Y: int32(p.Y)*32 + offsetY, W: 32, H: 32})

	textStart := int32(float64(ui.winHeight) * 0.69)
	textWidth := int32(float64(ui.winWidth) * 0.25)

	ui.renderer.Copy(ui.eventBackground, nil, &sdl.Rect{X: 0, Y: textStart, W: textWidth, H: int32(ui.winHeight) - textStart})

	i := level.EventPos
	count := 0

	_, fontSizeY, _ := ui.fontSmall.SizeUTF8("A")

	for {
		event := level.Events[i]
		if event != "" {
			tex := ui.stringToTexture(event, sdl.Color{R: 255, G: 0, B: 0, A: 0}, FontSmall)
			_, _, w, h, err := tex.Query()
			if err != nil {
				panic(err)
			}
			ui.renderer.Copy(tex, nil, &sdl.Rect{X: 5, Y: int32(count*fontSizeY) + textStart, W: w, H: h})
		}
		i = (i + 1) % len(level.Events)
		if i == level.EventPos {
			break
		}
		count++
	}

	ui.renderer.Present()
}

func (ui *ui) keyDownOnce(key uint8) bool {
	return ui.keyboardState[key] == 1 && ui.prevKeyboardState[key] == 0
}

// func (ui *ui) keyPressed(key uint8) bool {
// 	return ui.keyboardState[key] == 0 && ui.prevKeyboardState[key] == 1
// }

func (ui *ui) GetSinglePixelTex(color sdl.Color) *sdl.Texture {
	tex, err := ui.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, 1, 1)
	if err != nil {
		panic(err)
	}
	pixels := make([]byte, 4)
	pixels[0] = color.R
	pixels[1] = color.G
	pixels[2] = color.B
	pixels[3] = color.A

	tex.Update(nil, pixels, 4)
	return tex
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
				switch newLevel.LastEvent {
				case game.Move:
					playRandomSound(ui.sounds.footsteps, 10)
				case game.DoorOpen:
					playRandomSound(ui.sounds.doorOpens, 32)
				default:

				}
				ui.Draw(newLevel)
			}
		default:
		}

		if sdl.GetKeyboardFocus() == ui.window || sdl.GetMouseFocus() == ui.window {
			var input game.Input
			if ui.keyDownOnce(sdl.SCANCODE_UP) {
				input.Type = game.Up
			}
			if ui.keyDownOnce(sdl.SCANCODE_DOWN) {
				input.Type = game.Down
			}
			if ui.keyDownOnce(sdl.SCANCODE_LEFT) {
				input.Type = game.Left
			}
			if ui.keyDownOnce(sdl.SCANCODE_RIGHT) {
				input.Type = game.Right
			}

			copy(ui.prevKeyboardState, ui.keyboardState)

			if input.Type != game.None {
				ui.inputChan <- &input
			}
		}
		sdl.Delay(10)
	}
}

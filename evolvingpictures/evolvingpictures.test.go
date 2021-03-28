package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/LucasK1/gameswithgo/apt"
	"github.com/LucasK1/gameswithgo/evolvingpictures/gui"
	"github.com/veandco/go-sdl2/sdl"
)

var winWidth, winHeight int = 800, 600

var rows, cols, numPics int = 3, 3, rows * cols

type pixelResult struct {
	pixels []byte
	index  int
}

// type audioState struct {
// 	explosionBytes []byte
// 	deviceID       sdl.AudioDeviceID
// 	audioSpec      *sdl.AudioSpec
// }

type picture struct {
	r apt.Node
	g apt.Node
	b apt.Node
}

func (p *picture) String() string {
	return "R" + p.r.String() + "\n" + "G" + p.g.String() + "\n" + "B" + p.b.String()
}

func newPicture() *picture {
	p := &picture{}

	p.r = apt.GetRandomNode()
	p.g = apt.GetRandomNode()
	p.b = apt.GetRandomNode()

	num := rand.Intn(20) + 5
	for i := 0; i < num; i++ {
		p.r.AddRandom(apt.GetRandomNode())
	}
	num = rand.Intn(20) + 5
	for i := 0; i < num; i++ {
		p.g.AddRandom(apt.GetRandomNode())
	}
	num = rand.Intn(20) + 5
	for i := 0; i < num; i++ {
		p.b.AddRandom(apt.GetRandomNode())
	}

	for p.r.AddLeaf(apt.GetRandomLeaf()) {
	}
	for p.g.AddLeaf(apt.GetRandomLeaf()) {
	}
	for p.b.AddLeaf(apt.GetRandomLeaf()) {
	}

	return p
}

func (p *picture) pickRandomColor() apt.Node {
	r := rand.Intn(3)
	switch r {
	case 0:
		return p.r
	case 1:
		return p.g
	case 2:
		return p.b
	default:
		panic("pickRandomColor failed")
	}

}

func cross(a *picture, b *picture) *picture {
	aCopy := &picture{apt.CopyTree(a.r, nil), apt.CopyTree(a.g, nil), apt.CopyTree(a.b, nil)}
	aColor := aCopy.pickRandomColor()
	bColor := b.pickRandomColor()

	aIndex := rand.Intn(aColor.NodeCount())
	aNode, _ := apt.GetNthNode(aColor, aIndex, 0)

	bIndex := rand.Intn(bColor.NodeCount())
	bNode, _ := apt.GetNthNode(bColor, bIndex, 0)
	bNodeCopy := apt.CopyTree(bNode, bNode.GetParent())

	apt.ReplaceNode(aNode, bNodeCopy)
	return aCopy
}

func evolve(survivors []*picture) []*picture {
	newPics := make([]*picture, numPics)
	i := 0
	for i < len(survivors) {
		a := survivors[i]
		b := survivors[rand.Intn(len(survivors))]
		newPics[i] = cross(a, b)
		i++
	}
	for i < len(newPics) {
		a := survivors[rand.Intn(len(survivors))]
		b := survivors[rand.Intn(len(survivors))]
		newPics[i] = cross(a, b)
		i++
	}

	for _, pic := range newPics {
		r := rand.Intn(4)
		for i := 0; i < r; i++ {
			pic.mutate()
		}
	}

	return newPics
}

func (p *picture) mutate() {
	r := rand.Intn(3)
	var nodeToMutate apt.Node
	switch r {
	case 0:
		nodeToMutate = p.r
	case 1:
		nodeToMutate = p.g
	case 2:
		nodeToMutate = p.b
	}

	count := nodeToMutate.NodeCount()

	r = rand.Intn(count)

	nodeToMutate, _ = apt.GetNthNode(nodeToMutate, r, 0)

	mutation := apt.Mutate(nodeToMutate)
	if nodeToMutate == p.r {
		p.r = mutation
	} else if nodeToMutate == p.g {
		p.g = mutation
	} else if nodeToMutate == p.b {
		p.b = mutation
	}
}

func pixelsToTexture(renderer *sdl.Renderer, pixels []byte, w, h int) *sdl.Texture {
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(w), int32(h))
	if err != nil {
		panic(err)
	}
	tex.Update(nil, pixels, w*4)

	return tex
}

func aptToPixels(p *picture, w, h int) []byte {
	scale := float32(255 / 2)
	offset := float32(-1 * scale)
	pixels := make([]byte, w*h*4)
	pixelIndex := 0

	for yi := 0; yi < h; yi++ {
		y := float32(yi)/float32(h)*2 - 1
		for xi := 0; xi < w; xi++ {
			x := float32(xi)/float32(w)*2 - 1

			r := p.r.Eval(x, y)
			g := p.g.Eval(x, y)
			b := p.b.Eval(x, y)

			pixels[pixelIndex] = byte(r*scale - offset)
			pixelIndex++
			pixels[pixelIndex] = byte(g*scale - offset)
			pixelIndex++
			pixels[pixelIndex] = byte(b*scale - offset)
			pixelIndex++
			pixelIndex++
		}
	}
	return pixels
}

func main() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Evolving Pictures", 200, 200, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	// explosionBytes, audioSpec := sdl.LoadWAV("explode.wav")
	// audioID, err := sdl.OpenAudioDevice("", false, audioSpec, nil, 0)
	// if err != nil {
	// 	panic(err)
	// }
	// defer sdl.FreeWAV(explosionBytes)

	// audioState := audioState{explosionBytes, audioID, audioSpec}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	var elapsedTime float32

	rand.Seed(time.Now().UTC().UnixNano())

	picTrees := make([]*picture, numPics)
	for i := range picTrees {
		picTrees[i] = newPicture()
	}

	picWidth := int(float32(winWidth/cols) * float32(0.9))
	picHeight := int(float32(winHeight/rows) * float32(0.8))

	pixelsChannel := make(chan pixelResult, numPics)
	buttons := make([]*gui.ImageButton, numPics)

	eveolveButtonTex := gui.GetSinglePixel(renderer, sdl.Color{R: 255, G: 255, B: 255, A: 0})
	evolveRect := sdl.Rect{X: int32(float32(winWidth)/2 - float32(picWidth)/2), Y: int32(float32(winHeight) - (float32(winHeight) * 0.10)), W: int32(picWidth), H: int32(float32(winHeight) * 0.08)}

	evolveButton := gui.NewImageButton(renderer, eveolveButtonTex, evolveRect, sdl.Color{R: 255, G: 255, B: 255, A: 0})

	for i := range picTrees {
		go func(i int) {
			pixels := aptToPixels(picTrees[i], picWidth, picHeight)
			pixelsChannel <- pixelResult{pixels, i}
		}(i)
	}

	keyboardState := sdl.GetKeyboardState()
	mouseState := gui.GetMouseState()

	for {
		frameStart := time.Now()

		mouseState.Update()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		if keyboardState[sdl.SCANCODE_ESCAPE] != 0 {
			return
		}

		select {
		case pixelsAndIndex, ok := <-pixelsChannel:
			if ok {
				tex := pixelsToTexture(renderer, pixelsAndIndex.pixels, picWidth, picHeight)
				xi := pixelsAndIndex.index % cols
				yi := (pixelsAndIndex.index - xi) / cols

				x := int32(xi * picWidth)
				y := int32(yi * picHeight)
				xPad := int32(float32(winWidth) * 0.1 / float32(cols+1))
				yPad := int32(float32(winHeight) * 0.1 / float32(rows+1))

				x += xPad * (int32(xi) + 1)
				y += yPad * (int32(yi) + 1)

				rect := sdl.Rect{X: x, Y: y, W: int32(picWidth), H: int32(picHeight)}
				button := gui.NewImageButton(renderer, tex, rect, sdl.Color{R: 255, G: 255, B: 255, A: 0})
				buttons[pixelsAndIndex.index] = button

			}
		default:

		}

		renderer.Clear()
		for _, button := range buttons {
			if button != nil {
				button.Update(mouseState)
				if button.WasLeftClicked {
					button.IsSelected = !button.IsSelected
				}
				button.Draw(renderer)
			}
		}

		evolveButton.Update(mouseState)
		if evolveButton.WasLeftClicked {
			selectedPictures := make([]*picture, 0)
			for i, button := range buttons {
				if button.IsSelected {
					selectedPictures = append(selectedPictures, picTrees[i])
				}
			}
			if len(selectedPictures) != 0 {
				for i := range buttons {
					buttons[i] = nil
				}
				picTrees = evolve(selectedPictures)
				for i, picTree := range picTrees {
					go func(i int, picTree *picture) {
						pixels := aptToPixels(picTree, picWidth, picHeight)
						pixelsChannel <- pixelResult{pixels, i}
					}(i, picTree)
				}
			}
		}
		evolveButton.Draw(renderer)

		renderer.Present()
		elapsedTime = float32(time.Since(frameStart).Milliseconds())

		if elapsedTime < 5 {
			sdl.Delay(5 - uint32(elapsedTime))
			// elapsedTime = float32(time.Since(frameStart).Milliseconds())
		}
	}
}

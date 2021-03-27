package main

import (
	"fmt"
	"time"

	. "github.com/LucasK1/gameswithgo/apt"
	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight, winDepth int = 800, 600, 100

type audioState struct {
	explosionBytes []byte
	deviceID       sdl.AudioDeviceID
	audioSpec      *sdl.AudioSpec
}

type mouseState struct {
	leftButton  bool
	rightButton bool
	x, y        int
}

func getMouseState() mouseState {
	mouseX, mouseY, mouseButtonState := sdl.GetMouseState()

	leftButton := mouseButtonState & sdl.ButtonLMask()
	rightButton := mouseButtonState & sdl.ButtonRMask()

	var result mouseState
	result.x = int(mouseX)
	result.y = int(mouseY)
	result.leftButton = !(leftButton == 0)
	result.rightButton = !(rightButton == 0)

	return result
}

type rgba struct {
	r, g, b byte
}

func pixelsToTexture(renderer *sdl.Renderer, pixels []byte, w, h int) *sdl.Texture {
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(w), int32(h))
	if err != nil {
		panic(err)
	}
	tex.Update(nil, pixels, w*4)

	return tex
}

func aptToTexture(redNode, greenNode, blueNode Node, w, h int, renderer *sdl.Renderer) *sdl.Texture {
	scale := float32(255 / 2)
	offset := float32(-1 * scale)
	pixels := make([]byte, w*h*4)
	pixelIndex := 0

	for yi := 0; yi < h; yi++ {
		y := float32(yi)/float32(h)*2 - 1
		for xi := 0; xi < w; xi++ {
			x := float32(xi)/float32(w)*2 - 1

			r := redNode.Eval(x, y)
			g := greenNode.Eval(x, y)
			b := blueNode.Eval(x, y)

			pixels[pixelIndex] = byte(r*scale - offset)
			pixelIndex++
			pixels[pixelIndex] = byte(g*scale - offset)
			pixelIndex++
			pixels[pixelIndex] = byte(b*scale - offset)
			pixelIndex++
			pixelIndex++
		}
	}
	return pixelsToTexture(renderer, pixels, w, h)
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

	// currentMouseState := getMouseState()
	// prevMouseState := currentMouseState

	x := &OpX{}
	y := &OpY{}
	sine := &OpSin{}
	noise := &OpNoise{}
	mult := &OpMult{}
	plus := &OpPlus{}

	mult.LeftChild = x
	mult.RightChild = noise

	noise.LeftChild = x
	noise.RightChild = y

	sine.Child = mult
	plus.LeftChild = y
	plus.RightChild = sine

	// fmt.Println(noise.String())
	texture := aptToTexture(plus, mult, sine, 800, 600, renderer)

	for {
		frameStart := time.Now()

		currentMouseState := getMouseState()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.TouchFingerEvent:
				if e.Type == sdl.FINGERDOWN {
					touchX := e.X * float32(winWidth)
					touchY := e.Y * float32(winHeight)
					currentMouseState.x = int(touchX)
					currentMouseState.y = int(touchY)
					currentMouseState.leftButton = true
				}
			}
		}

		renderer.Copy(texture, nil, nil)

		renderer.Present()
		elapsedTime = float32(time.Since(frameStart).Milliseconds())
		// fmt.Println(elapsedTime, "ms/frame")

		if elapsedTime < 5 {
			sdl.Delay(5 - uint32(elapsedTime))
			elapsedTime = float32(time.Since(frameStart).Milliseconds())
		}

		// prevMouseState = currentMouseState
	}
}

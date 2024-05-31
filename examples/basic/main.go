package main

import (
	"image/color"
	"log"
	rand "math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	prestige "github.com/mbrc12/prestige"
)

const (
	WIDTH  = 640
	HEIGHT = 360
)

// Just a nice color instead of a stark white
var (
	ALABASTER = color.RGBA{237, 234, 224, 255}
)

type State struct {
	Manager *prestige.SceneManager[State]
}

type BasicScene struct {
	x     float64
	color color.Color
}

func (self *BasicScene) Update(state *State) error {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		self.x -= 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		self.x += 10
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		log.Println("Transitioning ...")

		if err := state.Manager.Transition(&BasicScene{}, prestige.NewFadeTransition(1)); err != nil {
            log.Fatal(err)
		}
	}

	return nil
}

func (self *BasicScene) Draw(state *State, screen *ebiten.Image) {
	screen.Fill(self.color)
	vector.DrawFilledCircle(screen, float32(self.x), 180, 30, ALABASTER, true)
}

func randomColor() (col color.Color) {
	f := func() uint8 { return uint8(rand.Int() % 255) }
	col = color.RGBA{f(), f(), f(), 255}
	return
}

func (self *BasicScene) EnterStart(state *State) {
	self.x = 0
	self.color = randomColor()
}

func (self *BasicScene) EnterEnd(state *State) {
}

func (self *BasicScene) ExitStart(state *State) {
}

func (self *BasicScene) ExitEnd(state *State) {
}

func (scene *BasicScene) Layout(state *State, w, h int) (int, int) {
	return WIDTH, HEIGHT
}

func main() {
	basicScene := &BasicScene{}
	state := State{}
	manager := prestige.NewSceneManager(basicScene, &state)
	state.Manager = manager

	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Scenery Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(manager); err != nil {
		log.Fatal(err)
	}
}

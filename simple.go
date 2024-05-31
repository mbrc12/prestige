package prestige

import "github.com/hajimehoshi/ebiten/v2"

// A simple transition that immediately transfers from the current state to the next
type SimpleTransition struct {
}

func (self *SimpleTransition) Start() {
}

func (self *SimpleTransition) Update() bool {
	return true
}

func (self *SimpleTransition) Interpolate(screen *ebiten.Image, src *ebiten.Image, dest *ebiten.Image) {
}

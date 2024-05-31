package prestige

import "github.com/hajimehoshi/ebiten/v2"

type SimpleTransition struct {
}

func (self *SimpleTransition) Start() {
}

func (self *SimpleTransition) Update() bool {
	return true
}

func (self *SimpleTransition) Interpolate(screen *ebiten.Image, src *ebiten.Image, dest *ebiten.Image) {
}

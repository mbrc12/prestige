package prestige

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type FadeTransition struct {
	startTime time.Time
	duration  float64
}

func NewFadeTransition(duration float64) *FadeTransition {
	return &FadeTransition{duration: duration}
}

func (self *FadeTransition) Start() {
	self.startTime = time.Now()
}

func (self *FadeTransition) Update() bool {
	if time.Now().Sub(self.startTime).Seconds() > self.duration {
		return true
	}

	return false
}

func (self *FadeTransition) Interpolate(screen *ebiten.Image, src *ebiten.Image, dest *ebiten.Image) {
	t := time.Now().Sub(self.startTime).Seconds() / self.duration

	drawOptions := &ebiten.DrawImageOptions{}

	if t < 0.5 {
		drawOptions.ColorScale.ScaleAlpha(float32(1 - t/0.5))
		screen.DrawImage(src, drawOptions)
	} else {
		drawOptions.ColorScale.ScaleAlpha(float32(t/0.5 - 1))
		screen.DrawImage(dest, drawOptions)
	}
}

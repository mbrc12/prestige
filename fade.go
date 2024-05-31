package prestige

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// A simple fade transition, that uses half the time to fade out of the current scene
// and half the time to fade into the new scene. You may wish to freeze inputs during the fade,
// which you can accomplish in the `ExitStart` and `EnterStart` phase of your scenes,
// and unfreeze during `EnterEnd` for the new scene.
type FadeTransition struct {
	startTime time.Time
	duration  float64
}

// Pass in the total duration to construct a new fade transition.
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

package prestige

import (
	"errors"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Your should implement this on *YourScene, and T should be your global state across scenes.
// Your scene may have scene-local state as well. 
type Scene[T any] interface {
	// Transition has started to enter your scene, load your stuff now
	EnterStart(state *T)
	// Transition has finished entering your scene, enable scene mechanics if they were disabled
	EnterEnd(state *T)

	Draw(state *T, screen *ebiten.Image)
	Update(state *T) error
	Layout(state *T, outsideWidth, outsideHeight int) (screenWidth, screenHeight int)

	// Transition has started to exit your scene, disable scene mechanics if you like
	ExitStart(state *T)
	// Transition has finished exiting your scene, wrap up now
	ExitEnd(state *T)
}

type Transition interface {
	// Called before the transition begins. Initialize the state
	Start()

	// Interpolate between two images produced by the two scenes between whom this transition exists
	Interpolate(screen *ebiten.Image, src *ebiten.Image, dest *ebiten.Image)

	// Update the state of the transition and return whether the transition has finished
	Update() bool
}

// The only field accessible here is State. The only use for this should be to implement
// custom transitions that also depend on the current state (for instance if you wish to show
// scores on the transition screen).
type SceneManager[T any] struct {
	State      *T
	current    Scene[T]
	next       Scene[T]
	transition Transition

	currentImage *ebiten.Image
	nextImage    *ebiten.Image
}


// Construct a new scene manager, given the initial scene and initial state.
func NewSceneManager[T any](firstScene Scene[T], state *T) (manager *SceneManager[T]) {
	manager = &SceneManager[T]{}
	manager.State = state
	manager.current = firstScene
	firstScene.EnterStart(state)
	firstScene.EnterEnd(state)

	return
}

// Call this to transition to a new scene. 
func (self *SceneManager[T]) Transition(dest Scene[T], transition Transition) error {
	if self.transition != nil {
		return errors.New("New transition cannot start during another transition")
	}

	self.transition = transition
	self.transition.Start()

	self.currentImage = nil
	self.nextImage = nil

	self.current.ExitStart(self.State)

	self.next = dest
	self.next.EnterStart(self.State)

	return nil
}

func correctImageSizes(image **ebiten.Image, bounds image.Rectangle) {
	if *image == nil || !(**image).Bounds().Eq(bounds) {
		*image = ebiten.NewImage(bounds.Dx(), bounds.Dy())
	}
}

func (self *SceneManager[T]) Draw(screen *ebiten.Image) {
	if self.transition == nil {
		self.current.Draw(self.State, screen)
		return
	}

	// If the current image sizes for the target images are not what the screen sizes are,
	// or do not exist yet, make new images
	bounds := screen.Bounds()
	correctImageSizes(&self.currentImage, bounds)
	correctImageSizes(&self.nextImage, bounds)

	// Draw the images from the source and destination to their respective image targets
	self.current.Draw(self.State, self.currentImage)
	self.next.Draw(self.State, self.nextImage)

	// Ask the transition to interpolate the images and produce a new image
	self.transition.Interpolate(screen, self.currentImage, self.nextImage)
}

func (self *SceneManager[T]) Update() error {
	if self.transition == nil {
		return self.current.Update(self.State)
	}

	if err := self.current.Update(self.State); err != nil {
		return err
	}

	if err := self.next.Update(self.State); err != nil {
		return err
	}

	finished := self.transition.Update()

	// Transition has finished
	if finished {
		// call the necessary functions to inform the scenes that the transition has ended
		self.current.ExitEnd(self.State)
		self.next.EnterEnd(self.State)

		// deallocate images and reset the transition
		self.transition = nil
		self.currentImage.Deallocate()
		self.nextImage.Deallocate()

		// change current state to next
		self.current = self.next
		self.next = nil
	}

	return nil
}

func (self *SceneManager[T]) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return self.current.Layout(self.State, outsideWidth, outsideHeight)
}

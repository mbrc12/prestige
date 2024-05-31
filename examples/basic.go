package examples

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
    scenery "github.com/mbrc12/scenery"
)


type State struct {
    Manager *scenery.SceneManager[State]
}



func main() {
	basicScene := &BasicScene{}
    state := State{}
    manager := scene.NewSceneManager(basicScene, &state)
    state.Manager = manager

	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Elemental")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(manager); err != nil {
		log.Fatal(err)
	}
}

# Prestige

[![Go Reference](https://pkg.go.dev/badge/github.com/mbrc12/prestige.svg)](https://pkg.go.dev/github.com/mbrc12/prestige)

A simple scene manager for [ebitengine](https://ebitengine.org/).

---- 

`prestige` is heavily inspired by [stagehand](https://github.com/joelschutz/stagehand), but aims to simplify and clarify
some parts. It also comes with (way) fewer features, but should suffice, or be easily extendable, to accommodate other
use cases. This document contains some details about the structure of this package. Please read it for a better
understanding of how this package works; if anything here is unclear, it is considered a bug.

The main file is `scene-manager.go` which defines two crucial interfaces, `Scene[T]` and `Transition`. `fade.go` and
`simple.go` contain two simple implementations of `Transition`. `examples/basic` contains a simple example showing how
to use this package. You can run it using `go run .`.

The idea is as follows:

* You have a state struct, which has the same purpose as the usual `Game` struct in ebitengine. It should contain
state used across all scenes. It is also recommended to add a `SceneManager` field in your state. For example,
the state in `examples/basic` is just

```go
type State struct {
    Manager *prestige.SceneManager[State]
}
```
and it is initialized as 
```go
basicScene := &BasicScene{}
state := State{}
manager := prestige.NewSceneManager(basicScene, &state)
state.Manager = manager
```
Having the manager in your state will allow you to call methods on it, for instance, to do transitions, like so:
```go
// Fade transition for 1s to a new instance of BasicScene
state.Manager.Transition(&BasicScene{}, prestige.NewFadeTransition(1)) 
```

* For each scene you have a struct (say `BasicScene`) such that `*BasicScene` implements `Scene[T]` 
where `T` is the type of your state. The methods needed of `Scene[T]` contain those needed of `*Game` in ebitengine,
but with the additional `state *T` parameter passed in:
```go
Draw(state *T, screen *ebiten.Image)
Update(state *T) error
Layout(state *T, outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
```
Take a look at `examples/basic/main.go` for an example implementation of these methods. These implementations should
look very similar to analogous implementations of `*Game`, but now you have two sources of state: the first is the method
received `*BasicState`, and the second is the overall state passed to the method `*T`.

* But there are also four lifecycle methods 
```go
// Transition has started to enter your scene, load your stuff now
EnterStart(state *T)
// Transition has finished entering your scene, enable scene mechanics if they were disabled
EnterEnd(state *T)

// Transition has started to exit your scene, disable scene mechanics if you like
ExitStart(state *T)
// Transition has finished exiting your scene, wrap up now
ExitEnd(state *T)
```
The functionality here is clear. If you do not wish to do anything special to accommodate transitions, all four of these
would be empty. However, during a transition, you may wish to freeze inputs and such, which may be accomplished by
implementing these methods.

_When the first scene is loaded, `EnterStart` and `EnterEnd` are called in succession (see the implementation of 
`NewSceneManager`), so you may use them consistently to load dependencies_.

* Transitions are completely decoupled from the scenes, so they have no type parameter. They implement the interface
```go
// Called before the transition begins. Initialize the state
Start()

// Interpolate between two images produced by the two scenes between whom this transition exists
Interpolate(screen *ebiten.Image, src *ebiten.Image, dest *ebiten.Image)

// Update the state of the transition and return whether the transition has finished
Update() bool
```
The functionality is clear, but I will elaborate a bit on `Interpolate`. Essentially, interpolate is given the
images output by the current scene and the next scene (loading and initializing them are taken care of by the scene
manager). The job of this function is to just interpolate the two images anyway you see fit. For instance, 
see the implementation in `fade.go`. 

_If you really need the current state in the transition, you may query that from `SceneManager.State` and insert it into
the state for your Transition during construction_.

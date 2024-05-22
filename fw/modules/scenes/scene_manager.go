// SceneManager provides a manager for game scenes, automating the calling
// of their Init(), Deinit(), Draw(), ... functions,
// A SceneManager also features enabling/disabling, and ordering of scenes
// for drawing operations.
//
// Usage:
//    - Create a new SceneManager with `NewSceneManager`.
//    - Register scenes using `RegisterScene(name, scene)`.
//    - Control scene state with `EnableScene` and `DisableScene`.
//    - Modify draw order using `MoveSceneToFront`, `MoveSceneToBack`, and `MoveSceneBefore`.
//    - In the game loop, use `DrawScenes` and `DrawScenesGUI` to render scenes in their specified order.

package scenes

import (
	"gorl/fw/core/gem"
	"gorl/fw/core/logging"
	"gorl/fw/util"
)

type sceneManager struct {
	scenes         map[string]Scene
	enabled_scenes map[string]bool
	should_exit    bool
}

// Create a new SceneManager. A SceneManager will automatically take care of
// your Scenes (calling their Init(), Deinit(), Draw(), DrawGUI() functions).
func newSceneManager() *sceneManager {
	return &sceneManager{
		scenes:         make(map[string]Scene),
		enabled_scenes: make(map[string]bool),
	}
}

// The global instance of the SceneManager
var sm *sceneManager = newSceneManager()

// Register a scene with the SceneManager for automatic control
func RegisterScene(name string, scene Scene) {
	if _, exists := sm.scenes[name]; exists {
		logging.Fatal("A scene with name \"%v\" is already registered.", name)
	}
	sm.scenes[name] = scene
}

// Enable the Scene. The Scenes Init() function will be called.
func EnableScene(name string) {
	scene, exists := sm.scenes[name]
	if !exists {
		logging.Fatal("Scene with name \"%v\" not found.", name)
	}

	// Initialize the scene if it's not already enabled
	if !sm.enabled_scenes[name] {
		gem.AddEntity(gem.GetRoot(), scene.GetRoot(), gem.DefaultLayer)
		scene.Init()
		sm.enabled_scenes[name] = true
	}
}

// Disable the Scene. The Scenes Deinit() function will be called.
func DisableScene(name string) {
	scene, exists := sm.scenes[name]
	if !exists {
		logging.Fatal("Scene with name \"%v\" not found.", name)
	}

	// De-initialize the scene if it's currently enabled
	if sm.enabled_scenes[name] {
		scene.Deinit()
		gem.RemoveEntity(scene.GetRoot())
		sm.enabled_scenes[name] = false
	}
}

// Disable all Scenes that are currently enabled.
func DisableAllScenes() {
	for name, _ := range sm.scenes {
		if sm.enabled_scenes[name] {
			sm.scenes[name].Deinit()
			gem.RemoveEntity(sm.scenes[name].GetRoot())
			sm.enabled_scenes[name] = false
		}
	}
}

// Disable all Scenes that are currently enabled, except for the ones specified
// by name in the `exception_slice` parameter.
func DisableAllScenesExcept(exception_slice []string) {
	for name, _ := range sm.scenes {
		if sm.enabled_scenes[name] && !util.SliceContains(exception_slice, name) {
			sm.scenes[name].Deinit()
			sm.enabled_scenes[name] = false
		}
	}
}

// Calls the Update() functions of all the registered Scenes
func UpdateScenes() {
	for name, _ := range sm.scenes {
		if sm.enabled_scenes[name] {
			sm.scenes[name].Update()
		}
	}
}

// Calls the FixedUpdate() functions of all the registered Scenes
func FixedUpdateScenes() {
	for name, _ := range sm.scenes {
		if sm.enabled_scenes[name] {
			sm.scenes[name].FixedUpdate()
		}
	}
}

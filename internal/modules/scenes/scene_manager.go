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
	"gorl/internal/logging"
	"gorl/internal/util"
)

type SceneManager struct {
	scenes         map[string]Scene
	enabled_scenes map[string]bool
	scene_order    []string // slice to maintain order, since map is unordered
	should_exit    bool
}

// Create a new SceneManager. A SceneManager will automatically take care of
// your Scenes (calling their Init(), Deinit(), Draw(), DrawGUI() functions).
func newSceneManager() *SceneManager {
	return &SceneManager{
		scenes:         make(map[string]Scene),
		enabled_scenes: make(map[string]bool),
		scene_order:    make([]string, 0),
	}
}

// The global instance of the SceneManager
var Sm *SceneManager = newSceneManager()

// Register a scene with the SceneManager for automatic control
func (sm *SceneManager) RegisterScene(name string, scene Scene) {
	if _, exists := sm.scenes[name]; exists {
		logging.Fatal("A scene with name \"%v\" is already registered.", name)
	}
	sm.scenes[name] = scene
	sm.scene_order = append(sm.scene_order, name) // Add to the end by default
}

// MoveSceneToFront moves the scene to the front of the draw order
func (sm *SceneManager) MoveSceneToFront(name string) {
	sm.reorderScene(name, 0)
}

// MoveSceneToBack moves the scene to the end of the draw order
func (sm *SceneManager) MoveSceneToBack(name string) {
	sm.reorderScene(name, len(sm.scene_order)-1)
}

// MoveSceneBefore moves the scene right before another scene in the draw order
func (sm *SceneManager) MoveSceneBefore(sceneName, beforeSceneName string) {
	index, found := sm.getSceneOrderIndex(beforeSceneName)
	if found {
		sm.reorderScene(sceneName, index)
	}
}

func (sm *SceneManager) reorderScene(name string, index int) {
	current_idx, found := sm.getSceneOrderIndex(name)
	if !found {
		return
	}
	sm.scene_order = append(sm.scene_order[:current_idx], sm.scene_order[current_idx+1:]...)
	sm.scene_order = append(sm.scene_order[:index], append([]string{name}, sm.scene_order[index:]...)...)
}

func (sm *SceneManager) getSceneOrderIndex(name string) (int, bool) {
	for i, scene_name := range sm.scene_order {
		if scene_name == name {
			return i, true
		}
	}
	return -1, false
}

// Enable the Scene. The Scenes Init() function will be called.
func (sm *SceneManager) EnableScene(name string) {
	scene, exists := sm.scenes[name]
	if !exists {
		logging.Fatal("Scene with name \"%v\" not found.", name)
	}

	// Initialize the scene if it's not already enabled
	if !sm.enabled_scenes[name] {
		scene.Init()
		sm.enabled_scenes[name] = true
	}
}

// Disable the Scene. The Scenes Deinit() function will be called.
func (sm *SceneManager) DisableScene(name string) {
	scene, exists := sm.scenes[name]
	if !exists {
		logging.Fatal("Scene with name \"%v\" not found.", name)
	}

	// De-initialize the scene if it's currently enabled
	if sm.enabled_scenes[name] {
		scene.Deinit()
		sm.enabled_scenes[name] = false
	}
}

// Disable all Scenes that are currently enabled.
func (sm *SceneManager) DisableAllScenes() {
	for _, name := range sm.scene_order {
		if sm.enabled_scenes[name] {
			sm.scenes[name].Deinit()
			sm.enabled_scenes[name] = false
		}
	}
}

// Disable all Scenes that are currently enabled, except for the ones specified
// by name in the `exception_slice` parameter.
func (sm *SceneManager) DisableAllScenesExcept(exception_slice []string) {
	for _, name := range sm.scene_order {
		if sm.enabled_scenes[name] && !util.SliceContains(exception_slice, name) {
			sm.scenes[name].Deinit()
			sm.enabled_scenes[name] = false
		}
	}
}

// Call the Draw() functions of all the registered Scenes in their defined order.
func (sm *SceneManager) DrawScenes() {
	for _, name := range sm.scene_order {
		if sm.enabled_scenes[name] {
			sm.scenes[name].Draw()
		}
	}
}

// Call the DrawGUI() functions of all the registered Scenes in their defined order.
func (sm *SceneManager) DrawScenesGUI() {
	for _, name := range sm.scene_order {
		if sm.enabled_scenes[name] {
			sm.scenes[name].DrawGUI()
		}
	}
}

// TODO: this should not be part of the scene manager, come up with a better solution
func (sm *SceneManager) ExitGame() {
	sm.should_exit = true
}
func (sm *SceneManager) ShouldExit() bool {
	return sm.should_exit
}

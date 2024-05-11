package main

import (
	"fmt"
	//"gorl/internal/audio"
	//"gorl/internal/physics"

	//"gorl/internal/collision"
	"gorl/internal/core/gem"
	//"gorl/pkg/entities"
	//"gorl/internal/gui"

	//"gorl/internal/lighting"
	"gorl/internal/logging"
	//"gorl/internal/physics"
	"gorl/internal/core/render"
	"gorl/internal/modules/scenes"
	"gorl/internal/settings"
	uscenes "gorl/pkg/scenes"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// PRE-INIT

	// settings
	settings_path := "settings.json"
	err := settings.LoadSettings(settings_path)
	if err != nil {
		fmt.Println("Error loading settings:", err)
		fmt.Println("Using fallback settings.")
		settings.UseFallbackSettings()
	}

	// logging
	logging.Init(settings.CurrentSettings().LogPath)
	logging.Info("Logging initialized")
	if err == nil {
		logging.Info("Settings loaded successfully.")
	} else {
		logging.Warning("Settings loading unsuccessful, using fallback.")
	}

	// INITIALIZATION
	// raylib window
	rl.InitWindow(
		int32(settings.CurrentSettings().ScreenWidth),
		int32(settings.CurrentSettings().ScreenHeight),
		settings.CurrentSettings().Title)
	defer rl.CloseWindow()
	rl.SetTargetFPS(int32(settings.CurrentSettings().TargetFps))

	// rendering
	renderSystem := render.NewRenderSystem(rl.NewVector2(
		float32(settings.CurrentSettings().ScreenWidth),
		float32(settings.CurrentSettings().ScreenHeight)))

	renderRatio := float32(settings.CurrentSettings().RenderWidth) / float32(settings.CurrentSettings().ScreenWidth)
	// renders at default resolution
	defaultRenderStage := render.NewRenderStage(rl.NewVector2(
		float32(settings.CurrentSettings().RenderWidth),
		float32(settings.CurrentSettings().RenderHeight)), renderRatio)

	// renders at double resolution
	//doubleResRenderStage := render.NewRenderStage(rl.NewVector2(
	//	float32(settings.CurrentSettings().RenderWidth*2),
	//	float32(settings.CurrentSettings().RenderHeight*2)), 2)

	logging.Info("Custom rendering initialized.")

	// initialize audio
	//audio.InitAudio()
	//defer audio.DeinitAudio()

	// collision
	//collision.InitCollision()
	//defer collision.DeinitCollision()

	// physics
	//physics.InitPhysics((1.0 / 60.0), rl.Vector2Zero(), (1.0 / 32.0))
	//defer physics.DeinitPhysics()

	//gem.InitGem(physics.GetTimestep())
	//gem.InitGem(0)
	gem.Init()
	defer gem.Deinit()

	// lighting
	//lighting.InitLighting()
	//defer lighting.DeinitLighting()

	// animtion (premades need init and update)
	//animation.InitPremades(render.Rs.CurrentStage.Camera, render.GetWSCameraOffset())

	// register audio tracks
	//audio.RegisterMusic("aza-tumbleweeds", "audio/music/azakaela/azaFMP2_field7_Tumbleweeds.ogg")
	//audio.RegisterMusic("aza-outwest", "audio/music/azakaela/azaFMP2_scene1_OutWest.ogg")
	//audio.RegisterMusic("aza-frontier", "audio/music/azakaela/azaFMP2_town_Frontier.ogg")
	//audio.CreatePlaylist("main-menu", []string{"aza-tumbleweeds", "aza-outwest", "aza-frontier"})
	//audio.SetGlobalVolume(0.9)
	//audio.SetMusicVolume(0.7)
	//audio.SetSFXVolume(0.9)

	// gui
	//gui.InitBackend()

	// cursor
	//rl.HideCursor()

	// scenes
	//scenes.Sm.RegisterScene("dev", &scenes.DevScene{})
	//scenes.Sm.EnableScene("dev")

	scenes.RegisterScene("some_name", &uscenes.TemplateScene{})
	scenes.EnableScene("some_name")
	scenes.DisableScene("some_name")

	//rl.DisableCursor()

	// GAME LOOP
	//rl.SetExitKey(rl.KeyEnd) // Set a key to exit the game
	shouldExit := false
	for !shouldExit {
		//animation.UpdatePremades()
		//render.UpdateEffects()

		scenes.UpdateScenes()
		// scenes.FixedUpdateScenes()

		rl.BeginDrawing()

		// begin drawing the world
		//render.BeginCustomRenderWorldspace()

		rl.ClearBackground(rl.RayWhite)

		renderSystem.EnableRenderStage(defaultRenderStage)
		rl.ClearBackground(rl.Blank)

		gem.Draw(gem.GetByLayer(gem.DefaultLayer))

		//renderSystem.EnableRenderStage(doubleResRenderStage)
		//rl.ClearBackground(rl.Blank)

		//gem.Draw(gem.GetByLayer(gem.DefaultLayer + 1))

		renderSystem.FlushRenderStage()
		renderSystem.RenderToScreen()

		//rl.DrawCircleV(rl.NewVector2(100, 300), 50, rl.Blue)

		// Draw all registered Scenes
		//gem.UpdateEntities()
		//gem.DrawEntities()
		//scenes.Sm.DrawScenes()
		//collision.Update()
		//physics.Update()

		// lighting
		//lighting.DrawLight()

		//physics.DrawColliders(true, false, false)

		// begin drawing the gui
		//render.BeginCustomRenderScreenspace()

		//rl.ClearBackground(rl.Blank) // clear the main rendertex

		//scenes.Sm.DrawScenesGUI()
		//gem.DrawEntitiesGUI()

		//render.EndCustomRender()
		//mousecursor.Draw()

		// Draw Debug Info
		//rl.DrawFPS(10, 10)
		//render.DebugDrawStageViewports(
		//	rl.NewVector2(10, 10), 4, renderSystem,
		//	[]*render.RenderStage{defaultRenderStage, doubleResRenderStage},
		//)

		rl.EndDrawing()

		//audio.Update()
		shouldExit = rl.WindowShouldClose() // || scenes.Sm.ShouldExit()
	}

	//scenes.Sm.DisableAllScenes()
	//gem.RemoveEntity(gem.Root())
}

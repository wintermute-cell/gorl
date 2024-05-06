package main

import (
	"fmt"
	//"gorl/internal/audio"
	//"gorl/internal/physics"

	//"gorl/internal/collision"
	"gorl/internal/core/entities/gem"
	//"gorl/internal/gui"

	//"gorl/internal/lighting"
	"gorl/internal/logging"
	//"gorl/internal/physics"
	"gorl/internal/core/render"
	"gorl/internal/modules/scenes"
	"gorl/internal/settings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// PRE-INIT

	// TODO: impl this
	//now := time.Now()
	//i := messages.ImplementsInterfaceGeneric[MyInterface](&MyStruct{})
	//fmt.Println(i) // should
	//fmt.Println("Time taken:", time.Since(now))

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

	// renders at default resolution
	defaultRenderStage := render.NewRenderStage(rl.NewVector2(
		float32(settings.CurrentSettings().RenderWidth),
		float32(settings.CurrentSettings().RenderHeight)), 1)

	// renders at double resolution
	doubleResRenderStage := render.NewRenderStage(rl.NewVector2(
		float32(settings.CurrentSettings().RenderWidth*2),
		float32(settings.CurrentSettings().RenderHeight*2)), 2)
	doubleResRenderStage.SetCameraOffset(rl.NewVector2(
		float32(settings.CurrentSettings().RenderWidth/2),
		float32(settings.CurrentSettings().RenderHeight/2)))

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
	gem.InitGem(0)

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

	//rl.DisableCursor()

	// GAME LOOP
	//rl.SetExitKey(rl.KeyEnd) // Set a key to exit the game
	shouldExit := false
	for !shouldExit {
		//animation.UpdatePremades()
		//render.UpdateEffects()

		rl.BeginDrawing()

		// begin drawing the world
		//render.BeginCustomRenderWorldspace()

		if rl.IsKeyDown(rl.KeyJ) {
			curretCameraTarget := doubleResRenderStage.GetCameraTarget()
			doubleResRenderStage.SetCameraTarget(rl.Vector2Add(curretCameraTarget, rl.NewVector2(0, 10)))
		}

		if rl.IsKeyDown(rl.KeyK) {
			curretCameraTarget := doubleResRenderStage.GetCameraTarget()
			doubleResRenderStage.SetCameraTarget(rl.Vector2Add(curretCameraTarget, rl.NewVector2(0, -10)))
		}

		if rl.IsKeyDown(rl.KeyH) {
			curretCameraTarget := doubleResRenderStage.GetCameraTarget()
			doubleResRenderStage.SetCameraTarget(rl.Vector2Add(curretCameraTarget, rl.NewVector2(-10, 0)))
		}

		if rl.IsKeyDown(rl.KeyL) {
			curretCameraTarget := doubleResRenderStage.GetCameraTarget()
			doubleResRenderStage.SetCameraTarget(rl.Vector2Add(curretCameraTarget, rl.NewVector2(10, 0)))
		}

		rl.ClearBackground(rl.RayWhite)

		renderSystem.EnableRenderStage(defaultRenderStage)
		rl.ClearBackground(rl.Blank)
		rl.DrawCircleV(rl.NewVector2(100, 100), 50, rl.Red)

		renderSystem.EnableRenderStage(doubleResRenderStage)
		rl.ClearBackground(rl.Blank)
		// mark each corner with a circle
		rl.DrawCircleV(rl.NewVector2(0, 0), 50, rl.Green)
		rl.DrawCircleV(rl.NewVector2(0, float32(settings.CurrentSettings().RenderHeight)), 50, rl.Green)
		rl.DrawCircleV(rl.NewVector2(float32(settings.CurrentSettings().RenderWidth), 0), 50, rl.Green)
		rl.DrawCircleV(rl.NewVector2(float32(settings.CurrentSettings().RenderWidth), float32(settings.CurrentSettings().RenderHeight)), 50, rl.Green)
		// and the center
		rl.DrawCircleV(rl.NewVector2(float32(settings.CurrentSettings().RenderWidth/2), float32(settings.CurrentSettings().RenderHeight/2)), 50, rl.Blue)

		// and at the camera target
		rl.DrawCircleV(doubleResRenderStage.GetCameraTarget(), 50, rl.Purple)

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
		rl.DrawFPS(10, 10)
		render.DebugDrawStageViewports(
			rl.NewVector2(10, 10), 4, renderSystem,
			[]*render.RenderStage{defaultRenderStage, doubleResRenderStage},
		)

		rl.EndDrawing()

		//audio.Update()
		shouldExit = rl.WindowShouldClose() || scenes.Sm.ShouldExit()
	}

	//scenes.Sm.DisableAllScenes()
	//gem.RemoveEntity(gem.Root())
}

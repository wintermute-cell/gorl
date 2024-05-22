package main

import (
	"fmt"
	"time"

	"gorl/fw/core/entities/proto"
	"gorl/fw/core/gem"
	input "gorl/fw/core/input/input_handling"
	"gorl/fw/core/logging"
	"gorl/fw/core/render"
	"gorl/fw/core/settings"
	"gorl/fw/modules/scenes"
	"gorl/game"
	"gorl/game/entities"

	rl "github.com/gen2brain/raylib-go/raylib"

	"net/http"
	_ "net/http/pprof"
)

func main() {
	// PRE-INIT
	go func() {
		err := http.ListenAndServe("localhost:6969", nil)
		if err != nil {
			panic(err)
		}
	}()

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
	render.Init(rl.NewVector2(
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

	//scenes.RegisterScene("some_name", &uscenes.TemplateScene{})
	//scenes.EnableScene("some_name")
	//scenes.DisableScene("some_name")

	//rl.DisableCursor()
	game.Init()

	// GAME LOOP
	//rl.SetExitKey(rl.KeyEnd) // Set a key to exit the game
	shouldExit := false

	// frame time measurement stuff
	now := time.Now()
	var avgTime time.Duration

	button1 := entities.NewButtonEntity2D(rl.NewVector2(100, 100), 0, rl.NewVector2(1, 1))
	button1.Name = "Button 1"
	button2 := entities.NewButtonEntity2D(rl.NewVector2(500, 100), 0, rl.NewVector2(1, 1))
	button2.Name = "Button 2"
	gem.AddEntity(gem.GetRoot(), button1, gem.DefaultLayer)
	gem.AddEntity(gem.GetRoot(), button2, gem.DefaultLayer)

	// keeps a list of entities in draw order, so we can pass input event in the correct order.
	orderedEntities := [][]proto.IEntity{}
	for !shouldExit {
		now = time.Now()

		orderedEntities = [][]proto.IEntity{}

		//animation.UpdatePremades()
		//render.UpdateEffects()

		game.Update()
		scenes.UpdateScenes()
		// scenes.FixedUpdateScenes()

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		render.EnableRenderStage(defaultRenderStage)
		rl.ClearBackground(rl.Blank)

		orderedEntities = append(orderedEntities, gem.DrawLayer(gem.DefaultLayer))

		//render.EnableRenderStage(doubleResRenderStage)
		//rl.ClearBackground(rl.Blank)
		//gem.Draw(gem.GetByLayer(gem.DefaultLayer + 1))

		render.FlushRenderStage()
		render.RenderToScreen()

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

		// input is processed at the end of the frame, because here we know in
		// what order the entities were drawn, and can be sure whatever the
		// user clicked was really visible at the front.
		input.HandleInputEvents(orderedEntities)

		// Draw Debug Info
		rl.DrawFPS(10, 10)
		rl.DrawText("dt: "+avgTime.String(), 10, 30, 20, rl.Lime)
		//render.DebugDrawStageViewports(
		//	rl.NewVector2(10, 10), 4, render,
		//	[]*render.RenderStage{defaultRenderStage},
		//)
		//gem.DebugDrawEntities(rl.NewVector2(10, 50), 12)
		gem.DebugDrawHierarchy(rl.NewVector2(10, 50), 8)

		rl.EndDrawing()

		//audio.Update()
		shouldExit = rl.WindowShouldClose() // || scenes.Sm.ShouldExit()

		// calculate the time it took to render the frame.
		// we must do this after the frame is drawn.
		if avgTime == 0 {
			avgTime = time.Since(now)
		} else {
			avgTime = (avgTime + time.Since(now)) / 2
		}

	}

	//scenes.Sm.DisableAllScenes()
}

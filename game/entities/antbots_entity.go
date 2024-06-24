package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"
	"gorl/fw/core/math"
	"gorl/fw/core/settings"
	"gorl/fw/util"
	"gorl/game/code"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that AntbotsEntity implements IEntity.
var _ entities.IEntity = &AntbotsEntity{}

// Antbots Entity
type AntbotsEntity struct {
	*entities.Entity // Required!

	// home base
	homeBasePoint  rl.Vector2
	homeBaseRadius float32

	// bots
	botSprite   rl.Texture2D
	bots        []*code.Antbot
	obstacleMap *code.ObstacleMap
	decayTimer  *util.Timer
	mapTexture  rl.Texture2D

	pathTex rl.RenderTexture2D

	cam *CameraEntity
}

// NewAntbotsEntity creates a new instance of the AntbotsEntity.
func NewAntbotsEntity(botAmount int, spawnPoint rl.Vector2, spawnRadius float32, cam *CameraEntity, foodpiles *FoodpilesEntity) *AntbotsEntity {
	new_ent := &AntbotsEntity{
		Entity:         entities.NewEntity("AntbotsEntity", rl.Vector2Zero(), 0, rl.Vector2One()),
		homeBasePoint:  spawnPoint,
		homeBaseRadius: spawnRadius,

		botSprite:   rl.LoadTexture("antbot.png"),
		bots:        make([]*code.Antbot, botAmount),
		obstacleMap: code.NewObstacleMap(math.NewVector2Int(1920/2, 1080/2), math.NewVector2Int(1920, 1080)),
		decayTimer:  util.NewTimer(40.0 / 200.0), // 20 seconds to decay to 0
		pathTex:     rl.LoadRenderTexture(1920, 1080),
		cam:         cam,
	}

	mapSize := new_ent.obstacleMap.GetSize()
	emptyImg := rl.GenImageColor(mapSize.X, mapSize.Y, rl.Blank)
	new_ent.mapTexture = rl.LoadTextureFromImage(emptyImg)
	rl.UnloadImage(emptyImg)

	for i := 0; i < botAmount; i++ {
		randPoint := rl.NewVector2(
			spawnPoint.X+((rand.Float32()-0.5)*2)*spawnRadius,
			spawnPoint.Y+((rand.Float32()-0.5)*2)*spawnRadius,
		)
		pointAngle := -util.Vector2Angle(rl.Vector2Subtract(randPoint, spawnPoint), rl.NewVector2(0, -1)) * rl.Rad2deg
		new_ent.bots[i] = code.NewAntbot(randPoint, pointAngle, new_ent.obstacleMap, foodpiles.Foodpiles, spawnPoint, spawnRadius)
	}

	return new_ent
}

func (ent *AntbotsEntity) Init() {
	// Initialization logic for the entity
	// ...
}

func (ent *AntbotsEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *AntbotsEntity) Update() {
	if ent.decayTimer.Check() {
		ent.obstacleMap.DecayPheromones(1)
	}
	for _, bot := range ent.bots {
		bot.Move()
	}
}

func (ent *AntbotsEntity) Draw() {

	ent.obstacleMap.ToRlTexture(ent.mapTexture)
	rl.DrawTexturePro(
		ent.mapTexture,
		rl.NewRectangle(0, 0, float32(ent.mapTexture.Width), float32(ent.mapTexture.Height)),
		rl.NewRectangle(0, 0, float32(settings.CurrentSettings().RenderWidth), float32(settings.CurrentSettings().RenderHeight)),
		rl.NewVector2(0, 0),
		0, rl.White,
	)

	rl.DrawCircleV(ent.homeBasePoint, ent.homeBaseRadius, rl.Fade(rl.Green, 0.7))
	// draw new pixels to the path texture
	//pt := ent.cam.GetCamera().GetTexture()
	//rl.EndTextureMode()
	//rl.BeginTextureMode(ent.pathTex)
	//for _, bot := range ent.bots {
	//	// Draw the bot's position
	//	rl.DrawPixel(int32(bot.Transform.GetPosition().X), int32(bot.Transform.GetPosition().Y), rl.Fade(rl.Red, 0.9))
	//}
	//rl.EndTextureMode()
	//rl.BeginTextureMode(pt)

	// Draw the path to the screen
	rl.DrawTexturePro(
		ent.pathTex.Texture,
		rl.NewRectangle(0, 0, float32(ent.pathTex.Texture.Width), -float32(ent.pathTex.Texture.Height)),
		rl.NewRectangle(0, 0, float32(ent.pathTex.Texture.Width), float32(ent.pathTex.Texture.Height)),
		rl.NewVector2(0, 0),
		0,
		rl.White,
	)

	for _, bot := range ent.bots {
		// Draw the bot
		sprtDims := rl.NewVector2(float32(ent.botSprite.Width), float32(ent.botSprite.Height))
		rl.DrawTexturePro(
			ent.botSprite,
			rl.NewRectangle(0, 0, sprtDims.X, sprtDims.Y),
			rl.NewRectangle(bot.Transform.GetPosition().X, bot.Transform.GetPosition().Y, sprtDims.X, sprtDims.Y),
			rl.NewVector2(sprtDims.X/2, sprtDims.Y/2),
			bot.Transform.GetRotation(),
			rl.White,
		)
		if bot.BotMode == code.BotModeReturning {
			rl.DrawCircleV(bot.Transform.GetPosition(), 3, rl.Green)
		}
	}
}

func (ent *AntbotsEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...
	return true
}

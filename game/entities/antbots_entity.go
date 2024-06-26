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
	botSprite       rl.Texture2D
	bots            []*code.Antbot
	obstacleMap     *code.PheromoneMap
	decayTimer      *util.Timer
	pheromoneMapTex rl.Texture2D
	obstacleMapTex  rl.Texture2D

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
		obstacleMap: code.NewPheromoneMap(math.NewVector2Int(1920/3, 1080/3), math.NewVector2Int(1920, 1080)),
		decayTimer:  util.NewTimer(25.0 / 200.0), // 30 seconds to decay to 0
		pathTex:     rl.LoadRenderTexture(1920, 1080),
		cam:         cam,
	}

	mapSize := new_ent.obstacleMap.GetSize()
	emptyImg := rl.GenImageColor(mapSize.X, mapSize.Y, rl.Blank)
	new_ent.pheromoneMapTex = rl.LoadTextureFromImage(emptyImg)
	new_ent.obstacleMapTex = rl.LoadTextureFromImage(emptyImg)
	rl.UnloadImage(emptyImg)

	for i := 0; i < botAmount; i++ {
		randPoint := rl.NewVector2(
			spawnPoint.X+((rand.Float32()-0.5)*2)*5,
			spawnPoint.Y+((rand.Float32()-0.5)*2)*5,
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

	rl.ClearBackground(rl.Black)
	ent.obstacleMap.ObstaclesToTexture(ent.obstacleMapTex)
	rl.DrawTexturePro(
		ent.obstacleMapTex,
		rl.NewRectangle(0, 0, float32(ent.obstacleMapTex.Width), float32(ent.obstacleMapTex.Height)),
		rl.NewRectangle(0, 0, float32(settings.CurrentSettings().RenderWidth), float32(settings.CurrentSettings().RenderHeight)),
		rl.NewVector2(0, 0),
		0, rl.White,
	)

	ent.obstacleMap.PheromoneToTexture(ent.pheromoneMapTex)
	rl.DrawTexturePro(
		ent.pheromoneMapTex,
		rl.NewRectangle(0, 0, float32(ent.pheromoneMapTex.Width), float32(ent.pheromoneMapTex.Height)),
		rl.NewRectangle(0, 0, float32(settings.CurrentSettings().RenderWidth), float32(settings.CurrentSettings().RenderHeight)),
		rl.NewVector2(0, 0),
		0, rl.White,
	)

	rl.DrawCircleV(ent.homeBasePoint, ent.homeBaseRadius, rl.Fade(rl.Green, 0.7))

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
		tint := rl.RayWhite
		if bot.BotMode == code.BotModeReturning {
			tint = rl.Lime
		}
		rl.DrawTexturePro(
			ent.botSprite,
			rl.NewRectangle(0, 0, sprtDims.X, sprtDims.Y),
			rl.NewRectangle(bot.Transform.GetPosition().X, bot.Transform.GetPosition().Y, sprtDims.X/2, sprtDims.Y/2),
			rl.NewVector2(sprtDims.X/4, sprtDims.Y/4),
			bot.Transform.GetRotation(),
			tint,
		)
	}
}

func (ent *AntbotsEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...
	return true
}

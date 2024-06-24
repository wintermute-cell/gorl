package code

import (
	"gorl/fw/core/math"
	"gorl/fw/util"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type mode int32

const (
	BotModeLeaving mode = iota
	BotModeReturning
	BotModeNoFood
)

type Antbot struct {
	Transform *math.Transform2D

	steerForce float32 // How fast the bot can rotate/steer

	maxSpeed float32 // Maximum speed the bot can reach

	isTurning bool
	turnTimer *util.Timer
	turnDir   rl.Vector2

	rotVelocity float32
	linVelocity float32

	botRadius float32

	sensorPoints         []rl.Vector2
	pheroSensorDist      float32
	pheroSensorRadius    float32
	obstacleSensorDist   float32
	obstacleSensorRadius float32

	obstacleMap *ObstacleMap
	foodpiles   *FoodPiles

	homeBasePoint  rl.Vector2
	homeBaseRadius float32

	BotMode        mode
	pheromoneTimer *util.Timer

	// wander
	wanderTheta float32

	// wall avoidance
	avoidanceStrength float32
	avoidanceDir      rl.Vector2
}

func NewAntbot(position rl.Vector2, rotation float32, obstacleMap *ObstacleMap, foodpiles *FoodPiles, homeBasePoint rl.Vector2, homeBaseRadius float32) *Antbot {
	tf := math.NewTransform2D(position, rotation, rl.Vector2One())

	sensorPoints := []rl.Vector2{}
	visionAngle := float32(90.0)
	sensorCount := 3 // don't changes this, we rely on it being 3 by this point
	anglePerSensor := (visionAngle / float32(sensorCount-1)) * rl.Deg2rad
	centerOffset := (visionAngle / 2) * rl.Deg2rad
	for i := 0; i < sensorCount; i++ {
		sensorPoints = append(sensorPoints, rl.Vector2Rotate(rl.NewVector2(0, -1), anglePerSensor*float32(i)-centerOffset))
	}

	return &Antbot{
		Transform: &tf,

		BotMode:        BotModeLeaving,
		pheromoneTimer: util.NewTimer(0.05),

		homeBasePoint:  homeBasePoint,
		homeBaseRadius: homeBaseRadius,

		isTurning: false,
		turnTimer: util.NewTimer(0.5),

		botRadius: 4,

		maxSpeed: 100,

		rotVelocity: 0,
		linVelocity: 0,

		sensorPoints:         sensorPoints,
		pheroSensorDist:      30,
		pheroSensorRadius:    20,
		obstacleSensorDist:   8,
		obstacleSensorRadius: 4,

		obstacleMap: obstacleMap,
		foodpiles:   foodpiles,
	}
}

func (bot *Antbot) Move() {
	//bot.linVelocity = util.Clamp(bot.linVelocity, 0, bot.maxSpeed)

	bot.rotVelocity = 0
	bot.linVelocity = bot.maxSpeed
	rot, lin := float32(0), float32(0)

	switch bot.BotMode {
	case BotModeLeaving:
		if bot.pheromoneTimer.Check() {
			bot.obstacleMap.SetPheromone(
				math.Vector2IntFromRl(bot.Transform.GetPosition()),
				PheromoneLeaving,
			)
		}
		foundFood := bot.detectAndTakeFood()
		if foundFood {
			bot.isTurning = true
			bot.turnDir = bot.Transform.Down()
			bot.turnTimer.ResetTime()
			bot.BotMode = BotModeReturning
		}

		wanderStrength := float32(2)
		rot, lin = bot.steerWander()
		bot.rotVelocity += rot * wanderStrength
		bot.linVelocity += lin * wanderStrength

		trackStrength := float32(1.9)
		rot, lin = bot.steerTrackPheromone(PheromoneReturning)
		bot.rotVelocity += rot * trackStrength
		bot.linVelocity += lin

	case BotModeReturning:
		if bot.pheromoneTimer.Check() {
			bot.obstacleMap.SetPheromone(
				math.Vector2IntFromRl(bot.Transform.GetPosition()),
				PheromoneReturning,
			)
		}
		wanderStrength := float32(2)
		rot, lin = bot.steerWander()
		bot.rotVelocity += rot * wanderStrength
		bot.linVelocity += lin

		trackStrength := float32(1.9)
		rot, lin = bot.steerTrackPheromone(PheromoneLeaving)
		bot.rotVelocity += rot * trackStrength
		bot.linVelocity += lin

		if rl.CheckCollisionCircles(bot.Transform.GetPosition(), bot.botRadius, bot.homeBasePoint, bot.homeBaseRadius) {
			bot.isTurning = true
			bot.turnDir = bot.Transform.Down()
			bot.turnTimer.ResetTime()
			bot.BotMode = BotModeLeaving
		}

	case BotModeNoFood:
		//bot.botMode = modeLeaving
	}

	if bot.isTurning {
		const turnStrength = 10
		rot, lin = bot.steerTurn()
		if bot.turnTimer.Check() {
			bot.isTurning = false
		}
		bot.rotVelocity += rot * turnStrength
		bot.linVelocity += 0
	}

	wallAvoidanceStrength := float32(20)
	rot, lin = bot.steerWallAvoidance()
	bot.rotVelocity += rot * wallAvoidanceStrength
	bot.linVelocity += lin

	// ensure the bot doesn't go too fast
	bot.linVelocity = util.Clamp(bot.linVelocity, 0, bot.maxSpeed)

	bot.Transform.AddRotation(bot.rotVelocity * rl.Rad2deg * rl.GetFrameTime())
	bot.Transform.AddPosition(rl.Vector2Scale(bot.Transform.Up(), bot.linVelocity*rl.GetFrameTime()))

	bot.wrapPosition()

}

func (bot *Antbot) steerSeek(target rl.Vector2) (float32, float32) {
	desiredVelocity := rl.Vector2Subtract(
		target,
		bot.Transform.GetPosition(),
	)
	steerAngle := util.Vector2Angle(bot.Transform.Up(), desiredVelocity)
	steerMove := rl.Vector2Length(desiredVelocity) * bot.maxSpeed

	return steerAngle, steerMove
}

func (bot *Antbot) steerWander() (float32, float32) {
	wanderLen := float32(100.0)
	wanderSector := float32(45.0)
	wanderRadius := float32(90.0)
	center := rl.Vector2Add(bot.Transform.GetPosition(), rl.Vector2Scale(bot.Transform.Up(), wanderLen))

	bot.wanderTheta += float32((rand.Float32()-0.5)*2) * wanderSector * rl.GetFrameTime()
	wanderDir := rl.Vector2Rotate(bot.Transform.Up(), bot.wanderTheta)
	wanderTarget := rl.Vector2Add(center, rl.Vector2Scale(wanderDir, wanderRadius))

	//rl.DrawCircleV(center, 5, rl.Blue)
	//rl.DrawCircleV(wanderTarget, 5, rl.Red)
	rot, _ := bot.steerSeek(wanderTarget)

	return rot, 0
}

func (bot *Antbot) steerTrackPheromone(pheromoneType PheromoneType) (float32, float32) {
	rot, lin := float32(0), float32(0)

	// a point in time that is guaranteed to be before any pheromone
	concentrations := []float32{0, 0, 0}
	points := []rl.Vector2{}

	found := false
	for idx, sensor := range bot.sensorPoints {
		rotPoint := rl.Vector2Rotate(rl.Vector2Scale(sensor, bot.pheroSensorDist), bot.Transform.GetRotation()*rl.Deg2rad)
		absPoint := rl.Vector2Add(
			bot.Transform.GetPosition(),
			rotPoint,
		)

		count, agedCount := bot.obstacleMap.HasInCircle(
			math.Vector2IntFromRl(absPoint),
			bot.pheroSensorRadius,
			pheromoneType,
		)

		if count > 0 { // use the integer value to check if there is any pheromone at all
			found = true
			concentrations[idx] = agedCount
		}
		points = append(points, absPoint)
	}

	// use the sensor with the freshest and most hits as the direction to follow
	//const freshWeight = 1
	//const hitWeight = 1

	// use the sensor with the freshest pheromone, prefer middle sensor
	if found {
		// if center > left and center > right
		if concentrations[1] > concentrations[0] && concentrations[1] > concentrations[2] {
			rot, lin = bot.steerSeek(points[1])
		} else if concentrations[0] > concentrations[2] {
			// if left > right
			rot, lin = bot.steerSeek(points[0])
		} else if concentrations[2] > concentrations[0] {
			// if right > left
			rot, lin = bot.steerSeek(points[2])
		}
	}

	return rot, lin
}

func (bot *Antbot) steerTurn() (float32, float32) {
	lin := float32(0)
	rot := float32(0)
	turnPoint := rl.Vector2Add(bot.Transform.GetPosition(), rl.Vector2Scale(bot.turnDir, 10))
	rl.DrawCircleV(turnPoint, 5, rl.Red)
	rot, lin = bot.steerSeek(turnPoint)
	return rot, lin
}

func (bot *Antbot) steerWallAvoidance() (float32, float32) {
	bot.avoidanceStrength = max(0, bot.avoidanceStrength-0.2*rl.GetFrameTime())
	if bot.avoidanceStrength <= 0 {
		bot.avoidanceDir = rl.Vector2Zero()
	}
	mustAvoid := float32(0)
	for _, sensor := range bot.sensorPoints {
		rotPoint := rl.Vector2Rotate(rl.Vector2Scale(sensor, bot.obstacleSensorDist), bot.Transform.GetRotation()*rl.Deg2rad)
		absPoint := rl.Vector2Add(
			bot.Transform.GetPosition(),
			rotPoint,
		)
		// Drawing the sensor points
		//rl.DrawCircleV(absPoint, bot.sensorRadius, rl.Red)

		hits, _ := bot.obstacleMap.HasInCircle(
			math.Vector2IntFromRl(absPoint),
			bot.obstacleSensorRadius,
			PheromoneEdgeObstacle,
		)
		if hits > 0 {
			awayDir := rl.Vector2Scale(rotPoint, -1)
			bot.avoidanceDir = rl.Vector2Add(bot.avoidanceDir, awayDir)
			bot.avoidanceStrength += 1 * rl.GetFrameTime()
			mustAvoid = -1
		}
	}

	if bot.avoidanceStrength > 0.0001 {
		// Drawing the avoid direction
		//rl.DrawLineV(bot.Transform.GetPosition(), rl.Vector2Add(bot.Transform.GetPosition(), bot.avoidanceDir), rl.Red)

		rot, _ := bot.steerSeek(rl.Vector2Add(bot.Transform.GetPosition(), bot.avoidanceDir))
		lin := mustAvoid * bot.maxSpeed
		return rot, lin
	}
	return 0, 0
}

// detectAndTakeFood returns true if the bot has found food.
// the food will be removed from the map.
func (bot *Antbot) detectAndTakeFood() bool {
	hasFood := bot.foodpiles.CheckForFoodInCircle(bot.Transform.GetPosition(), bot.botRadius)
	if hasFood {
		return true
	}
	return false
}

func (bot *Antbot) wrapPosition() {

	if bot.Transform.GetPosition().X > float32(rl.GetScreenWidth()) {
		bot.Transform.SetPosition(rl.NewVector2(0, bot.Transform.GetPosition().Y))
	}
	if bot.Transform.GetPosition().X < 0 {
		bot.Transform.SetPosition(rl.NewVector2(float32(rl.GetScreenWidth()), bot.Transform.GetPosition().Y))
	}
	if bot.Transform.GetPosition().Y > float32(rl.GetScreenHeight()) {
		bot.Transform.SetPosition(rl.NewVector2(bot.Transform.GetPosition().X, 0))
	}
	if bot.Transform.GetPosition().Y < 0 {
		bot.Transform.SetPosition(rl.NewVector2(bot.Transform.GetPosition().X, float32(rl.GetScreenHeight())))
	}
}
package code

import (
	"gorl/fw/core/math"
	"gorl/fw/util"
	gomath "math"
	"math/rand"
	"time"

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

	rotVelocity float32
	linVelocity float32

	botRadius float32

	sensorPoints         []rl.Vector2
	pheroSensorDist      float32
	pheroSensorRadius    float32
	obstacleSensorDist   float32
	obstacleSensorRadius float32

	obstacleMap *PheromoneMap
	foodpiles   *FoodPiles

	homeBasePoint  rl.Vector2
	homeBaseRadius float32

	BotMode            mode
	pheromoneIntensity float32
	strictness         float32
	pheromoneTimer     *util.Timer
	pheromoneRefreshed time.Time
	pheromoneLifetime  time.Duration

	// wander
	wanderTheta float32

	// wall avoidance
	avoidanceStrength float32
	avoidanceDir      rl.Vector2
}

func NewAntbot(position rl.Vector2, rotation float32, obstacleMap *PheromoneMap, foodpiles *FoodPiles, homeBasePoint rl.Vector2, homeBaseRadius float32) *Antbot {
	tf := math.NewTransform2D(position, rotation, rl.Vector2One())

	sensorPoints := []rl.Vector2{}
	visionAngle := float32(100.0)
	sensorCount := 3 // don't changes this, we rely on it being 3 by this point
	anglePerSensor := (visionAngle / float32(sensorCount-1)) * rl.Deg2rad
	centerOffset := (visionAngle / 2) * rl.Deg2rad
	for i := 0; i < sensorCount; i++ {
		sensorPoints = append(sensorPoints, rl.Vector2Rotate(rl.NewVector2(0, -1), anglePerSensor*float32(i)-centerOffset))
	}

	return &Antbot{
		Transform: &tf,

		BotMode:            BotModeLeaving,
		strictness:         4.7,
		pheromoneIntensity: 1,
		pheromoneTimer:     util.NewTimer(0.1),
		pheromoneLifetime:  10 * time.Second,

		homeBasePoint:  homeBasePoint,
		homeBaseRadius: homeBaseRadius,

		botRadius: 3,

		maxSpeed: 130,

		rotVelocity: 0,
		linVelocity: 0,

		sensorPoints:         sensorPoints,
		pheroSensorDist:      24,
		pheroSensorRadius:    12,
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
		// if we walk through the home base, we refresh our pheromone supply
		if rl.CheckCollisionCircles(
			bot.Transform.GetPosition(), bot.botRadius,
			bot.homeBasePoint, bot.homeBaseRadius,
		) {
			bot.pheromoneRefreshed = time.Now()
		}

		if bot.pheromoneTimer.Check() {
			toPlace := bot.getPheroIntensity(0.35)
			bot.obstacleMap.SetPheromone(
				math.Vector2IntFromRl(bot.Transform.GetPosition()),
				PheromoneLeaving,
				uint8(toPlace),
			)
		}
		foundFood := bot.detectAndTakeFood()
		if foundFood {
			bot.pheromoneRefreshed = time.Now()
			bot.Transform.AddRotation(180)
			bot.BotMode = BotModeReturning
		}

		wanderStrength := float32(2.0)
		rot, lin = bot.steerWander()
		bot.rotVelocity += rot * wanderStrength
		bot.linVelocity += lin * wanderStrength

		trackStrength := wanderStrength * bot.strictness
		rot, lin = bot.steerTrackPheromone(PheromoneReturning)
		bot.rotVelocity += rot * trackStrength * 1.2
		bot.linVelocity += lin

	case BotModeReturning:
		if bot.pheromoneTimer.Check() {
			toPlace := bot.getPheroIntensity(0.35)
			bot.obstacleMap.SetPheromone(
				math.Vector2IntFromRl(bot.Transform.GetPosition()),
				PheromoneReturning,
				uint8(toPlace),
			)
		}
		wanderStrength := float32(2.0)
		rot, lin = bot.steerWander()
		bot.rotVelocity += rot * wanderStrength
		bot.linVelocity += lin

		trackStrength := wanderStrength * bot.strictness
		rot, lin = bot.steerTrackPheromone(PheromoneLeaving)
		bot.rotVelocity += rot * trackStrength * 1.2
		bot.linVelocity += lin

		if rl.CheckCollisionCircles(bot.Transform.GetPosition(), bot.botRadius, bot.homeBasePoint, bot.homeBaseRadius) {
			bot.Transform.AddRotation(180)
			bot.BotMode = BotModeLeaving
			bot.pheromoneRefreshed = time.Now()
		}

	case BotModeNoFood:
		//bot.botMode = modeLeaving
	}

	foodAndHomeStrength := float32(3)
	rot, lin = bot.steerToFoodAndHome()
	bot.rotVelocity += rot * foodAndHomeStrength
	bot.linVelocity += lin

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
	absPoints := []rl.Vector2{}
	relPoints := []rl.Vector2{}

	found := false
	for idx, sensor := range bot.sensorPoints {
		rotPoint := rl.Vector2Rotate(rl.Vector2Scale(sensor, bot.pheroSensorDist), bot.Transform.GetRotation()*rl.Deg2rad)
		relPoints = append(relPoints, rotPoint)
		absPoint := rl.Vector2Add(
			bot.Transform.GetPosition(),
			rotPoint,
		)

		//rl.DrawCircleV(absPoint, bot.pheroSensorRadius, rl.Green)

		count, agedCount := bot.obstacleMap.HasInCircle(
			math.Vector2IntFromRl(absPoint),
			bot.pheroSensorRadius,
			pheromoneType,
		)

		if count > 0 { // use the integer value to check if there is any pheromone at all
			found = true
			concentrations[idx] = agedCount
		}
		absPoints = append(absPoints, absPoint)
	}

	// use the sensor with the freshest and most hits as the direction to follow
	//const freshWeight = 1
	//const hitWeight = 1

	// use the sensor with the freshest pheromone, prefer middle sensor
	if found {
		// if center > left and center > right
		if concentrations[1] > concentrations[0] && concentrations[1] > concentrations[2] {
			rot, lin = bot.steerSeek(absPoints[1])
		} else if concentrations[0] > concentrations[2] {
			// if left > right
			rot, lin = bot.steerSeek(absPoints[0])
		} else if concentrations[2] > concentrations[0] {
			// if right > left
			rot, lin = bot.steerSeek(absPoints[2])
		}
	}

	// we add the directions together, scaling each one by the concentration.
	// the sum is our direction. see the "VizSensors" scene for a visual representation.
	//if found {
	//	target := rl.Vector2Zero()
	//	rot, lin = float32(0), float32(0)
	//	for idx, point := range relPoints {
	//		point = util.Vector2NormalizeSafe(point)
	//		logging.Debug("%v", concentrations[idx])
	//		c := util.Clamp(concentrations[idx]*100000, 0, 100)
	//		target = rl.Vector2Add(target, rl.Vector2Scale(point, c))
	//	}
	//	absTarget := rl.Vector2Add(bot.Transform.GetPosition(), target)
	//	rl.DrawCircleV(absTarget, 5, rl.Red)
	//	rot, lin = bot.steerSeek(absTarget)
	//}

	return rot, lin
}

// steerToFoodAndHome implements short range steering to food and home points
// within the bot's sensor range.
func (bot *Antbot) steerToFoodAndHome() (float32, float32) {
	rot, _ := float32(0), float32(0)
	validDirections := []rl.Vector2{}
	for _, sensor := range bot.sensorPoints {
		rotPoint := rl.Vector2Rotate(rl.Vector2Scale(sensor, bot.pheroSensorDist), bot.Transform.GetRotation()*rl.Deg2rad)
		absPoint := rl.Vector2Add(
			bot.Transform.GetPosition(),
			rotPoint,
		)
		// Drawing the sensor points
		//rl.DrawCircleV(absPoint, bot.sensorRadius, rl.Red)
		hasFood := bot.foodpiles.CheckForFoodInCircle(
			absPoint,
			bot.pheroSensorRadius,
			false,
		)
		isHome := rl.CheckCollisionCircles(
			absPoint,
			bot.pheroSensorRadius,
			bot.homeBasePoint,
			bot.homeBaseRadius,
		)
		if hasFood || isHome {
			validDirections = append(validDirections, absPoint)
		}
	}
	targetAvg := rl.Vector2Zero()
	for _, dir := range validDirections {
		targetAvg = rl.Vector2Add(targetAvg, dir)
	}
	if len(validDirections) > 0 {
		targetAvg = rl.Vector2Scale(targetAvg, 1/float32(len(validDirections)))
		rot, _ = bot.steerSeek(targetAvg)
	}
	return rot, 0
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
		//rl.DrawCircleV(absPoint, bot.obstacleSensorRadius, rl.Red)

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
	hasFood := bot.foodpiles.CheckForFoodInCircle(bot.Transform.GetPosition(), bot.botRadius, true)
	if hasFood {
		return true
	}
	return false
}

// wrapPosition wraps the bot around the screen if it goes out of bounds.
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

func (bot *Antbot) getPheroIntensity(coef float32) uint8 {
	timeSince := time.Since(bot.pheromoneRefreshed).Seconds()
	intensity := bot.pheromoneIntensity * float32(gomath.Exp(-float64(coef)*timeSince)) * 255
	intensity = util.Clamp(intensity, 0, 255)
	return uint8(intensity)
}

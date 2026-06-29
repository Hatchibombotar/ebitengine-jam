package main

import (
	"fmt"
	"hatchi/disconnect/superui"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

const CONVEYOR_SPEED = 2.0

type Game struct {
	t      int
	player *Character

	uiContext   *superui.UIContext
	hudUI       *superui.UIContainer
	craftingUI  *superui.UIContainer
	todoUI      *superui.UIContainer
	mainMenuUi  *superui.UIContainer
	lossScreen  *superui.UIContainer
	winScreen   *superui.UIContainer
	introScreen *superui.UIContainer

	inventory    [8]*Item
	selectedSlot int

	inCraftingUi bool
	inTodoUI     bool
	inEndScreen  bool

	selectedRecipe int

	currentSpaceId string

	spaces map[string]*Sublevel

	// Relative to grid
	targetX, targetY float64

	selectionName    string
	selectionInRange bool

	day int

	progressBar float64

	timeRemaining int
	lockDownStart int

	endScreenUI *superui.UIContainer

	tasks          []*Task
	startTask      int
	endTask        int
	completedTasks []*Task

	currentTask *Task

	Health    int
	MaxHealth int

	templateItem *Item

	inMainMenu  bool
	inWinScreen bool

	inExclusiveUIMode bool

	dayEndedInDeath bool
	inIntroScreen   bool

	isElectricityDown            bool
	electricityDownRemainingTime int

	isLockDown bool

	testIndex int

	audioContext *audio.Context

	prevLevelDirection string
	prevLevelMoveTime  int
}

func (g *Game) HasItem(s string) bool {
	for _, item := range g.inventory {
		if item == nil {
			continue
		}
		if item.id == s {
			return true
		}
	}
	return false
}

func (g *Game) CurrentSublevel() *Sublevel {
	sublevel, ok := g.spaces[g.currentSpaceId]
	if !ok {
		panic(fmt.Sprint("sublevel doesn't exist.", g.currentSpaceId))
	}
	return sublevel
}

const ITEM_REACH_RANGE = 1.0

func (g *Game) handlePlayerMovement() {
	g.player.startLerpT = 1

	playerCenter := VectorAdd(g.player.position, Vector{.5, .5})

	xSpeed, ySpeed := 0.0, 0.0

	speed := 0.1

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		xSpeed += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		xSpeed -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		ySpeed -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		ySpeed += speed
	}

	xOffset := 0.0
	if xSpeed > 0 {
		xOffset = 4.0 / 16
	} else if xSpeed < 0 {
		xOffset = -4.0 / 16
	}

	yOffset := 0.0
	if ySpeed > 0 {
		yOffset = 0.5
	} else if ySpeed < 0 {
		yOffset = -5.0 / 16
	}

	adjacentVectorX := VectorFloor(VectorAdd(playerCenter, Vector{X: xSpeed + xOffset, Y: 0}))
	if adjacentVectorX.X < 0 {
		if g.CurrentSublevel().adjacentSpaces.West != "" {
			g.player.position.X = 19
			g.currentSpaceId = g.CurrentSublevel().adjacentSpaces.West
			g.prevLevelDirection = "West"
			g.prevLevelMoveTime = g.t
		}
		return
	} else if adjacentVectorX.X > 19 {
		if g.CurrentSublevel().adjacentSpaces.East != "" {
			g.player.position.X = 0
			g.currentSpaceId = g.CurrentSublevel().adjacentSpaces.East
			g.prevLevelDirection = "East"
			g.prevLevelMoveTime = g.t
		}
		return
	}

	adjacentVectorY := VectorFloor(VectorAdd(playerCenter, Vector{X: 0, Y: ySpeed + yOffset}))
	if adjacentVectorY.Y < 0 {
		if g.CurrentSublevel().adjacentSpaces.North != "" {
			g.player.position.Y = 14
			g.currentSpaceId = g.CurrentSublevel().adjacentSpaces.North
			g.prevLevelDirection = "North"
			g.prevLevelMoveTime = g.t
		}
		return
	} else if adjacentVectorY.Y > 14 {
		if g.CurrentSublevel().adjacentSpaces.South != "" {
			g.player.position.Y = 0
			g.currentSpaceId = g.CurrentSublevel().adjacentSpaces.South
			g.prevLevelDirection = "South"
			g.prevLevelMoveTime = g.t
		}
		return
	}

	if TileIsSolid(g, g.CurrentSublevel().tileMap[int(adjacentVectorX.Y)][int(adjacentVectorX.X)]) {
		xSpeed = 0
	}
	if TileIsSolid(g, g.CurrentSublevel().tileMap[int(adjacentVectorY.Y)][int(adjacentVectorY.X)]) {
		ySpeed = 0
	}

	if xSpeed != 0 || ySpeed != 0 {
		g.player.startLerpT = g.t

		if g.t%20 == 0 {
			PlaySound(
				g.audioContext,
				RandomSound(footstepSounds),
				0.3,
			)
		}
	}

	if xSpeed != 0 && ySpeed != 0 {
		xSpeed *= 0.7071
		ySpeed *= 0.7071
	}

	if xSpeed > 0 {
		g.player.facingDirection = Vector{X: 1, Y: 0}
	} else if xSpeed < 0 {
		g.player.facingDirection = Vector{X: -1, Y: 0}
	} else if ySpeed > 0 {
		g.player.facingDirection = Vector{X: 0, Y: 1}
	} else if ySpeed < 0 {
		g.player.facingDirection = Vector{X: 0, Y: -1}
	}

	g.player.position.X += xSpeed
	g.player.position.Y += ySpeed
	if g.player.position.X < 0 && g.CurrentSublevel().adjacentSpaces.West == "" {
		g.player.position.X = 0
	}
}

// func (g *Game) setTileToWall() {
// 	cursorX, cursorY := ebiten.CursorPosition()
// 	targetX, targetY := (cursorX / 16), (cursorY / 16)

// 	if ebiten.IsKeyPressed(ebiten.KeyL) {
// 		g.CurrentSublevel().tileMap[targetY][targetX] = nil
// 	}
// 	if ebiten.IsKeyPressed(ebiten.KeyP) {
// 		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
// 			Type: "wall",
// 		}
// 	}
// 	if ebiten.IsKeyPressed(ebiten.Key0) {
// 		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
// 			Type: "conveyor_left",
// 		}
// 	}
// 	if ebiten.IsKeyPressed(ebiten.Key9) {
// 		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
// 			Type: "conveyor_down",
// 		}
// 	}
// 	if ebiten.IsKeyPressed(ebiten.Key8) {
// 		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
// 			Type: "machine",
// 		}
// 	}
// 	if ebiten.IsKeyPressed(ebiten.KeyV) {
// 		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
// 			Type: "vent_down",
// 		}
// 	}
// 	if ebiten.IsKeyPressed(ebiten.KeyK) {
// 		g.CurrentSublevel().tileMap[targetY][targetX] = &Tile{
// 			Type: "wire",
// 		}
// 	}
// 	if inpututil.IsKeyJustPressed(ebiten.Key7) {
// 		g.CurrentSublevel().conveyorItems = append(g.CurrentSublevel().conveyorItems,
// 			&ConveyorItem{
// 				X: float64(cursorX / 16),
// 				Y: float64(cursorY / 16),
// 				itemType: &Item{
// 					id: "hacking_usb",
// 				},
// 			},
// 		)
// 	}
// }

func (g *Game) StartDay() {
	g.inExclusiveUIMode = false
	g.inEndScreen = false
	g.timeRemaining = 3 * 60 * 60
	g.day += 1
	g.Health = g.MaxHealth

	g.isLockDown = false
	g.isElectricityDown = false

	g.currentSpaceId = "sewer_entrance"

	UpdateTaskRange(g)

	g.completedTasks = []*Task{}

	for _, space := range g.spaces {
		space.Enemies = []*Enemy{}
		if !space.isSafeArea {
			if rand.IntN(2) == 0 {
				continue
			}
			space.Enemies = append(space.Enemies,
				createEnemy(),
			)
		}
	}

	g.spaces["electrical_corridor"].Enemies = append(g.spaces["electrical_corridor"].Enemies, createEnemy(), createEnemy())

}

func PreUpdateConveyors(g *Game) {
	for range 60 * 60 {
		g.t += 1
		for _, space := range g.spaces {

			if space.Spawners != nil {
				for _, spawner := range space.Spawners {
					if g.t%(CONVEYOR_SPEED*16*4) != 0 {
						continue
					}

					space.conveyorItems = append(space.conveyorItems,
						&ConveyorItem{
							X: float64(spawner.X),
							Y: float64(spawner.Y),
							itemType: &Item{
								id: spawner.item,
							},
						},
					)
				}
			}

			UpdateConveyors(g, space)
		}
	}
}
func (g *Game) EndDay() {
	g.inEndScreen = true
	g.inExclusiveUIMode = true
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle(`"Join Us"`)

	g := &Game{}
	g.audioContext = audio.NewContext(SAMPLE_RATE)
	g.Init()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Init() {
	g.t = 0
	g.uiContext = superui.NewUIContext()

	g.inventory = [8]*Item{
		{id: "box"},
	}
	g.currentSpaceId = "component_pnp"
	g.inWinScreen = false
	g.inMainMenu = true
	g.inExclusiveUIMode = true
	g.dayEndedInDeath = false
	g.spaces = createSublevels()

	g.tasks = getTaskList()

	g.Health = 20
	g.MaxHealth = 20

	g.templateItem = &Item{
		id:         "template",
		resultData: "Mind Control",
	}

	g.hudUI = createHudUi(g.uiContext, g)
	g.craftingUI = createCraftingUi(g.uiContext, g)
	g.endScreenUI = CreateEndScreen(g.uiContext, g)
	g.todoUI = createTodoUI(g.uiContext, g)
	g.mainMenuUi = CreateMainMenuUi(g.uiContext, g)
	g.lossScreen = CreateLossScreen(g.uiContext, g)
	g.winScreen = CreateWinScreen(g.uiContext, g)
	g.introScreen = CreateIntroScreen(g.uiContext, g)

	p := &Character{
		position:        Vector{14, 6.5},
		startLerpT:      -1000,
		facingDirection: Vector{-1, 0},
		walkSpeed:       .11,
		speedMultiplier: 1,
	}
	g.player = p

	g.day = 0

	g.prevLevelMoveTime = 0
	g.prevLevelDirection = ""

	PreUpdateConveyors(g)
}

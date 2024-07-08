package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Tile struct {
	x       int
	y       int
	terrain terrain
}

type Entity struct {
	x       float32
	y       float32
	zone    zone
	health  int
	damage  int
	isAlive bool
	money   int
	xp      int
}

const (
	screenWidth  int32 = 800
	screenHeight int32 = 600
	worldWidth         = 20 // in tiles
	worldHeight        = 20
	tileWidth          = 8 // in pixels
	tileHeight         = 8
)

const maxTextures int = 1
const (
	textureFileTilemap = iota
)

type terrain int

const (
	terrainDirt terrain = iota
	terrainGrass
	terrainTree
	terrainStone
)

type zone int

const (
	zoneAll zone = iota
	zoneWorld
	zoneDungeon
)

var textures [maxTextures]rl.Texture2D

const maxSounds int = 5

type soundAsset int

const (
	soundFootGrass soundAsset = iota
	soundFootStone
	soundAttack
	soundDeath
	soundCoins
)

var sounds [maxSounds]rl.Sound

const maxMusics int = 2

type musicAsset int

const (
	musicLightAmbience musicAsset = iota
	musicDarkAmbience
)

var musics [maxMusics]rl.Music

var world [worldWidth][worldHeight]Tile
var dungeon [worldWidth][worldHeight]Tile
var camera rl.Camera2D
var player Entity
var orc Entity
var chest Entity
var dungeonGate Entity

func main() {
	defer gameShutdown()

	gameStartup()

	for !rl.WindowShouldClose() {
		gameUpdate()
		rl.BeginDrawing()
		rl.ClearBackground(rl.Gray)
		gameRender()
		uiRender()
		rl.EndDrawing()
	}
}

func gameStartup() {
	rl.InitWindow(screenWidth, screenHeight, "gocastle main window")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)
	rl.InitAudioDevice()

	image := rl.LoadImage("res/tilemap/colored_tilemap_packed.png")
	textures[textureFileTilemap] = rl.LoadTextureFromImage(image)
	rl.UnloadImage(image)

	for i := 0; i < worldWidth; i++ {
		for j := 0; j < worldHeight; j++ {
			world[i][j] = Tile{
				x:       i,
				y:       j,
				terrain: terrain(rl.GetRandomValue(int32(terrainDirt), int32(terrainTree))),
			}
			dungeon[i][j] = Tile{
				x:       i,
				y:       j,
				terrain: terrainDirt,
			}
		}
	}

	player = Entity{x: tileWidth * 9, y: tileHeight * 9, zone: zoneWorld, health: 100, damage: 0, isAlive: true, money: 1000, xp: 50}
	orc = Entity{x: tileWidth * 5, y: tileHeight * 5, zone: zoneDungeon, health: 100, damage: 0, isAlive: true}
	dungeonGate = Entity{x: tileWidth * 10, y: tileHeight * 10, zone: zoneAll}

	camera = rl.Camera2D{
		Target:   rl.Vector2{X: player.x, Y: player.y},
		Offset:   rl.Vector2{X: float32(screenWidth / 2), Y: float32(screenHeight / 2)},
		Rotation: 0,
		Zoom:     3,
	}

	sounds[soundFootGrass] = rl.LoadSound("res/sounds/grass.wav")
	sounds[soundFootStone] = rl.LoadSound("res/sounds/concrete.wav")
	sounds[soundAttack] = rl.LoadSound("res/sounds/07_human_atk_sword_2.wav")
	sounds[soundDeath] = rl.LoadSound("res/sounds/24_orc_death_spin.wav")
	sounds[soundCoins] = rl.LoadSound("res/sounds/coins.ogg")
	musics[musicLightAmbience] = rl.LoadMusicStream("res/sounds/music_light.wav")
	musics[musicDarkAmbience] = rl.LoadMusicStream("res/sounds/music_dark.wav")

	rl.PlayMusicStream(musics[musicLightAmbience])
}

func gameUpdate() {

	if player.zone == zoneWorld {
		rl.UpdateMusicStream(musics[musicLightAmbience])
	} else {
		rl.UpdateMusicStream(musics[musicDarkAmbience])
	}

	newX := player.x
	newY := player.y
	hasMoved := false

	if rl.IsKeyPressed(rl.KeyA) {
		newX -= tileWidth
		hasMoved = true
	}
	if rl.IsKeyPressed(rl.KeyD) {
		newX += tileWidth
		hasMoved = true
	}
	if rl.IsKeyPressed(rl.KeyW) {
		newY -= tileHeight
		hasMoved = true
	}
	if rl.IsKeyPressed(rl.KeyS) {
		newY += tileHeight
		hasMoved = true
	}

	if hasMoved {
		if player.zone == zoneWorld {
			rl.PlaySound(sounds[soundFootGrass])
		} else {
			rl.PlaySound(sounds[soundFootStone])
		}
	}

	if rl.IsKeyPressed(rl.KeyE) {
		newZone := zoneWorld
		if player.zone == zoneWorld {
			newZone = zoneDungeon
			rl.PlayMusicStream(musics[musicDarkAmbience])
		} else {
			rl.PlayMusicStream(musics[musicLightAmbience])
		}
		if player.x == dungeonGate.x && player.y == dungeonGate.y {
			player.zone = newZone
		}
	}

	if rl.IsKeyPressed(rl.KeyG) {
		if chest.isAlive && player.zone == chest.zone && chest.x == player.x && chest.y == player.y {
			chest.isAlive = false
			player.money += chest.money
			rl.PlaySound(sounds[soundCoins])
		}
	}

	if orc.isAlive && player.zone == orc.zone && orc.x == newX && orc.y == newY {
		damage := rl.GetRandomValue(2, 20)
		orc.health -= int(damage)
		orc.damage = int(damage)
		if orc.health <= 0 {
			orc.isAlive = false
			player.xp += orc.xp
			chest.isAlive = true
			chest.x = orc.x
			chest.y = orc.y
			chest.zone = orc.zone
			chest.money = int(rl.GetRandomValue(10, 100))
			rl.PlaySound(sounds[soundDeath])
		} else {
			rl.PlaySound(sounds[soundAttack])
		}
	} else {
		player.x = newX
		player.y = newY
	}

	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		camera.Zoom += wheel * 0.125
		if camera.Zoom < 3 {
			camera.Zoom = 3
		}
		if camera.Zoom > 8 {
			camera.Zoom = 8
		}
	}
	camera.Target = rl.Vector2{X: player.x, Y: player.y}
}

func uiRender() {
	rl.DrawRectangle(5, 5, 330, 120, rl.Fade(rl.SkyBlue, 0.5))
	rl.DrawRectangleLines(5, 5, 330, 120, rl.Blue)
	rl.DrawText(fmt.Sprintf("zoom: %f", camera.Zoom), 15, 10, 14, rl.Yellow)
	rl.DrawText(fmt.Sprintf("cam target: %f %f", camera.Target.X, camera.Target.Y), 15, 30, 14, rl.Yellow)
	rl.DrawText(fmt.Sprintf("player health: %d", player.health), 15, 50, 14, rl.Yellow)
	rl.DrawText(fmt.Sprintf("player xp: %d", player.xp), 15, 70, 14, rl.Yellow)
	rl.DrawText(fmt.Sprintf("player money: %d", player.money), 15, 90, 14, rl.Yellow)
	rl.DrawText(fmt.Sprintf("orc health: %d", orc.health), 15, 110, 14, rl.Yellow)
}

func gameRender() {
	rl.BeginMode2D(camera)
	defer rl.EndMode2D()

	var (
		tile                         Tile
		textureIndexX, textureIndexY int
	)

	for i := 0; i < worldWidth; i++ {
		for j := 0; j < worldHeight; j++ {
			if player.zone == zoneWorld {
				tile = world[i][j]
			} else {
				tile = dungeon[i][j]
			}
			switch tile.terrain {
			case terrainDirt:
				textureIndexX = 4
				textureIndexY = 4
			case terrainGrass:
				textureIndexX = 5
				textureIndexY = 4
			case terrainTree:
				textureIndexX = 4
				textureIndexY = 5
			}
			drawTile(float32(tile.x*tileWidth), float32(tile.y*tileHeight), textureIndexX, textureIndexY)
		}
	}

	if player.zone == orc.zone && orc.isAlive {
		drawTile(orc.x, orc.y, 11, 0)
	}

	if player.zone == chest.zone && chest.isAlive {
		drawTile(chest.x, chest.y, 9, 3)
	}

	drawTile(dungeonGate.x, dungeonGate.y, 8, 9)
	drawTile(player.x, player.y, 4, 0)

}

func drawTile(xPos float32, yPos float32, textureIndexX int, textureIndexY int) {
	var (
		source, dest rl.Rectangle
		origin       rl.Vector2
	)
	origin = rl.Vector2{X: 0, Y: 0}

	source = rl.Rectangle{
		X:      float32(textureIndexX) * tileWidth,
		Y:      float32(textureIndexY) * tileHeight,
		Width:  tileWidth,
		Height: tileHeight,
	}
	dest = rl.Rectangle{
		X:      xPos,
		Y:      yPos,
		Width:  tileWidth,
		Height: tileHeight,
	}
	rl.DrawTexturePro(textures[textureFileTilemap], source, dest, origin, 0, rl.White)
}

func gameShutdown() {
	for i := 0; i < len(textures); i++ {
		rl.UnloadTexture(textures[i])
	}
	for i := 0; i < len(sounds); i++ {
		rl.UnloadSound(sounds[i])
	}
	for i := 0; i < len(musics); i++ {
		rl.UnloadMusicStream(musics[i])
	}
	rl.CloseAudioDevice()
	rl.CloseWindow()
}

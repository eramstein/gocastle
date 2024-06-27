package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Tile struct {
	x int
	y int
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

var textures [maxTextures]rl.Texture2D

var world [worldWidth][worldHeight]Tile

var camera rl.Camera2D

func main() {
	defer gameShutdown()

	gameStartup()

	for !rl.WindowShouldClose() {
		gameUpdate()
		rl.BeginDrawing()
		rl.ClearBackground(rl.Gray)
		gameRender()
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
				x: i,
				y: j,
			}
		}
	}

	camera = rl.Camera2D{
		Target:   rl.Vector2{X: 0, Y: 0},
		Offset:   rl.Vector2{X: float32(screenWidth / 2), Y: float32(screenHeight / 2)},
		Rotation: 0,
		Zoom:     3,
	}
}

func gameUpdate() {
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
	camera.Target = rl.Vector2{X: 0, Y: 0}
}

func gameRender() {
	rl.BeginMode2D(camera)
	var (
		tile                         Tile
		source, dest                 rl.Rectangle
		origin                       rl.Vector2
		textureIndexX, textureIndexY int
	)
	origin = rl.Vector2{X: 0, Y: 0}
	for i := 0; i < worldWidth; i++ {
		for j := 0; j < worldHeight; j++ {
			tile = world[i][j]
			textureIndexX = 4
			textureIndexY = 4
			source = rl.Rectangle{X: float32(textureIndexX) * tileWidth, Y: float32(textureIndexY) * tileHeight, Width: tileWidth, Height: tileHeight}
			dest = rl.Rectangle{X: float32(tile.x * tileWidth), Y: float32(tile.y * tileHeight), Width: tileWidth, Height: tileHeight}
			rl.DrawTexturePro(textures[textureFileTilemap], source, dest, origin, 0, rl.White)
		}
	}
	rl.EndMode2D()
}

func gameShutdown() {
	for i := 0; i < len(textures); i++ {
		rl.UnloadTexture(textures[i])
	}
	rl.CloseAudioDevice()
	rl.CloseWindow()
}

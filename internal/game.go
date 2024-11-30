package internal

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

const (
	screenWidth  = 640
	screenHeight = 480
	mapWidth     = 8
	mapHeight    = 8
	tileSize     = 64
	fov          = math.Pi / 3
	rayCount     = 120
	moveSpeed    = 3.0
	rotSpeed     = 0.05
	maxDist      = 800
)

var worldMap = [mapWidth][mapHeight]int{
	{1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1},
}

type Player struct {
	x, y, angle float64
}

func (p *Player) moveForward() {
	newX := p.x + math.Cos(p.angle)*moveSpeed
	newY := p.y + math.Sin(p.angle)*moveSpeed

	if !isCollision(newX, newY) {
		p.x = newX
		p.y = newY
	}
}

func (p *Player) moveBackward() {
	newX := p.x - math.Cos(p.angle)*moveSpeed
	newY := p.y - math.Sin(p.angle)*moveSpeed

	if !isCollision(newX, newY) {
		p.x = newX
		p.y = newY
	}
}

func (p *Player) rotate(angle float64) {
	p.angle += angle
	if p.angle < 0 {
		p.angle += 2 * math.Pi
	} else if p.angle > 2*math.Pi {
		p.angle -= 2 * math.Pi
	}
}

func isCollision(x, y float64) bool {
	mapX := int(x / tileSize)
	mapY := int(y / tileSize)

	if mapX < 0 || mapY < 0 || mapX >= mapWidth || mapY >= mapHeight {
		return true
	}
	return worldMap[mapY][mapX] == 1
}

func castRays(renderer *sdl.Renderer, player *Player) {
	for i := 0; i < rayCount; i++ {
		rayAngle := player.angle - (fov / 2) + (fov * float64(i) / rayCount)
		distToWall := raycasting(player.x, player.y, rayAngle)

		// Calculate brightness based on distance
		brightness := 1 - distToWall/maxDist
		if brightness < 0.2 {
			brightness = 0.2 // Minimum brightness
		}

		// Change wall color based on brightness
		color := uint8(255 * brightness)

		// Draw wall
		wallHeight := int32(float64(screenHeight) / distToWall * 200)
		renderer.SetDrawColor(color, color, color, 255) // Wall color based on lighting
		renderer.DrawLine(int32(i*int(screenWidth)/rayCount), (screenHeight-wallHeight)/2,
			int32(i*int(screenWidth)/rayCount), (screenHeight+wallHeight)/2)
	}
}

func raycasting(x, y, angle float64) float64 {
	for distance := 0.0; distance < maxDist; distance += 1 {
		rayX := x + distance*math.Cos(angle)
		rayY := y + distance*math.Sin(angle)
		mapX := int(rayX / tileSize)
		mapY := int(rayY / tileSize)

		if mapX >= 0 && mapX < mapWidth && mapY >= 0 && mapY < mapHeight && worldMap[mapY][mapX] == 1 {
			return distance
		}
	}
	return maxDist // Max distance
}

func RunGame(renderer *sdl.Renderer) error {
	player := &Player{x: 300, y: 300, angle: 0}

	for {
		// Handle SDL events
		event := sdl.PollEvent()
		for event != nil {
			switch event.(type) {
			case *sdl.QuitEvent:
				return nil // Exit the game loop if the window is closed
			}
			event = sdl.PollEvent()
		}

		// Handle keyboard inputs
		keys := sdl.GetKeyboardState()

		if keys[sdl.SCANCODE_W] == 1 {
			player.moveForward()
		}
		if keys[sdl.SCANCODE_S] == 1 {
			player.moveBackward()
		}
		if keys[sdl.SCANCODE_A] == 1 {
			player.rotate(-rotSpeed)
		}
		if keys[sdl.SCANCODE_D] == 1 {
			player.rotate(rotSpeed)
		}

		// Update rendering
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		castRays(renderer, player)

		renderer.Present()
		sdl.Delay(16) // ~60 FPS
	}
}

package main

import (
	"image/color"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth, screenHeigh = 640, 360
	boidCound                = 500
	viewRadius               = 13
	adjRate                  = 0.015
)

var (
	green   = color.RGBA{R: 10, G: 255, B: 50, A: 255}
	boids   [boidCound]*Boid
	boidMap [screenWidth + 1][screenHeigh + 1]int
	rWlock  = sync.RWMutex{}
)

func update(screen *ebiten.Image) error {
	if !ebiten.IsDrawingSkipped() {
		for _, boid := range boids {
			screen.Set(int(boid.position.x+1), int(boid.position.y), green)
			screen.Set(int(boid.position.x-1), int(boid.position.y), green)
			screen.Set(int(boid.position.x), int(boid.position.y+1), green)
			screen.Set(int(boid.position.x), int(boid.position.y-1), green)
		}
	}
	return nil
}

func main() {
	for i, row := range boidMap {
		for j := range row {
			boidMap[i][j] = -1
		}
	}
	for i := 0; i < boidCound; i++ {
		createBoid(i)
	}
	if err := ebiten.Run(update, screenWidth, screenHeigh, 2, "Boids in a box"); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"math"
	"math/rand"
	"time"
)

type Boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func (b *Boid) calcAcceleration() Vector2D {
	upper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)
	avgVelocity := Vector2D{0, 0}
	count := 0.0
	for i := math.Max(lower.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j <= math.Min(upper.y, screenHeigh); j++ {
			// 之前用过下面的代码，会得到 runtime error: index out of range [-1]
			// 原因应该是
			// if boidMap[int(i)][int(j)] != -1 && boidMap[int(i)][int(j)] != b.id && boids[boidMap[int(i)][int(j)]].position.Distance(b.position) < viewRadius {
			// 	count++
			//  avgVelocity = avgVelocity.Add(boids[boidMap[int(i)][int(j)]].velocity) <--- 这里有可能boidMap[int(i)][int(j)]这个点B的线程在A进行计算的时候同时计算了，所以当A走到这一部的时候B已经走了
			// }
			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.id {
				if dist := boids[otherBoidId].position.Distance(b.position); dist < viewRadius {
					count++
					avgVelocity = avgVelocity.Add(boids[otherBoidId].velocity)
				}
			}
		}
	}

	accel := Vector2D{0, 0}
	if count > 0 {
		avgVelocity = avgVelocity.DivisionV(count)
		accel = avgVelocity.Sub(b.velocity).MultiplyV(adjRate)
	}
	return accel
}

func (b *Boid) moveOne() {
	b.velocity = b.velocity.Add(b.calcAcceleration()).limit(-1, 1)
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.Add(b.velocity)
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	next := b.position.Add(b.velocity)
	if next.x >= screenWidth || next.x < 0 {
		b.velocity.x = -b.velocity.x
	}
	if next.y >= screenHeigh || next.y < 0 {
		b.velocity.y = -b.velocity.y
	}

}

func (b *Boid) start() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func createBoid(bid int) {
	b := Boid{
		position: Vector2D{x: rand.Float64() * screenWidth, y: rand.Float64() * screenHeigh},
		velocity: Vector2D{x: (rand.Float64() * 2) - 1.0, y: (rand.Float64() * 2) - 1.0},
		id:       bid,
	}
	boids[bid] = &b
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	go b.start()
}

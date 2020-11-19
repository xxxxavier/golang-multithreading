package main

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type Point2D struct {
	x int
	y int
}

const numberOfThreads int = 4

var (
	r         = regexp.MustCompile(`\((\d*),(\d*)\)`)
	waitGroup = sync.WaitGroup{} // 为什么加waitgroup呢，因为89行关闭channel后，还是会有线程在执行计算
)

func randomPolygon() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	pointsNum := r1.Intn(7) + 3
	pointStr := ""
	for i := 0; i < pointsNum; i++ {
		pointStr += "(" + strconv.Itoa(r1.Intn(50)) + "," + strconv.Itoa(r1.Intn(50)) + ")"
	}
	return pointStr
}

func findArea(inputChannel chan string) {
	for pointsStr := range inputChannel {
		var points []Point2D
		for _, p := range r.FindAllStringSubmatch(pointsStr, -1) {
			x, _ := strconv.Atoi(p[1])
			y, _ := strconv.Atoi(p[2])
			points = append(points, Point2D{x, y})
		}
		area := 0.0
		for i := 0; i < len(points); i++ {
			a, b := points[i], points[(i+1)%len(points)]
			area += float64(a.x*b.y) - float64(a.y*b.x)
		}
		fmt.Println(math.Abs(area) / 2)
	}
	waitGroup.Done()
}

// golang 不支持重载
func findAreaFromString(pointsStr string) {
	var points []Point2D
	for _, p := range r.FindAllStringSubmatch(pointsStr, -1) {
		x, _ := strconv.Atoi(p[1])
		y, _ := strconv.Atoi(p[2])
		points = append(points, Point2D{x, y})
	}
	area := 0.0
	for i := 0; i < len(points); i++ {
		a, b := points[i], points[(i+1)%len(points)]
		area += float64(a.x*b.y) - float64(a.y*b.x)
	}
	fmt.Println(math.Abs(area) / 2)
}

// func main() {
// 	t1 := time.Now()
// 	for i := 0; i < 100000; i++ {
// 		findAreaFromString(randomPolygon())
// 	}
// 	elapsed := time.Since(t1)
// 	fmt.Println("App elapsed: ", elapsed)
// }

func main() {
	t1 := time.Now()
	inputChannel := make(chan string, 1000)
	for i := 0; i < numberOfThreads; i++ {
		waitGroup.Add(1)
		go findArea(inputChannel)
	}
	for i := 0; i < 100000; i++ {
		inputChannel <- randomPolygon()
	}
	close(inputChannel)
	waitGroup.Wait()
	elapsed := time.Since(t1)
	fmt.Println("App elapsed: ", elapsed)
}

package main

import (
	"fmt"
	"time"
)

func main() {
	fibonacci := func() chan uint64 {
		c := make(chan uint64)
		go func() {
			var x, y uint64 = 0, 1
			for ; y < (1 << 63); c <- y { // here
				fmt.Println(x, y)
				x, y = y, x+y // channel如果一直不被消费，就会堵塞，比如现在如果22行的x, ok := <-c不执行，这里就不会下一个循环
			}
			close(c)
		}()
		return c
	}
	c := fibonacci()
	for x, ok := <-c; ok; x, ok = <-c { // here
		time.Sleep(time.Minute)
		fmt.Println(x)
	}
}

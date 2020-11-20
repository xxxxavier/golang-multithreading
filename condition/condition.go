package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	money     = 100
	lock      = sync.Mutex{}
	deposited = sync.NewCond(&lock)
)

func stingy() {
	for i := 1; i <= 1000; i++ {
		lock.Lock()
		money += 10
		deposited.Signal()
		lock.Unlock()
		time.Sleep(1 * time.Millisecond)
	}
	println("Stingy Done")
}

func spendy() {
	for i := 1; i <= 400; i++ {
		lock.Lock()
		for money < 20 {
			deposited.Wait()
		}
		money -= 20
		fmt.Println("run ", i, " times")
		lock.Unlock()
		time.Sleep(1 * time.Millisecond)
	}
	println("Spendy Done")
}

func main() {
	go stingy()
	go spendy()
	time.Sleep(3000 * time.Millisecond)
	print(money)
	// 其实spendy只执行了505次，到最后也没有执行完，只是程序运行到了，进程关闭，自然线程关闭
	// condition在我的理解下，是这样，他会自己解开锁，
	// spendy在发现自己没有足够的值后，通过condition解开锁，有点类似与lock和channel的结合体
	// 在不满足要求的前提下解锁，等待channel的，等channel更新了，再锁
}

package main

import (
	"fmt"
	"time"

	"github.com/colt3k/utils/profile"
)

func main() {
	profile.CpuUsagePercent()
	fmt.Println("start")

	var count = 0
	for {
		profile.MemUsage()
		profile.CpuUsagePercent()
		if count > 10 {
			break
		}
		time.Sleep(1*time.Second)
		count++
	}

	fmt.Println("end")
	profile.MemUsage()
	profile.CpuUsagePercent()
}

package game

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func checkPort(port int) bool {
	portAllocated := fmt.Sprintf(":%d", port)
	// ensure that VPN port is free to allocate
	conn, _ := net.DialTimeout("tcp", portAllocated, time.Second)
	if conn != nil {
		_ = conn.Close()
		fmt.Printf("Checking VPN port %s\n", portAllocated)
		// true means port is already allocated
		return true
	}
	return false
}

func getRandomPort(min, max int) uint {
	port := rand.Intn(max-min) + min
	for checkPort(port) {
		port = rand.Intn(max-min) + min
	}
	return uint(port)
}

func (env *environment) getRandomIp() (string, error) {
	var ip string
	if env.controller.IPPool != nil {
		ipAddress, err := env.controller.IPPool.Get()
		if err != nil {
			return "", err
		}
		ip = ipAddress
	}
	return ip, nil
}

var doThingCounter = 0

func IPcounter() int {
	// Do the thing...
	doThingCounter = doThingCounter + 2
	if doThingCounter >= 240 {
		doThingCounter = 3

	}
	return doThingCounter
}

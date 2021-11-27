package game

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
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

func waitWireguard(ctx context.Context, host, port string) {
	var ready bool
	for ready != true {
		tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		conn, err := grpc.DialContext(tctx, host+port, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			cancel()
			continue
		}
		if conn != nil {
			conn.Close()
			log.Debug().Msg("connected to wireguard VM")
			cancel()
			ready = true
		}
		cancel()
	}
	return
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
	doThingCounter =doThingCounter +2
	if doThingCounter >= 240{
		doThingCounter= 3

	}
	return doThingCounter
}
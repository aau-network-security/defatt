package controller

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
)

var (
	ErrNoAvailableIPs = errors.New("no available IPs")
)

type IPPool struct {
	m       sync.Mutex
	ips     map[string]struct{}
	weights map[string]int
}

func randomPickWeighted(m map[string]int) string {
	var totalWeight int
	for _, w := range m {
		totalWeight += w
	}

	r := rand.Intn(totalWeight)

	for k, w := range m {
		r -= w
		if r <= 0 {
			return k
		}
	}

	return ""
}

func (ipp *IPPool) Get() (string, error) {
	ipp.m.Lock()
	defer ipp.m.Unlock()

	if len(ipp.ips) > 60000 {
		return "", ErrNoAvailableIPs
	}

	genIP := func() string {
		ip := randomPickWeighted(ipp.weights)
		switch ip {
		case "5":
			ip += fmt.Sprintf(".%d", rand.Intn(255))
		case "103":
			ip += fmt.Sprintf(".%d", rand.Intn(255))
		}

		ip += fmt.Sprintf(".%d", rand.Intn(255))

		return ip
	}

	var ip string
	exists := true
	for exists {
		ip = genIP()
		_, exists = ipp.ips[ip]
	}

	ipp.ips[ip] = struct{}{}

	return ip, nil
}

func NewIPPoolFromHost() *IPPool {
	ips := map[string]struct{}{}
	weights := map[string]int{
		"5":   1 * 255,   // 172.{2nd}.{0-255}.{0-255} => 2nd => 25-31 => 6 + 1 => 7
		"103": 255 * 255, // 10.{2nd}.{0-255}.{0-255} => 2nd => 0-254 => 254 + 1 => 255
	}

	ifaces, err := net.Interfaces()
	if err == nil {
		for _, i := range ifaces {
			addrs, err := i.Addrs()
			if err != nil {
				continue
			}

			for _, a := range addrs {
				addr, ok := a.(*net.IPNet)
				//fmt.Printf("addrs: %s\n", addr.String())
				if !ok {
					continue
				}

				if addr.IP.To4() == nil {
					// not v4
					continue
				}

				ipParts := strings.Split(addr.IP.String(), ".")
				lvl1 := ipParts[0]
				if _, ok = weights[lvl1]; !ok {
					// not relevant ip
					continue
				}

				ipStr := strings.Join(ipParts[0:3], ".")
				ips[ipStr] = struct{}{}

				weights[lvl1] = weights[lvl1] - 1
			}
		}
	}

	return &IPPool{
		ips:     ips,
		weights: weights,
	}
}

package utils

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"time"
)

var PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

func IsFoundHost(host string, port uint16) bool {
	target := fmt.Sprintf("%s:%d", host, port)

	_, err := net.DialTimeout("tcp", target, time.Second)
	fmt.Println("Dialing", target, "...")
	return err == nil
}

func FindNeighbors(currHost string, currPort uint16, startIP, endIP uint8, startPort, endPort uint16) []string {
	currAddress := fmt.Sprintf("%s:%d", currHost, currPort)

	m := PATTERN.FindStringSubmatch(currHost)
	if len(m) == 0 {
		println("Invalid host name")
		return nil
	}

	prefix := m[1]

	lastIP, _ := strconv.Atoi(m[len(m)-1])

	neighbors := make([]string, 0)

	for port := startPort; port <= endPort; port++ {
		for ip := startIP; ip <= endIP; ip++ {
			host := fmt.Sprintf("%s%d", prefix, lastIP+int(ip))
			address := fmt.Sprintf("%s:%d", host, port)
			if address != currAddress && IsFoundHost(host, port) {
				neighbors = append(neighbors, address)
			}
		}
	}

	return neighbors
}

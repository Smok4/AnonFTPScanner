//Made By Namz
// For educational prupose only

package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
)

func main() {
	ipRanges := []string{
		"185.0.0.0/8",            //put the range you want scan
		"188.0.0.0/8",
		"195.0.0.0/8",
		"212.0.0.0/8",
	}

	resultFile, err := os.Create("results.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer resultFile.Close()

	for _, ipRange := range ipRanges {
		_, network, err := net.ParseCIDR(ipRange)
		if err != nil {
			log.Fatal(err)
		}

		for ip := network.IP; network.Contains(ip); incIP(ip) {
			go func(ip net.IP) {
				if connectAnonymousFTP(ip, resultFile) {
					result := fmt.Sprintf("Connected to: %s\n", ip)
					fmt.Print(result) // Print to console
				}
			}(ip)
		}
	}

	// Wait for the connections to finish
	time.Sleep(5 * time.Second)
}

func connectAnonymousFTP(ip net.IP, resultFile *os.File) bool {
	// Simulating a connection attempt
	fmt.Printf("Attempting FTP connection to: %s\n", ip)

	// Simulating network latency
	time.Sleep(time.Duration(randInt(500, 2000)) * time.Millisecond)

	// Simulating a successful FTP connection for demonstration purposes
	success := randInt(0, 2) == 0
	if success {
		fmt.Printf("Successfully connected to FTP at: %s\n", ip)
		resultFile.WriteString(fmt.Sprintf("Connected to FTP at: %s\n", ip)) // Write to file
	} else {
		fmt.Printf("Failed to connect to FTP at: %s\n", ip)
	}

	return success
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func randInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

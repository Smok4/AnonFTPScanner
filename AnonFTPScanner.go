//Made By Namz
// For educational purpose only

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var mutex sync.Mutex // Mutex to protect IP scanning

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./monscript <iprange>")
		return
	}

	ipRange := os.Args[1]

	_, network, err := net.ParseCIDR(ipRange)
	if err != nil {
		log.Fatal(err)
	}

	resultFile, err := os.Create("results.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer resultFile.Close()

	wg := sync.WaitGroup{}

	for ip := network.IP; network.Contains(ip); incIP(ip) {
		wg.Add(1)
		go func(ip net.IP) {
			defer wg.Done()

			if connectAnonymousFTP(ip) {
				result := fmt.Sprintf("Connected to: %s\n", ip)
				fmt.Print(result) // Print to console

				mutex.Lock()
				if _, err := resultFile.WriteString(result); err != nil {
					log.Printf("Error writing to file: %s", err)
				}
				mutex.Unlock()
			}
		}(ip)
	}

	wg.Wait()
}

func connectAnonymousFTP(ip net.IP) bool {
	fmt.Printf("Attempting FTP connection to: %s\n", ip)

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:21", ip), 5)
	if err != nil {
		fmt.Printf("Failed to connect to FTP at: %s\n", ip)
		return false
	}
	defer conn.Close()

	// Read the server response
	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		fmt.Printf("Failed to read server response: %s\n", err)
		return false
	}

	// Send the USER command
	userCmd := fmt.Sprintf("USER %s\r\n", "anonymous")
	_, err = conn.Write([]byte(userCmd))
	if err != nil {
		fmt.Printf("Failed to send USER command: %s\n", err)
		return false
	}

	// Read the response to the USER command
	_, err = conn.Read(response)
	if err != nil {
		fmt.Printf("Failed to read server response: %s\n", err)
		return false
	}

	// Check if the response code indicates success (code 230)
	if string(response[0:3]) != "230" {
		fmt.Printf("Failed to log in anonymously to FTP at: %s\n", ip)
		return false
	}

	fmt.Printf("Successfully connected to FTP at: %s\n", ip)
	return true
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

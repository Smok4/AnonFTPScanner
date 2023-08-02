// Made By Namz
// For educational purpose only
// Update 02/08/2023
// apt-get install golang
// go build AnonFTPScanner
// Usage ./AnonFTPScanner youriplist.txt
// remmoved github.com/jlaffaye/ftp

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var mutex sync.Mutex // Mutex to protect IP scanning

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./monscript <filename>")
		return
	}

	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer file.Close()

	resultFile, err := os.Create("results.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer resultFile.Close()

	// Create a buffered writer to reduce the number of system calls
	writer := bufio.NewWriter(resultFile)
	defer writer.Flush()

	wg := sync.WaitGroup{}
	successChannel := make(chan string)

	// Create a worker pool with 10 workers (adjust as needed)
	workerCount := 1000
	workerPool := make(chan struct{}, workerCount)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ipStr := scanner.Text()
		ip := net.ParseIP(ipStr)
		if ip == nil {
			fmt.Printf("Invalid IP address: %s\n", ipStr)
			continue
		}

		wg.Add(1)
		workerPool <- struct{}{} // Acquire a worker slot
		go func(ip net.IP) {
			defer wg.Done()
			defer func() { <-workerPool }() // Release the worker slot

			if connectAnonymousFTP(ip) {
				result := fmt.Sprintf("Connected to: %s\n", ip)
				fmt.Print(result) // Print to console

				// Lock the mutex before writing to the file
				mutex.Lock()
				if _, err := writer.WriteString(result); err != nil {
					log.Printf("Error writing to file: %s", err)
				}
				mutex.Unlock()
				// Unlock the mutex after writing

				successChannel <- result
			}
		}(ip)
	}

	go func() {
		for range successChannel {
			// No need to process the result here since we are writing it to the file in the goroutine
		}
	}()

	wg.Wait()
	close(successChannel)
}

func connectAnonymousFTP(ip net.IP) bool {
	fmt.Printf("Attempting FTP connection to: %s\n", ip)

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:21", ip), 5*time.Second)
	if err != nil {
		fmt.Printf("Failed to connect to FTP at: %s\n", ip)
		return false
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Read the server response
	response, err := reader.ReadString('\n')
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
	response, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Failed to read server response: %s\n", err)
		return false
	}

	// Check if the response code indicates success (code 230)
	if response[0] != '2' {
		fmt.Printf("Failed to log in anonymously to FTP at: %s\n", ip)
		return false
	}

	// Send the PASS command
	passCmd := fmt.Sprintf("PASS %s\r\n", "anonymous")
	_, err = conn.Write([]byte(passCmd))
	if err != nil {
		fmt.Printf("Failed to send PASS command: %s\n", err)
		return false
	}

	// Read the response to the PASS command
	response, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Failed to read server response: %s\n", err)
		return false
	}

	// Check if the response code indicates success (code 230)
	if response[0] != '2' {
		fmt.Printf("Failed to log in anonymously to FTP at: %s\n", ip)
		return false
	}

	fmt.Printf("Successfully connected to FTP at: %s\n", ip)
	return true
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	conn, err := net.Dial("tcp", ":7007")
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	var wg sync.WaitGroup

	var name string
	fmt.Print("Name? ")
	fmt.Scan(&name)

	log.Printf("connected to server %s\n", conn.RemoteAddr())

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			var msg string
			// fmt.Printf("%s: ", name)
			reader := bufio.NewReader(os.Stdin)
			msg, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("error reading input: %v\n", err)
				return
			}

			msg = strings.TrimSpace(msg)

			if msg == "exit" || msg == "quit" {
				log.Println("connection closed by client")
				conn.Close()
				return
			}

			_, err = conn.Write([]byte(fmt.Sprintf("%s: %s\n", name, msg)))
			if err != nil {
				log.Printf("error sending msg %v\n", err)
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		wg.Done()

		for {
			buffer := make([]byte, 4096)
			_, err = conn.Read(buffer)
			if err != nil {
				log.Printf("error reading msg %v\n", err)
				return
			}

			fmt.Printf("%s", buffer)
		}
	}()

	wg.Wait()
	// select {}
}

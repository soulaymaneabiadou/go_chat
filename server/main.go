package main

import (
	"io"
	"log"
	"net"
	"sync"
)

var conns = &sync.Map{}

func main() {
	ln, err := net.Listen("tcp", ":7007")
	if err != nil {
		log.Panicf("error listening %v", err)
	}

	log.Println("listening for tcp on port 7007")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("error accepting conn %v", err)
		}

		conns.Store(conn.RemoteAddr(), conn)

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		conns.Delete(conn.RemoteAddr())
	}()

	for {
		buffer := make([]byte, 4096)
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Printf("connection from %s closed\n", conn.RemoteAddr())
				return
			}

			log.Printf("error reading from %s: %v\n", conn.RemoteAddr(), err)
			return
		}

		//
		conns.Range(func(key, value any) bool {
			c := value.(net.Conn)

			if key != conn.RemoteAddr() {
				_, err = c.Write(buffer[:n])
				if err != nil {
					log.Printf("error broadcasting message from %s: %v\n", conn.RemoteAddr(), err)
				}
			}

			return true
		})

		// _, err = conn.Write(buffer[:n])
		// if err != nil {
		// 	log.Printf("error writing to %s: %v\n", conn.RemoteAddr(), err)
		// 	return
		// }
	}
}

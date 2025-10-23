package main

import (
	"fmt"
	"net"
	"time"
)

func StartServer(addr string) {
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		fmt.Printf("Server connection err: %v\n", err)
		return
	}

	for {

		conn, err := listener.Accept()

		if err != nil {
			fmt.Printf("Error with accepting connections: %v\n", err)
			return
		}

		buf := make([]byte, 1024)
		done := make(chan struct{})

		go func() {
			defer conn.Close()

			for {
				//extra for loop is to terminate connection after 3 success round of msg exchange
				for i := 0; i < 3; i++ {
					if err = conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
						fmt.Printf("SetReadDeadline error: %v\n", err)
						return
					}

					n, err := conn.Read(buf)

					if err != nil {
						fmt.Printf("Cannot read: %q\n", err)
						return
					}

					fmt.Printf("Client sent: %q\n", buf[:n])

					if err = conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
						fmt.Printf("SetReadDeadline error: %v\n", err)
						return
					}

					_, err = conn.Write([]byte("I have Listened.\n"))

					if err != nil {
						fmt.Printf("Client write failed: %v\n", err)
						return
					}

				}

				<-done //just blocks it so the connection can terminate
			}

		}()

	}
}

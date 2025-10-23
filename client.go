package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

func StartClient(addr string) {
	ctxClient, cancelClient := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancelClient()

	d := net.Dialer{}

	conn, err := d.DialContext(ctxClient, "tcp", addr)

	if err != nil {
		fmt.Printf("Err with client connection %q\n", err)
		return
	}

	ctxPinger, cancelPinger := context.WithCancel(context.Background())

	resetTimer := make(chan time.Duration, 1)

	resetTimer <- 2 * time.Second

	go Pinger(ctxPinger, conn, resetTimer)

	buf := make([]byte, 1024)

	go func() {
		defer func() {
			close(resetTimer)
			cancelPinger()
			cancelClient()
			conn.Close()
		}()

		for {

			if err = conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
				fmt.Printf("SetDeadline error: %v\n", err)
				return
			}

			n, err := conn.Read(buf)

			if err != nil {
				fmt.Printf("Error reading from server %q\n", err)
				return
			}

			resetTimer <- 0

			fmt.Printf("Data processed and recieved is: %v\n", buf[:n])

			if err = conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
				fmt.Printf("SetDeadline error: %v\n", err)
				return
			}

			_, err = conn.Write([]byte("This is a message from client\n"))

			if err != nil {
				fmt.Printf("Error reading from server %q\n", err)
				return
			}

		}

	}()

}

func Pinger(ctx context.Context, conn net.Conn, resetTimer <-chan time.Duration) {
	var interval time.Duration
	var err error
	select {
	case <-ctx.Done():
		return
	case interval = <-resetTimer:
	default:
	}

	if interval <= 0 {
		interval = 30 * time.Second
	}

	timer := time.NewTimer(interval)

	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	for {

		select {

		case <-ctx.Done():
			return

		case newTimer := <-resetTimer:

			if !timer.Stop() {
				<-timer.C
			}

			if newTimer > 0 {
				interval = newTimer
			}

		case <-timer.C:

			if err = conn.SetWriteDeadline(time.Now().Add(2 * time.Second)); err != nil {
				fmt.Printf("SetDeadline error: %v\n", err)
				return
			}

			_, err := conn.Write([]byte("ping"))

			if err != nil {
				fmt.Printf("failed to write: %v\n", err)
				return
			}
		}

		timer.Reset(interval)

	}
}

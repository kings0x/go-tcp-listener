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
		fmt.Printf("Err with client connection %q", err)
		return
	}

	defer conn.Close()

	ctxPinger, cancelPinger := context.WithCancel(context.Background())
	defer cancelPinger()

	resetTimer := make(chan time.Duration)

	go Pinger(ctxPinger, conn, resetTimer)

	if err = conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		fmt.Printf("SetDeadline error: %v\n", err)
		return
	}

	buf := make([]byte, 1024)

	go func() {

		for {
			n, err := conn.Read(buf)

			if err != nil {
				fmt.Printf("Error reading from server %q\n", err)
				return
			}

			resetTimer <- 0

			fmt.Printf("Data processed and recieved is: %v", buf[:n])

			if err = conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
				fmt.Printf("SetDeadline error: %v\n", err)
				return
			}

			_, err = conn.Write([]byte("This is a message from client"))

			if err != nil {
				fmt.Printf("Error reading from server %q", err)
				return
			}

			resetTimer <- 0
		}

	}()

}

func Pinger(ctx context.Context, conn net.Conn, resetTimer <-chan time.Duration) {
	var interval time.Duration
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
			_, err := conn.Write([]byte("ping"))

			if err != nil {
				fmt.Printf("failed to write: %v", err)
				return
			}
		}

		timer.Reset(interval)

	}
}

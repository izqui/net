package main

import (
	"fmt"
	"github.com/lab11/go-tuntap/tuntap"
	"io"
	"os/exec"
	"time"
)

func main() {

	tund, err := tuntap.Open("tap0", tuntap.DevTap, false)
	panicOnError(err)

	cmd := exec.Command("ifconfig", "tap0", "inet6", "beef::1/10", "up")
	out, err := cmd.Output()
	panicOnError(err)

	fmt.Println(string(out))

	time.Sleep(5 * time.Second)

	for {

		input, err := tund.ReadPacket()
		panicOnError(err)

		fmt.Println(string(input.Packet))

	}

}

func panicOnError(err error) {
	if err != nil && err != io.EOF {
		panic(err)
	}
}

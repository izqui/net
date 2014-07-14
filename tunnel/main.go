package main

import (
	"fmt"
	"github.com/lab11/go-tuntap/tuntap"
	"io"
	"os/exec"
	"time"
)

func main() {

	tund, err := tuntap.Open("tun0", tuntap.DevTun, false)
	panicOnError(err)

	cmd := exec.Command("ifconfig", "tun0", "inet6", "fe80::1", "up")
	out, err := cmd.Output()
	panicOnError(err)

	cmd = exec.Command("route", "-n", "add", "-inet6", "beef::/10", "fe80::1")
	out, err = cmd.Output()
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

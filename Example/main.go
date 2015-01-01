package main

import (
	"log"
	"syscall"
	"time"

	"github.com/tommessick/netSerial"
)

// Open a serial port at 1200 BPS and map it to localhost:8765
// Quit running after about 8 hours
func main() {
	err := netSerial.Open("/dev/ttyUSB0", syscall.B1200, ":8765")
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(30000 * time.Second)
}

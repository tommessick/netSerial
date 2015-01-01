// Package netserial allows accessing a serial port from a network connection
package netSerial

import (
	"fmt"
	"log"
	"net"
	"syscall"

	"github.com/schleibinger/sio"
)

// Open initializes the serial port and listens on the TCP port
func Open(dev string, rate uint32, port string) (err error) {

	var fromSerial = make(chan byte)

	p, err := sio.Open("/dev/fluke", syscall.B9600)
	if err != nil {
		return err
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		return err
	}

	fmt.Println(tcpAddr)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	fmt.Println(listener)

	// Read serial port and write to channel
	go func() {
		rb := make([]byte, 128)
		for {
			count, err := p.Read(rb)
			if err != nil {
				log.Fatal(err)
			}
			for i := 0; i < count; i++ {
				fromSerial <- rb[i]
			}
		}

	}()

	// Listen for network connections
	go func() {
		var count = 0
		for {
			conn, err := listener.Accept()
			if count == 0 {
				count++
				if err != nil {
					log.Fatal(err)
				}
				// quit chan is used to kill the serial goroutine
				//    when the network goroutine exits
				quit := make(chan struct{})

				go readNet(conn, *p, quit)
				go readSerial(fromSerial, conn, *p, quit)
				go func() {
					<-quit
					count--
				}()
			} else {
				conn.Write([]byte("Connection already in use\n"))
				conn.Close()

			}
		}
	}()

	return err
}

func readNet(conn net.Conn, p sio.Port, quit chan struct{}) {
	defer conn.Close()
	defer close(quit)

	rb := make([]byte, 1)
	for {
		_, err := conn.Read(rb[0:])
		if err != nil {
			return
		}
		_, err = p.Write(rb)
		if err != nil {
			return
		}
	}
}

func readSerial(fromSerial chan byte, conn net.Conn, p sio.Port, quit chan struct{}) {
	for {
		wb := make([]byte, 1)
		var x byte

		select {
		case x = <-fromSerial:
			wb[0] = x
			_, err := conn.Write(wb)
			if err != nil {
				return
			}
		case <-quit:
			return
		}
	}
}

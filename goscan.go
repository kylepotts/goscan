package goscan

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"math"
	"net"
	"runtime"
	"time"
)

type scanner struct {
	host string
}

type tcpResponse struct {
	port int
	name string
}

//    Create a new scanner object
func NewScanner(host string) *scanner {
	runtime.GOMAXPROCS(runtime.NumCPU())
	return &scanner{host}
}

//    String representation of the scanner object
func (s scanner) String() string {
	return fmt.Sprintf("Scanner - Host:%s", s.host)
}

//    print format of "192.168.1.1:8080"
func (s scanner) hostAndPort(port int) string {
	return fmt.Sprintf("%s:%d", s.host, port)
}

// Check if current port is open
// sends a tcp ping to the host and port and if there is a response returns true
func (s scanner) IsPortOpen(port int) bool {
	durration, err := time.ParseDuration("0.5s")
	if err != nil {
		return false
	}
	conn, err := net.DialTimeout("tcp4", s.hostAndPort(port), durration)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// Go Routine that go from start port to end port, add it to a list if it is open
func (s scanner) getPorts(start int, end int, c chan []int) {
	//ports:= <- c
	ports := []int{}
	for i := start; i < end; i++ {
		if s.IsPortOpen(i) {
			ct.ChangeColor(ct.Red, false, ct.None, false)
			//fmt.Printf("Rountine[%d-%d], port[%d] success\n", start, length, i)
			ct.ResetColor()
			ports = append(ports, i)
		} else {
			//fmt.Printf("Rountine[%d-%d], port[%d] failure\n", start, length, i)

		}

	}
	c <- ports
	return

}

// Get all the open ports from start to end, splitting each thread to have (start-end)/(length) ports to work with
func (s scanner) GetOpenPortsRoutine(start int, end int, split int) []int {
	//n is the number of ports each routines will handle
	n := int(math.Floor(float64((end - start) / (split))))
	numChannels := 0
	openPorts := []int{}
	k := 0
	// find how many channels we need to make for each rountines
	for j := 0; j < end; j += n {
		numChannels++
	}
	// create and initalize each channel
	channels := make([]chan []int, numChannels)

	for j := 0; j < numChannels; j++ {
		channels[j] = make(chan []int)
	}

	// main loop to start each rountines
	for i := start; i < end; i += n {
		if i+n >= end {
			go s.getPorts(i, end, channels[k])
			break
		}
		go s.getPorts(i, i+n, channels[k])
		k++
	}

	// loop through each  channels and get the open ports
	for k := 0; k < numChannels; k++ {
		p := <-channels[k]
		for _, port := range p {
			openPorts = append(openPorts, port)
		}
		close(channels[k])
	}
	return openPorts
}

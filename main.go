package main

import (
	"fmt"
	"math"
	"net"
	"runtime"
	"time"
    "github.com/daviddengcn/go-colortext"
)

type scanner struct {
	host string
}

type tcpResponse struct {
    port int
    name string
}

func main() {
	s := newScanner("kptechblog.com")
	fmt.Println(s)
	runtime.GOMAXPROCS(runtime.NumCPU())
	s.getOpenPortsRoutine(15, 120, 4)
}

func newScanner(host string) *scanner {
	return &scanner{host}
}

func (s scanner) String() string {
	return fmt.Sprintf("Scanner - Host:%s", s.host)
}

func (s scanner) hostAndPort(port int) string {
	return fmt.Sprintf("%s:%d", s.host, port)
}

func (s scanner) IsPortOpen(port int) bool {
	durration, err := time.ParseDuration("0.5s")
	if err != nil {
		fmt.Println("error")
		return false
	}
	conn, err := net.DialTimeout("tcp4", s.hostAndPort(port), durration)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func (s scanner) getPorts(start int, length int,c chan []int) {
	//ports:= <- c
    ports:= []int{}
	for i := start; i < length; i++ {
		if s.IsPortOpen(i) {
            ct.ChangeColor(ct.Red,false,ct.None,false)
			fmt.Printf("Rountine[%d-%d], port[%d] success\n", start, length, i)
            ct.ResetColor()
			ports = append(ports, i)
		} else {
			fmt.Printf("Rountine[%d-%d], port[%d] failure\n", start, length, i)

		}

	}
	c <- ports
	fmt.Println("Done")
	return

}

func (s scanner) getOpenPortsRoutine(start int, end int, split int) {
    //n is the number of ports each routines will handle
	n :=int( math.Floor(float64((end - start) / (split))))
    numChannels:=0
    openPorts:= []int{}
    k:=0
    // find how many channels we need to make for each rountines
	for j := 0; j < end; j += n {
		numChannels++
	}
    // create and initalize each channel
    channels := make([]chan []int,numChannels)

    for j:=0; j<numChannels; j++{
        channels[j] = make(chan []int)
    }

    // main loop to start each rountines
	for i := start; i < end; i += n {
		if i+n >= end {
			go s.getPorts(i, end,channels[k])
			break
		}
		go s.getPorts(i, i+n,channels[k])
        k++
	}

    // loop through each  channels and get the open ports
    for k:=0; k<numChannels; k++{
        p:= <-channels[k]
        for _, port:=range p{
            openPorts = append(openPorts,port)
        }
        close(channels[k])
    }
    fmt.Println(openPorts)
}



package main

import (
	"fmt"
	"math"
	"net"
	"runtime"
	"time"
)

type scanner struct {
	host string
}

func main() {
	s := newScanner("kyle-potts.com")
	fmt.Println(s)
	runtime.GOMAXPROCS(2)
	s.getOpenPortsRoutine(15, 90, 4)
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

func (s scanner) getPorts(start int, length int) {
	//ports:= <- c
	ports := []int{}
	for i := start; i < length; i++ {
		if s.IsPortOpen(i) {
			fmt.Printf("Rountine[%d-%d], port[%d] success\n", start, length, i)
			ports = append(ports, i)
		} else {
			fmt.Printf("Rountine[%d-%d], port[%d] failure\n", start, length, i)

		}

	}
	//c <- ports
	fmt.Println("Done")
	return

}

func (s scanner) getOpenPortsRoutine(start int, end int, split int) {
	numRountines := math.Floor(float64((end - start) / (split)))
	n := int(numRountines)
	var numChannels = 0
	for j := 0; j < end; j += n {
		numChannels++
	}

	fmt.Printf("numChannels=%d\n", numChannels)
	//channels := make([]chan []int,numChannels)

	fmt.Printf("NumRoutines = %d\n", int(numRountines+0.5))
	for i := start; i < end; i += n {
		if i+n >= end {
			fmt.Printf("%d %d\n", i, end)
			go s.getPorts(i, end)
			break
		}
		fmt.Printf("%d %d\n", i, i+n)
		go s.getPorts(i, i+n)
	}
	//ports:=[]int{}
	//c<-ports
	//p:=<-c
	//fmt.Println(p)
	//close(c)

}

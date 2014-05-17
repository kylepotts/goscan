package main
import "fmt"
import "net"
import "time"
import "math"
import "sync"

type scanner struct {
    host string
}

func main(){
    s := newScanner("kyle-potts.com")
    fmt.Println(s)
    p := make(chan []int)
    var wg sync.WaitGroup
    s.getOpenPortsRoutine(15,20,2,p,wg)
}

func newScanner(host string) *scanner {
    return &scanner{host}
}

func (s scanner)String() string {
    return fmt.Sprintf("Scanner - Host:%s",s.host)
}

func (s scanner) hostAndPort(port int) string {
    return fmt.Sprintf("%s:%d", s.host,port)
}

func (s scanner)IsPortOpen(port int) bool {
    durration,err := time.ParseDuration("2s")
    if(err != nil){
        fmt.Println("error")
        return false
    }
    fmt.Println(durration)
    conn,err := net.DialTimeout("tcp4",s.hostAndPort(port),durration)
    if(err != nil){
        return false
    }
    defer conn.Close()
    return true
}

func (s scanner) getPorts(start int, length int, c chan []int,wg sync.WaitGroup){
    defer wg.Done()
    fmt.Printf("Starting getPorts from %d to %d\n",start, start+length)
    ports := <- c
    for i:=start; i<start+length; i++{
        if(s.IsPortOpen(i)){
            ports = append(ports,i)
        }

    }
    c <- ports
    fmt.Println("Done")


}

func (s scanner)getOpenPortsRoutine(start int, end int,split int,c chan []int, wg sync.WaitGroup){
    numRountines := math.Ceil(float64((end-start)/(split)))
    fmt.Printf("numRountines = %lf\n")
    for i:=start; i<end; i+=int(numRountines){
            wg.Add(1)
        go s.getPorts(i,i+int(numRountines),c,wg)
    }
    ports := []int{}
    c<-ports
    wg.Wait()
    p:=<-c
    fmt.Println(p)
    close(c)



}

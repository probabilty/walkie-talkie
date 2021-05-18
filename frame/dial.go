package frame

import (
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	// "time"
	"walki-talki/phonebook"
)

var activeCalls map[string]string

// var connMutex sync.Mutex
var frameMutex sync.Mutex

func Init(conn *net.UDPConn) {
	activeCalls = make(map[string]string)
	// go func() {
	// 	for {
	// 		pingClients(conn)
	// 	}
	// }()
}
func Dial(conn *net.UDPConn, addr string, channel string) {
	//check if there is a call currently
	if activeCalls[addr] != channel {
		Hungup(conn, addr)
	}
	frameMutex.Lock()
	activeCalls[addr] = channel
	frameMutex.Unlock()
}
func IsInACall(addr string) bool {
	if activeCalls[addr] != "" {
		return true
	}
	return false
}
func Hungup(conn *net.UDPConn, addr string) {
	deadline := time.Now().Add(time.Second)
	conn.SetWriteDeadline(deadline)
	conn.WriteToUDP([]byte("ok\n"), getAddress(addr))
	frameMutex.Lock()
	delete(activeCalls, addr)
	frameMutex.Unlock()
}
func Relay(conn *net.UDPConn, addr string, data []byte) {
	if string(data) == "Hangup" {
		Hungup(conn, addr)
		return
	}
	callChannel := activeCalls[addr]
	reciptors := phonebook.Get(callChannel)

	for i := 0; i < len(reciptors); i++ {
		go func(i int) {
			if len(data) == 0 {
				return
			}
			// connMutex.Lock()
			deadline := time.Now().Add(time.Second)
			conn.SetWriteDeadline(deadline)
			_, err := conn.WriteTo(data, reciptors[i])
			// connMutex.Unlock()
			if err != nil {
				log.Println(err)
				Hungup(conn, addr)
				return
			}
		}(i)
	}

}
func SendOK(conn *net.UDPConn, addr *net.UDPAddr) {
	// connMutex.Lock()
	deadline := time.Now().Add(time.Second)
	conn.SetWriteDeadline(deadline)
	conn.WriteToUDP([]byte("ok\n"), addr)
	// connMutex.Unlock()
}
func pingClients(conn *net.UDPConn) {
	for k := range activeCalls {
		parts := strings.SplitAfter(k, ":")
		port := parts[1]
		portNum, _ := strconv.Atoi(port)
		parts[0] = strings.TrimSuffix(parts[0], ":")
		// ip := net.ParseIP(parts[0])
		ip := net.ParseIP("10.15.14.30")
		// zone := net.parseIPZone(parts[0])
		deadline := time.Now().Add(time.Second)
		conn.SetWriteDeadline(deadline)
		n, err := conn.WriteToUDP([]byte("ok\n"), &net.UDPAddr{
			IP:   ip,
			Port: portNum,
			Zone: "",
		})
		// connMutex.Unlock()
		if err != nil || n == 0 {
			frameMutex.Lock()
			delete(activeCalls, k)
			frameMutex.Unlock()
		}
	}
	log.Printf("%+v\n", activeCalls)
	time.Sleep(5 * time.Second)
}
func getAddress(add string) (addr *net.UDPAddr) {
	parts := strings.SplitAfter(add, ":")
	port := parts[1]
	portNum, _ := strconv.Atoi(port)
	parts[0] = strings.TrimSuffix(parts[0], ":")
	ip := net.ParseIP(parts[0])
	addr = &net.UDPAddr{
		IP:   ip,
		Port: portNum,
		Zone: "",
	}
	return

}

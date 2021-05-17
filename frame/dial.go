package frame

import (
	"net"
	"sync"

	// "time"
	"walki-talki/phonebook"
)

var activeCalls map[string]string
var connMutex sync.Mutex
var frameMutex sync.Mutex

func Init() {
	activeCalls = make(map[string]string)
}
func Dial(addr string, channel string) {
	//check if there is a call currently
	if activeCalls[addr] != channel {
		Hungup(addr)
	}
	activeCalls[addr] = channel
}
func IsInACall(addr string) bool {
	if activeCalls[addr] != "" {
		return true
	}
	return false
}
func Hungup(addr string) {
	frameMutex.Lock()
	delete(activeCalls, addr)
	frameMutex.Unlock()
}
func Relay(conn *net.UDPConn, addr string, data []byte) {
	if string(data) == "Hangup" {
		Hungup(addr)
		return
	}
	callChannel := activeCalls[addr]
	reciptors := phonebook.Get(callChannel)

	for i := 0; i < len(reciptors); i++ {
		go func(i int) {
			connMutex.Lock()
			// deadline := time.Now().Add(time.Second)
			// conn.SetReadDeadline(deadline)
			_, err := conn.WriteToUDP(data, reciptors[i])
			connMutex.Unlock()
			if err != nil {
				Hungup(addr)
				return
			}
		}(i)
	}

}

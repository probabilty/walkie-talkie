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

var activeCalls map[string]call

const MAX_CALL_DURATION = 60

type call struct {
	startedAt time.Time
	Channel   string
}

// var connMutex sync.Mutex
var frameMutex sync.Mutex
var conn *net.UDPConn

func Init(connetion *net.UDPConn) {
	conn = connetion
	activeCalls = make(map[string]call)
	go durationWatcher()
}
func Dial(addr string, channel string) {
	//check if there is a call currently
	if activeCalls[addr].Channel != channel {
		Hangup(addr)
	}
	frameMutex.Lock()
	activeCalls[addr] = call{
		Channel:   channel,
		startedAt: time.Now(),
	}
	frameMutex.Unlock()
}
func IsInACall(addr string) bool {
	if activeCalls[addr].Channel != "" {
		return true
	}
	return false
}
func Hangup(addr string) {
	deadline := time.Now().Add(time.Second)
	conn.SetWriteDeadline(deadline)
	conn.WriteToUDP([]byte("ok\n"), getAddress(addr))
	frameMutex.Lock()
	delete(activeCalls, addr)
	frameMutex.Unlock()
}
func Relay(addr string, data []byte) {
	if string(data) == "Hangup" {
		Hangup(addr)
		return
	}
	callChannel := activeCalls[addr]
	receptors := phonebook.Get(callChannel.Channel)

	for i := 0; i < len(receptors); i++ {
		go func(i int) {
			if len(data) == 0 {
				return
			}
			// connMutex.Lock()
			deadline := time.Now().Add(time.Second)
			conn.SetWriteDeadline(deadline)
			_, err := conn.WriteTo(data, receptors[i])
			// connMutex.Unlock()
			if err != nil {
				log.Println(err)
				Hangup(addr)
				return
			}
		}(i)
	}

}
func SendOK(addr *net.UDPAddr) {
	// connMutex.Lock()
	deadline := time.Now().Add(time.Second)
	conn.SetWriteDeadline(deadline)
	conn.WriteToUDP([]byte("ok\n"), addr)
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
func durationWatcher() {
	for {
		time.Sleep(time.Minute)
		for key := range activeCalls {
			if time.Since(activeCalls[key].startedAt) >= time.Minute {
				Hangup(key)
			}
		}
	}
}

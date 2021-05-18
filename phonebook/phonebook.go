package phonebook

import (
	"net"
)

var phonebook map[string][]*net.UDPAddr

func Init() {
	phonebook = make(map[string][]*net.UDPAddr)
}
func Register(addr *net.UDPAddr, channel string) {
	found := false
	for i := 0; i < len(phonebook[channel]); i++ {
		if phonebook[channel][i].String() == addr.String() {
			found = true
			break
		}
	}
	if !found {
		phonebook[channel] = append(phonebook[channel], addr)
	}
}
func Reset() {
	phonebook = map[string][]*net.UDPAddr{}
}
func Get(channel string) []*net.UDPAddr {
	return phonebook[channel]
}

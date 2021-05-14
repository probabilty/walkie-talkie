package phonebook

import "net"

var phonebook map[string][]*net.UDPAddr

func Init() {
	phonebook = make(map[string][]*net.UDPAddr)
}
func Register(addr *net.UDPAddr, channel string) {
	phonebook[channel] = append(phonebook[channel], addr)
}
func Reset() {
	phonebook = map[string][]*net.UDPAddr{}
}
func Get(channel string) []*net.UDPAddr {
	return phonebook[channel]
}

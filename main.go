package main

import (
	"fmt"
	"net"
	"strings"
	"walki-talki/frame"
	"walki-talki/phonebook"
	"walki-talki/utils"
)

func init() {

}
func main() {
	PORT := ":8844"
	var hangChan chan *net.UDPAddr = make(chan *net.UDPAddr)
	s, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}
	phonebook.Init(hangChan)
	frame.Init(connection)
	go utils.Hangupcalls(hangChan)
	defer connection.Close()
	for {
		buffer := make([]byte, 10240000)
		// fmt.Printf("address of buffer %p  \n", &buffer)
		n, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}
		go func(buffer []byte, addr *net.UDPAddr, n int) {
			if frame.IsInACall(addr.String()) {
				frame.Relay(addr.String(), (buffer[0 : n-1]))
				return
			}
			if strings.HasPrefix(string(buffer[0:n-1]), "Dial") {
				channel := strings.Split(string(buffer[0:n-1]), " ")
				if len(channel) == 2 {
					frame.Dial(addr.String(), channel[1])
					// frame.SendOK(connection, addr)
				}
			}
			if strings.HasPrefix(string(buffer[0:n-1]), "Register") {
				channel := strings.Split(string(buffer[0:n-1]), " ")
				phonebook.Register(addr, channel[1:])
				frame.SendOK(addr)
			}
		}(buffer, addr, n)
	}
}

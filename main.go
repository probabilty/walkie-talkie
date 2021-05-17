package main

import (
	"fmt"
	"net"
	"strings"
	"walki-talki/frame"
	"walki-talki/phonebook"
)

func init() {
	phonebook.Init()
	frame.Init()
}
func main() {
	PORT := ":8844"

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

	defer connection.Close()
	buffer := make([]byte, 1024)
	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Print("-> ", string(buffer[0:n]))
		go func() {
			if frame.IsInACall(addr.Network()) {
				frame.Relay(connection, addr.Network(), (buffer[0:n]))
				return
			}
			if strings.HasPrefix(string(buffer[0:n]), "Dial") {
				channel := strings.Split(string(buffer[0:n]), " ")
				if len(channel) != 0 {
					frame.Dial(addr.Network(), channel[1])
				}
			}
			if strings.HasPrefix(string(buffer[0:n]), "Register") {
				channel := strings.Split(string(buffer[0:n]), " ")
				for i := 1; i < len(channel); i++ {
					phonebook.Register(addr, channel[i])
				}
			}
		}()
	}
}

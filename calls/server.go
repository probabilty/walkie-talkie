package calls

import (
	"fmt"
	"log"
	"net"
	"strings"
	"walki-talki/frame"
	"walki-talki/phonebook"
	"walki-talki/utils"
)

//Serve Starts a UDP Walki Talki server on port
//The format for port parameter is ":8844"
//To Test the server on terminal:
//nc 127.0.0.1 8844 -u
//To Register new client Type
//Register phon1 phone2 ...
//To Dial phone2 Type
//Dial phone2
//Start sending data
//When Done Type
//Hangup
//Kindly note that if the client address changes, this particular client has to register itself
func Serve(port string) {
	if port == "" {
		port = ":8842"
	}
	PORT := ":8844"
	var hangChan chan *net.UDPAddr = make(chan *net.UDPAddr)
	s, err := net.ResolveUDPAddr("udp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	connection, err := net.ListenUDP("udp", s)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("UDP server is running on%+v\n", s.String())
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

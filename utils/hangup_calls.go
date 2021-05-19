package utils

import (
	"net"
	"walki-talki/frame"
)

func Hangupcalls(notifier chan *net.UDPAddr) {
	for {
		select {
		case msg := <-notifier:
			frame.Hangup(msg.String())
		}
	}
}

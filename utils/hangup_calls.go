package utils

import (
	"net"

	"github.com/probabilty/walki-talki/frame"
)

func Hangupcalls(notifier chan *net.UDPAddr) {
	for {
		select {
		case msg := <-notifier:
			frame.Hangup(msg.String())
		}
	}
}

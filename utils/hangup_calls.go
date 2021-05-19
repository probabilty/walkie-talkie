package utils

import (
	"net"

	"github.com/probabilty/walkie-talkie/frame"
)

func Hangupcalls(notifier chan *net.UDPAddr) {
	for {
		select {
		case msg := <-notifier:
			frame.Hangup(msg.String())
		}
	}
}

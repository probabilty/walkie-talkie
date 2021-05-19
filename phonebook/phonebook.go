package phonebook

import (
	"net"
	"reflect"
	"sync"
)

var phonebook map[string][]*net.UDPAddr
var phonebookLock sync.Mutex

type registration struct {
	address  net.UDPAddr
	channels []string
}

var registrations []registration
var hangChan chan *net.UDPAddr

func Init(hangs chan *net.UDPAddr) {
	phonebook = make(map[string][]*net.UDPAddr)
	hangChan = hangs
}
func Register(addr *net.UDPAddr, channel []string) {
	//findout the old address
	var oldAddress *net.UDPAddr
	//check previous registrations
	for ri, v := range registrations {
		//if the same registration happend before
		if reflect.DeepEqual(v.channels, channel) {
			oldAddress = &v.address
			//remove from phonebook
			for channel := range phonebook {
				for index := 0; index < len(phonebook[channel]); index++ {
					phonebookLock.Lock()
					phonebook[channel] = append(phonebook[channel][:index], phonebook[channel][index+1:]...)
					phonebookLock.Unlock()
				}
			}
			hangChan <- oldAddress
			registrations = append(registrations[:ri], registrations[ri+1:]...)
			break
		}
	}
	for _, v := range channel {
		phonebookLock.Lock()
		phonebook[v] = append(phonebook[v], addr)
		phonebookLock.Unlock()
	}
	registrations = append(registrations, registration{
		address:  *addr,
		channels: channel,
	})
}
func Reset() {
	phonebookLock.Lock()
	phonebook = map[string][]*net.UDPAddr{}
	phonebookLock.Unlock()
}
func Get(channel string) []*net.UDPAddr {
	return phonebook[channel]
}

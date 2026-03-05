package garp

import (
	"log"
)

type Device interface {
	Listen()
	Unicast() UnicastChan
}

type device struct {
	name string
	ip   string
	mac  string

	unicast    UnicastChan // The channel device received message from switch
	switchChan chan string // The channel device send ack message to switch
}

func NewDevice(name, ip, mac string, switchChan chan string) *device {
	return &device{
		name:       name,
		ip:         ip,
		mac:        mac,
		unicast:    make(chan string),
		switchChan: switchChan,
	}
}

func (d *device) Unicast() UnicastChan {
	return d.unicast
}

func (d *device) Listen() {
	for askIP := range d.unicast {
		if d.ip == askIP {
			log.Printf("%s (%s): my MAC is %s", d.name, d.ip, d.mac)
			d.switchChan <- d.mac
		} else {
			log.Printf("%s (%s): IP %s???", d.name, d.ip, askIP)
			d.switchChan <- ""
		}
	}
}

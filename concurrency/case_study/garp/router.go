package garp

import (
	"log"
	"time"
)

type Router interface {
	SendArp()
}

type router struct {
	ipList     []string
	broadcast  BroadcastChan
	listenChan UnicastChan
	jitter     time.Duration
}

func NewRouter(
	ipList []string, broadcast BroadcastChan,
	listenChan UnicastChan, jitter time.Duration,
) (*router, error) {
	if len(ipList) == 0 {
		return nil, ErrEmptyIPList
	}
	if broadcast == nil {
		return nil, ErrMissingBroadcastChannel
	}
	if jitter == 0 {
		jitter = time.Millisecond * 100
	}

	return &router{
		ipList:     ipList,
		broadcast:  broadcast,
		listenChan: listenChan,
		jitter:     jitter,
	}, nil
}

func (r *router) SendArp() {
	for _, ip := range r.ipList {
		// Send ip want to ask to broadcast channel
		time.Sleep(r.jitter)
		log.Printf("Router: I want to ask IP [%s]", ip)
		r.broadcast <- ip

		// Listening for MAC IP reply from switch
		select {
		case ack := <-r.listenChan:
			log.Printf("Router: IP: %s - MAC: %s", ip, ack)
			log.Println("----------------------------------------------")
		case <-time.After(5 * time.Second):
			log.Printf("Router: Timeout waiting for MAC reply for IP: %s", ip)
		}
	}
	close(r.broadcast)
}

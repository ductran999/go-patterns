package garp_test

import (
	"patterns/concurrency/case_study/garp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_SwitchListen(t *testing.T) {
	t.Parallel()

	broadcast := make(chan string, 1)
	ack := make(chan string, 1)
	routerUnicast := make(chan string, 1)
	deviceReceived := make(chan string, 1)

	s := garp.NewSwitch(broadcast, ack, routerUnicast)
	s.RegisterDeviceUnicast(deviceReceived)

	go s.Listen()

	// send IP want to ask
	s.BroadcastChan() <- "192.168.1.23"
	close(broadcast)

	// Get IP from switch and verify it matches
	select {
	case receivedIP := <-deviceReceived:
		assert.Equal(t, "192.168.1.23", receivedIP, "Switch should forward the broadcast IP")
	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for IP from switch")
	}

	// Device send fake MAC via ack channel
	ack <- "AA:BB:CC:DD:EE:FF"

	// Add small delay to ensure switch processes the MAC
	time.Sleep(100 * time.Millisecond)
}

package garp

import (
	"context"
	"log"
	"os/signal"
	"patterns/concurrency/case_study/garp"
	"syscall"
	"time"
)

// initRouter initializes the router with IPs, broadcast, and listen channels.
func initRouter(broadcast chan string, routerListen garp.UnicastChan) (garp.Router, error) {
	ips := []string{
		"192.186.1.12",
		"192.186.1.13",
		"192.186.1.14",
		"192.186.1.17",
	}

	return garp.NewRouter(ips, broadcast, routerListen, time.Second)
}

// GoNetSim is a Go-based network simulator that demonstrates how a router discovers devices in a LAN using IP addresses.
// It simulates ARP (Address Resolution Protocol) communication in a local network topology.
//
// Communication Flow:
// - The router broadcasts ARP requests using a fan-in pattern via broadcastChan.
// - The switch receives these broadcasts and distributes them to all connected devices using a fan-out pattern (unicast).
// - Devices respond with acknowledgment messages, which flow back through ackChan.
//
// This function sets up the core components (router, switch, and devices),
// establishes their communication channels, and runs the simulation until a termination signal is received.
func GoNetSim() {
	// Communication channels
	broadcastChan := make(chan string)
	routerListenChan := make(chan string)
	ackChan := make(chan string)

	// Initialize core garp
	switchDevice := garp.NewSwitch(broadcastChan, ackChan, routerListenChan)
	router, err := initRouter(broadcastChan, routerListenChan)
	if err != nil {
		log.Fatalf("failed to initialize router: %v", err)
	}

	// Create devices
	devices := []garp.Device{
		garp.NewDevice("Device 1", "192.186.1.12", "00:1A:2B:3C:4D:5E", ackChan),
		garp.NewDevice("Device 2", "192.186.1.13", "00:1A:2B:3C:4D:5F", ackChan),
		garp.NewDevice("Device 3", "192.186.1.14", "00:1A:2B:3C:4D:60", ackChan),
		garp.NewDevice("Device 4", "192.186.1.17", "00:1A:2B:3C:4D:50", ackChan),
	}

	// Register devices with the switch
	unicasts := make([]garp.UnicastChan, len(devices))
	for i, dev := range devices {
		unicasts[i] = dev.Unicast()
	}
	switchDevice.RegisterDeviceUnicast(unicasts...)

	// Start all goroutines
	go router.SendArp()
	go switchDevice.Listen()
	for _, dev := range devices {
		go dev.Listen()
	}

	// Graceful shutdown on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("Shutdown app")
}

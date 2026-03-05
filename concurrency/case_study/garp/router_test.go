package garp_test

import (
	"patterns/concurrency/case_study/garp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewRouter_MissingBroadcast(t *testing.T) {
	t.Parallel()

	listen := make(chan string)
	ipList := []string{"192.168.1.10"}

	router, err := garp.NewRouter(ipList, nil, listen, 0)

	require.ErrorIs(t, garp.ErrMissingBroadcastChannel, err)
	assert.Nil(t, router)
}

func Test_NewRouter_MissingIpList(t *testing.T) {
	t.Parallel()

	listen := make(chan string)
	broadcast := make(chan string)
	router, err := garp.NewRouter([]string{}, broadcast, listen, 0)

	require.ErrorIs(t, garp.ErrEmptyIPList, err)
	assert.Nil(t, router)
}

func Test_SendArp_FoundMAC(t *testing.T) {
	t.Parallel()

	listen := make(chan string)
	broadcast := make(chan string)
	ipList := []string{"192.168.1.10", "192.168.1.1"}

	router, err := garp.NewRouter(ipList, broadcast, listen, 0)
	require.NoError(t, err)

	go router.SendArp()
	broadcastIp := <-broadcast
	assert.Equal(t, "192.168.1.10", broadcastIp)

	// Fake MAC
	listen <- "00:1A:2B:3C:4D:50"

	broadcastIp = <-broadcast
	assert.Equal(t, "192.168.1.1", broadcastIp)

	// Fake MAC not found
	listen <- ""
}

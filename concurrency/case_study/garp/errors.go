package garp

import "errors"

var (
	ErrMissingBroadcastChannel = errors.New("missing broadcast channel")
	ErrEmptyIPList             = errors.New("ip list is empty")
)

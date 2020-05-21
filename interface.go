package main

import (
	"fmt"

	peerNumber "github.com/516108736/QKCMainnetMonitor/PeerNumber"

	blockTime "github.com/516108736/QKCMainnetMonitor/BlockTime"
	"github.com/516108736/QKCMainnetMonitor/common"
)

type ModuleInterface interface {
	Check() []string
	PreCheck() error
}

func NewRuler(config common.Config) (ModuleInterface, error) {
	switch config.Module {
	case "BlockTime":
		return blockTime.New(config)
	case "PeerNumber":
		return peerNumber.New(config)
	default:
		panic(fmt.Errorf("not support %v", config.Module))

	}
}

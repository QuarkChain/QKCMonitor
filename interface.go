package main

import (
	"fmt"

	"github.com/516108736/QKCMainnetMonitor/PeerNumber"

	"github.com/516108736/QKCMainnetMonitor/BlockTime"
	"github.com/516108736/QKCMainnetMonitor/common"
)

type ModuleInterface interface {
	Check() []string
	PreCheck() error
}

func NewRuler(config common.Config) (ModuleInterface, error) {
	switch config.Module {
	case "BlockTime":
		return BlockTime.New(config)
	case "PeerNumber":
		return PeerNumber.New(config)
	default:
		panic(fmt.Errorf("not support %v", config.Module))

	}
}

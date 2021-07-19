package main

import (
	"fmt"

	"github.com/516108736/QKCMonitor/poolStats"

	"github.com/516108736/QKCMonitor/blockTime"
	"github.com/516108736/QKCMonitor/common"
	"github.com/516108736/QKCMonitor/peerNumber"
)

type ModuleInterface interface {
	Check() []string
	PreCheck() error
	Summary() []string
}

func NewRuler(config common.Config) (ModuleInterface, error) {
	switch config.Module {
	case "BlockTime":
		return blockTime.New(config)
	case "PeerNumber":
		return peerNumber.New(config)
	case "qpoolStats":
		return poolStats.New(config)
	default:
		panic(fmt.Errorf("not support %v", config.Module))

	}
}

package poolStats

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/516108736/QKCMonitor/common"
)

type PoolStats struct {
	urls string
}

func New(config common.Config) (*PoolStats, error) {
	return &PoolStats{urls: config.TargetHosts[0]}, nil
}

func (b *PoolStats) makeError(host string, errMsg error) string {
	return fmt.Sprintf("PoolStats:%v \nerrMsg:%v", b.urls, errMsg)
}

func (b *PoolStats) Check() []string {
	resp, err := http.Get(b.urls)
	fmt.Println("resp", err, reflect.TypeOf(resp))
	fmt.Println("resp", resp)
	if err != nil {
		return []string{err.Error()}
	}
	return nil
}

func (b *PoolStats) Summary() []string {
	return nil
}

func (b *PoolStats) PreCheck() error {
	return nil
}

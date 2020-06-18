package qkcClient

import (
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ybbus/jsonrpc"
)

type Client struct {
	client jsonrpc.RPCClient
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
}

func HexToUint64(data string) (uint64, error) {
	res, err := hexutil.DecodeUint64(data)
	return res, err
}

func NewClient(host string) *Client {
	client := jsonrpc.NewClient(host)
	return &Client{client: client}
}

func (c *Client) tryCall(method string, params ...interface{}) (*jsonrpc.RPCResponse, error) {
	var (
		errMsg error                = nil
		cnt                         = 0
		resp   *jsonrpc.RPCResponse = nil
	)
	for cnt < 3 {
		time.Sleep(3 * time.Second)
		cnt++
		resp, err := c.client.Call(method, params)
		if err != nil {
			errMsg = err
			continue
		}
		if resp.Error != nil {
			errMsg = resp.Error
			continue
		}
		return resp, errMsg

	}
	return resp, errMsg
}

func (c *Client) GetRootBlockHeight() (uint64, error) {
	resp, err := c.tryCall("getRootBlockByHeight", nil, false)
	if err != nil {
		return 0, err
	}
	height, ok := resp.Result.(map[string]interface{})["height"]
	if !ok {
		return 0, errors.New("resp err")
	}

	return HexToUint64(height.(string))
}

func (c *Client) GetPeers() (int, error) {
	resp, err := c.tryCall("getPeers")
	if err != nil {
		return 0, err
	}
	return len(resp.Result.(map[string]interface{})["peers"].([]interface{})), nil
}

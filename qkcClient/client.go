package qkcClient

import (
	"errors"
	"math/big"

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

func (c *Client) GetRootBlockHeight() (uint64, error) {
	resp, err := c.client.Call("getRootBlockByHeight", nil, false)
	if err != nil {
		return 0, err
	}
	if resp.Error != nil {
		return 0, resp.Error
	}

	height, ok := resp.Result.(map[string]interface{})["height"]
	if !ok {
		return 0, errors.New("resp err")
	}

	return HexToUint64(height.(string))
}

func (c *Client) GetPeers() (int, error) {
	resp, err := c.client.Call("getPeers")
	if err != nil {
		return 0, err
	}
	if resp.Error != nil {
		return 0, resp.Error
	}

	return len(resp.Result.(map[string]interface{})["peers"].([]interface{})), nil
}

package fabric

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type FabricClient struct {
	baseURL       string
	ChaincodePath string
	ChaincodeID   string
}

func NewFabricClient(baseURL string) *FabricClient {
	return &FabricClient{
		baseURL: baseURL,
	}
}

type rpcRequest struct {
	Jsonrpc string        `json:"jsonrpc,omitempty"`
	ID      int64         `json:"id"`
	Method  string        `json:"method,omitempty"` // invoke, query, init
	Params  requestParams `json:"params"`
}

type requestParams struct {
	Type        int         `json:"type"`
	ChaincodeID chaincodeID `json:"chaincodeID"`
	CTorMsg     struct {
		Function string   `json:"function"`
		Args     []string `json:"args"`
	} `json:"ctorMsg"`
}

type chaincodeID struct {
	Path string `json:"path,omitempty"`
	Name string `json:"name,omitempty"`
}

func (c *FabricClient) SendChaincodeRequest(method, function string, args []string) (string, error) {
	req := c.buildCCRequest(method, function, args)
	return c.postCCRequest(req)
}

func (c *FabricClient) buildCCRequest(method, function string, args []string) []byte {
	req := rpcRequest{
		Jsonrpc: "2.0",
		ID:      123,
		Method:  method,
	}
	req.Params.Type = 1
	if method == "deploy" {
		req.Params.ChaincodeID.Path = c.ChaincodePath
	} else {
		req.Params.ChaincodeID.Name = c.ChaincodeID
	}
	req.Params.CTorMsg.Function = function
	req.Params.CTorMsg.Args = args

	ret, _ := json.Marshal(req)
	return ret
}

type rpcResponse struct {
	Jsonrpc string     `json:"jsonrpc,omitempty"`
	Result  *rpcResult `json:"result,omitempty"`
	Error   *rpcError  `json:"error,omitempty"`
	ID      *int64     `json:"id"`
}

type rpcResult struct {
	Status  string    `json:"status,omitempty"`
	Message string    `json:"message,omitempty"`
	Error   *rpcError `json:"error,omitempty"`
}

type rpcError struct {
	Code    int64  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`
}

func (c *FabricClient) postCCRequest(data []byte) (string, error) {
	res, err := http.Post(c.baseURL+"/chaincode", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	resp := new(rpcResponse)
	err = json.NewDecoder(res.Body).Decode(resp)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", errors.New(resp.Error.Data)
	}
	if resp.Result == nil {
		return "", nil
	}
	return resp.Result.Message, nil
}

func (c *FabricClient) WaitTx(txID string) error {
	timer := time.NewTimer(5 * time.Second)
	for {
		select {
		case <-timer.C:
			return fmt.Errorf("transaction timeout")
		case <-time.After(300 * time.Millisecond):
			if c.IsTxInLastBlock(txID) {
				return nil
			}
		}
	}
}

func (c *FabricClient) IsTxInLastBlock(txID string) bool {
	status, err := c.GetChainStatus()
	if err != nil {
		return false
	}
	block, err := c.GetBlock(status.Height - 1)
	if err != nil {
		return false
	}
	for _, tx := range block.Transactions {
		if tx.TxID == txID {
			return true
		}
	}
	return false
}

type Block struct {
	Transactions []Transaction `json:"transactions"`
	StateHash    string        `json:"stateHash"`
}

type Transaction struct {
	Type        int         `json:"type"`
	ChaincodeID chaincodeID `json:"chaincodeID"`
	Payload     string      `json:"payload"`
	TxID        string      `json:"txid"`
}

func (c *FabricClient) GetBlock(index int) (*Block, error) {
	res, err := http.Get(fmt.Sprintf("%s/blocks/%d", c.baseURL, index))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	block := new(Block)
	err = json.NewDecoder(res.Body).Decode(block)
	if err != nil {
		return nil, err
	}
	return block, nil
}

type ChainStatus struct {
	Height            int    `json:"height,omitempty"`
	CurrentBlockHash  string `json:"currentBlockHash,omitempty"`
	PreviousBlockHash string `json:"previousBlockHash,omitempty"`
}

func (c *FabricClient) GetChainStatus() (ChainStatus, error) {
	var chainStatus ChainStatus
	res, err := http.Get(c.baseURL + "/chain")
	if err != nil {
		return chainStatus, err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&chainStatus)
	return chainStatus, err
}

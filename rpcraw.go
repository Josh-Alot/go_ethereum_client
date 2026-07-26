package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const rpcURL = "http://127.0.0.1:8545"

type RpcRequestData struct {
	JsonRpc string   `json:"jsonrpc"`
	ID      int      `json:"id"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}

type RpcResponseData struct {
	JsonRpc string    `json:"jsonrpc"`
	ID      int       `json:"id"`
	Result  string    `json:"result"`
	Error   *RpcError `json:"error"`
}

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func BlockNumber() (uint64, error) {
	msg := RpcRequestData{JsonRpc: "2.0", ID: 1, Method: "eth_blockNumber"}

	req, err := json.Marshal(msg)
	if err != nil {
		return 0, fmt.Errorf("marshal %v: %w", msg, err)
	}

	res, err := http.Post(rpcURL, "application/json", bytes.NewReader(req))
	if err != nil {
		return 0, fmt.Errorf("post %s: %w", rpcURL, err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("read all: %w", err)
	}

	var rpcResponse RpcResponseData
	err = json.Unmarshal(body, &rpcResponse)
	if err != nil {
		return 0, fmt.Errorf("unmarshal: %w", err)
	}

	if rpcResponse.Error != nil {
		return 0, errors.New(rpcResponse.Error.Message)
	}

	height, err := strconv.ParseUint(rpcResponse.Result, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("parse uint %s: %w", rpcResponse.Result, err)
	}

	return height, nil
}

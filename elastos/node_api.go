/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package elastos

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/blocktree/go-owcdrivers/elastosTransaction"
	"github.com/blocktree/openwallet/log"
	"github.com/imroc/req"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

type ClientInterface interface {
	Call(path string, request []interface{}) (*gjson.Result, error)
}

// A Client is a Elastos RPC client. It performs RPCs over HTTP using JSON
// request and responses. A Client must be configured with a secret token
// to authenticate with other Cores on the network.
type Client struct {
	BaseURL string
	// AccessToken string
	Debug  bool
	client *req.Req
	//Client *req.Req
}

type Response struct {
	Code    int         `json:"code,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Message string      `json:"message,omitempty"`
	Id      string      `json:"id,omitempty"`
}

func NewClient(url string, debug bool) *Client {
	c := Client{
		BaseURL: url,
		//	AccessToken: token,
		Debug: debug,
	}

	api := req.New()
	//trans, _ := api.Client().Transport.(*http.Transport)
	//trans.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c.client = api

	return &c
}

// Call calls a remote procedure on another node, specified by the path.
func (c *Client) Call(path string, request []interface{}) (*gjson.Result, error) {

	var (
		body = make(map[string]interface{}, 0)
	)

	if c.client == nil {
		return nil, errors.New("API url is not setup. ")
	}

	authHeader := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic ", // + c.AccessToken,
	}

	//json-rpc
	body["jsonrpc"] = "2.0"
	body["id"] = "1"
	body["method"] = path
	body["params"] = request

	if c.Debug {
		log.Std.Info("Start Request API...")
	}

	r, err := c.client.Post(c.BaseURL, req.BodyJSON(&body), authHeader)

	if c.Debug {
		log.Std.Info("Request API Completed")
	}

	if c.Debug {
		log.Std.Info("%+v", r)
	}

	if err != nil {
		return nil, err
	}

	resp := gjson.ParseBytes(r.Bytes())
	err = isError(&resp)
	if err != nil {
		return nil, err
	}

	result := resp.Get("result")

	return &result, nil
}

// See 2 (end of page 4) http://www.ietf.org/rfc/rfc2617.txt
// "To receive authorization, the client sends the userid and password,
// separated by a single colon (":") character, within a base64
// encoded string in the credentials."
// It is not meant to be urlencoded.
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

//isError 是否报错
func isError(result *gjson.Result) error {
	var (
		err error
	)

	/*
		//failed 返回错误
		{
			"result": null,
			"error": {
				"code": -8,
				"message": "Block height out of range"
			},
			"id": "foo"
		}
	*/

	if !result.Get("error").IsObject() {

		if !result.Get("result").Exists() {
			return errors.New("Response is empty! ")
		}

		return nil
	}

	errInfo := fmt.Sprintf("[%d]%s",
		result.Get("error.code").Int(),
		result.Get("error.message").String())
	err = errors.New(errInfo)

	return err
}

//getBlockHeight 获取区块链高度
func (c *Client) getBlockHeight() (uint64, error) {

	result, err := c.Call("getblockcount", nil)
	if err != nil {
		return 0, err
	}

	return result.Uint() - 1, nil
}

//getBlockHash 根据区块高度获得区块hash
func (c *Client) getBlockHash(height uint64) (string, error) {

	request := []interface{}{
		height,
	}

	result, err := c.Call("getblockhash", request)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

//getBlock 获取区块数据
func (c *Client) getBlock(hash string, format ...uint64) (*Block, error) {

	request := []interface{}{
		hash,
	}

	if len(format) > 0 {
		request = append(request, format[0])
	}

	result, err := c.Call("getblock", request)
	if err != nil {
		return nil, err
	}

	return c.NewBlock(result), nil
}

//getTransaction 获取交易单
func (c *Client) getTransaction(txid string) (*Transaction, error) {

	var (
		result *gjson.Result
		err    error
	)

	request := []interface{}{
		txid,
		true,
	}

	result, err = c.Call("getrawtransaction", request)
	if err != nil {

		request = []interface{}{
			txid,
			1,
		}

		result, err = c.Call("getrawtransaction", request)
		if err != nil {
			return nil, err
		}
	}

	return c.newTx(result), nil
}

//getTxOut 获取交易单输出信息，用于追溯交易单输入源头
func (c *Client) getTxOut(txid string, vout uint64) (*Vout, error) {

	request := []interface{}{
		txid,
		true,
	}

	result, err := c.Call("getrawtransaction", request)
	if err != nil {
		return nil, err
	}

	outputs := gjson.Get(result.Raw, "vout").Array()
	if int(vout+1) > len(outputs) {
		return nil, errors.New("vout is too big in transction :" + txid)
	}

	output := &Vout{}
	for _, out := range outputs {
		if out.Get("n").Uint() == vout {
			output = newTxVout(&out)
		}
	}
	return output, nil
}

//getTxIDsInMemPool 获取待处理的交易池中的交易单IDs
func (c *Client) getTxIDsInMemPool() ([]string, error) {

	var (
		txids = make([]string, 0)
	)

	result, err := c.Call("getrawmempool", nil)
	if err != nil {
		return nil, err
	}

	if !result.IsArray() {
		return nil, errors.New("no query record")
	}

	for _, txid := range result.Array() {
		txids = append(txids, txid.Get("txid").String())
	}

	return txids, nil
}

//getTransaction 获取交易单
func (c *Client) getListUnspent(min uint64, addresses ...string) ([]*Unspent, error) {

	var (
		utxos = make([]*Unspent, 0)
	)

	request := []interface{}{
		addresses,
	}

	result, err := c.Call("listunspent", request)
	if err != nil {
		return nil, err
	}

	array := result.Array()
	for _, a := range array {
		if a.Get("assetid").String() == elastosTransaction.AssetID_ELA && a.Get("confirmations").Uint() >= min {
			utxos = append(utxos, NewUnspent(&a))
		}
	}

	return utxos, nil
}

//estimateFeeRate 预估的没KB手续费率
func (c *Client) estimateFeeRate() (decimal.Decimal, error) {

	feeRate := decimal.Zero

	//估算交易大小 手续费
	request := []interface{}{
		2,
	}

	estimatesmartfee, err := c.Call("estimatesmartfee", request)
	if err != nil {
		return decimal.Zero, errors.New("Failed to get estimatesmartfee!")
	} else {
		feeRate, _ = decimal.NewFromString(estimatesmartfee.String())
	}

	div, _ := decimal.NewFromString("100000000")
	feeRate = feeRate.Div(div)

	return feeRate, nil
}

//sendRawTransaction 广播交易
func (c *Client) sendRawTransaction(txHex string) (string, error) {

	request := []interface{}{
		txHex,
	}

	result, err := c.Call("sendrawtransaction", request)
	if err != nil {
		return "", err
	}

	return result.String(), nil

}

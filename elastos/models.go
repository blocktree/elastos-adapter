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
	"fmt"

	"github.com/blocktree/openwallet/crypto"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tidwall/gjson"
)

//BlockchainInfo 本地节点区块链信息
type BlockchainInfo struct {
	Chain                string `json:"chain"`
	Blocks               uint64 `json:"blocks"`
	Headers              uint64 `json:"headers"`
	Bestblockhash        string `json:"bestblockhash"`
	Difficulty           string `json:"difficulty"`
	Mediantime           uint64 `json:"mediantime"`
	Verificationprogress string `json:"verificationprogress"`
	Chainwork            string `json:"chainwork"`
	Pruned               bool   `json:"pruned"`
}

func NewBlockchainInfo(json *gjson.Result) *BlockchainInfo {
	b := &BlockchainInfo{}
	//解析json
	b.Chain = gjson.Get(json.Raw, "chain").String()
	b.Blocks = gjson.Get(json.Raw, "blocks").Uint()
	b.Headers = gjson.Get(json.Raw, "headers").Uint()
	b.Bestblockhash = gjson.Get(json.Raw, "bestblockhash").String()
	b.Difficulty = gjson.Get(json.Raw, "difficulty").String()
	b.Mediantime = gjson.Get(json.Raw, "mediantime").Uint()
	b.Verificationprogress = gjson.Get(json.Raw, "verificationprogress").String()
	b.Chainwork = gjson.Get(json.Raw, "chainwork").String()
	b.Pruned = gjson.Get(json.Raw, "pruned").Bool()
	return b
}

//Unspent 未花记录
type Unspent struct {

	/*
			{
		        "txid" : "d54994ece1d11b19785c7248868696250ab195605b469632b7bd68130e880c9a",
		        "vout" : 1,
		        "address" : "mgnucj8nYqdrPFh2JfZSB1NmUThUGnmsqe",
		        "account" : "test label",
		        "scriptPubKey" : "76a9140dfc8bafc8419853b34d5e072ad37d1a5159f58488ac",
		        "amount" : 0.00010000,
		        "confirmations" : 6210,
		        "spendable" : true,
		        "solvable" : true
		    }
	*/
	Key     string `storm:"id"`
	TxID    string `json:"txid"`
	Vout    uint64 `json:"vout"`
	Address string `json:"address"`
	//	AccountID string `json:"account" storm:"index"`
	// ScriptPubKey  string `json:"scriptPubKey"`
	Amount        string `json:"amount"`
	Confirmations uint64 `json:"confirmations"`
	Spendable     bool   `json:"spendable"`
	Solvable      bool   `json:"solvable"`
	HDAddress     openwallet.Address
}

func NewUnspent(json *gjson.Result) *Unspent {
	fmt.Println(json.Raw)
	obj := &Unspent{}
	//解析json
	obj.TxID = gjson.Get(json.Raw, "txid").String()
	obj.Vout = gjson.Get(json.Raw, "vout").Uint()
	obj.Address = gjson.Get(json.Raw, "address").String()
	//	obj.AccountID = gjson.Get(json.Raw, "account").String()
	// obj.ScriptPubKey = gjson.Get(json.Raw, "scriptPubKey").String()
	obj.Amount = gjson.Get(json.Raw, "amount").String()
	obj.Confirmations = gjson.Get(json.Raw, "confirmations").Uint()
	//obj.Spendable = gjson.Get(json.Raw, "spendable").Bool()
	obj.Spendable = true
	obj.Solvable = gjson.Get(json.Raw, "solvable").Bool()

	return obj
}

type UnspentSort struct {
	Values     []*Unspent
	Comparator func(a, b *Unspent) int
}

func (s UnspentSort) Len() int {
	return len(s.Values)
}
func (s UnspentSort) Swap(i, j int) {
	s.Values[i], s.Values[j] = s.Values[j], s.Values[i]
}
func (s UnspentSort) Less(i, j int) bool {
	return s.Comparator(s.Values[i], s.Values[j]) < 0
}

type User struct {
	UserKey string `storm:"id"`     // primary key
	Group   string `storm:"index"`  // this field will be indexed
	Email   string `storm:"unique"` // this field will be indexed with a unique constraint
	Name    string // this field will not be indexed
	Age     int    `storm:"index"`
}

type Block struct {
	Hash              string
	Confirmations     uint64
	Merkleroot        string
	tx                []string
	Previousblockhash string
	Height            uint64 `storm:"id"`
	Version           uint64
	Time              uint64
	Fork              bool
	txDetails         []*Transaction
	//	isVerbose         bool
}

func (c *Client) NewBlock(json *gjson.Result) *Block {
	obj := &Block{}
	//解析json
	obj.Height = gjson.Get(json.Raw, "height").Uint()
	obj.Hash = gjson.Get(json.Raw, "hash").String()
	obj.Confirmations = gjson.Get(json.Raw, "confirmations").Uint()
	obj.Merkleroot = gjson.Get(json.Raw, "merkleroot").String()
	obj.Previousblockhash = gjson.Get(json.Raw, "previousblockhash").String()
	obj.Version = gjson.Get(json.Raw, "version").Uint()
	obj.Time = gjson.Get(json.Raw, "time").Uint()

	txs := make([]string, 0)
	txDetails := make([]*Transaction, 0)
	for _, tx := range gjson.Get(json.Raw, "tx").Array() {
		txs = append(txs, tx.String())
	}

	obj.tx = txs
	obj.txDetails = txDetails

	return obj
}

//BlockHeader 区块链头
func (b *Block) BlockHeader(symbol string) *openwallet.BlockHeader {

	obj := openwallet.BlockHeader{}
	//解析json
	obj.Hash = b.Hash
	obj.Confirmations = b.Confirmations
	obj.Merkleroot = b.Merkleroot
	obj.Previousblockhash = b.Previousblockhash
	obj.Height = b.Height
	obj.Version = b.Version
	obj.Time = b.Time
	obj.Symbol = symbol

	return &obj
}

//UnscanRecords 扫描失败的区块及交易
type UnscanRecord struct {
	ID          string `storm:"id"` // primary key
	BlockHeight uint64
	TxID        string
	Reason      string
}

func NewUnscanRecord(height uint64, txID, reason string) *UnscanRecord {
	obj := UnscanRecord{}
	obj.BlockHeight = height
	obj.TxID = txID
	obj.Reason = reason
	obj.ID = common.Bytes2Hex(crypto.SHA256([]byte(fmt.Sprintf("%d_%s", height, txID))))
	return &obj
}

type Transaction struct {
	TxID          string
	Size          uint64
	Version       uint64
	LockTime      int64
	Hex           string
	BlockHash     string
	BlockHeight   uint64
	Confirmations uint64
	Blocktime     int64
	IsCoinBase    bool
	Fees          string
	Decimals      int32
	Type          int32
	Vins          []*Vin
	Vouts         []*Vout
}

type Vin struct {
	Coinbase string
	TxID     string
	Vout     uint64
	N        uint64
	Addr     string
	Value    string
}

type Vout struct {
	AssetID string
	N       uint64
	Addr    string
	Value   string
	Type    string
	Lock    uint64
}

func (c *Client) newTx(json *gjson.Result) *Transaction {
	obj := Transaction{}
	//解析json
	obj.TxID = gjson.Get(json.Raw, "txid").String()
	obj.Version = gjson.Get(json.Raw, "version").Uint()
	obj.LockTime = gjson.Get(json.Raw, "locktime").Int()
	obj.BlockHash = gjson.Get(json.Raw, "blockhash").String()
	obj.Confirmations = gjson.Get(json.Raw, "confirmations").Uint()
	obj.Blocktime = gjson.Get(json.Raw, "blocktime").Int()
	obj.Size = gjson.Get(json.Raw, "size").Uint()
	obj.Type = int32(gjson.Get(json.Raw, "type").Int())
	if obj.Type == 0 {
		obj.IsCoinBase = true
	}

	obj.Decimals = 8
	obj.Vins = make([]*Vin, 0)
	if vins := gjson.Get(json.Raw, "vin"); vins.IsArray() {
		for i, vin := range vins.Array() {
			input := newTxVin(&vin)
			input.N = uint64(i)
			obj.Vins = append(obj.Vins, input)
		}
	}

	obj.Vouts = make([]*Vout, 0)
	if vouts := gjson.Get(json.Raw, "vout"); vouts.IsArray() {
		for _, vout := range vouts.Array() {
			output := newTxVout(&vout)
			obj.Vouts = append(obj.Vouts, output)
		}
	}

	return &obj
}

func newTxVin(json *gjson.Result) *Vin {
	obj := Vin{}
	//解析json
	obj.TxID = gjson.Get(json.Raw, "txid").String()
	obj.Vout = gjson.Get(json.Raw, "vout").Uint()
	if obj.TxID == "0000000000000000000000000000000000000000000000000000000000000000" {
		obj.Coinbase = "true"
	}

	return &obj
}

func newTxVout(json *gjson.Result) *Vout {
	obj := Vout{}
	//解析json
	obj.AssetID = gjson.Get(json.Raw, "assetid").String()
	obj.Value = gjson.Get(json.Raw, "value").String()
	obj.N = gjson.Get(json.Raw, "n").Uint()
	obj.Addr = gjson.Get(json.Raw, "address").String()
	obj.Type = gjson.Get(json.Raw, "type").String()
	obj.Lock = gjson.Get(json.Raw, "outputlock").Uint()

	return &obj
}

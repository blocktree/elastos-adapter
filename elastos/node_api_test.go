package elastos

import (
	"fmt"
	"testing"
)

func Test_getBlockHeight(t *testing.T) {
	c := NewClient("http://127.0.0.1:20336", false)
	height, err := c.getBlockHeight()
	fmt.Println(err)
	fmt.Println(height)
}

func Test_getBlockHash(t *testing.T) {
	c := NewClient("http://127.0.0.1:20336", false)
	height := uint64(398814)
	hash, err := c.getBlockHash(height)
	fmt.Println(err)
	fmt.Println(hash)
}

func Test_getBlock(t *testing.T) {
	c := NewClient("http://127.0.0.1:20336", false)
	hash := "09f3a336ac22a56e32b814af6ef76111044aaac8f64e9ef7e37d4b949d3d9b38"
	block, err := c.getBlock(hash)

	fmt.Println(err)
	fmt.Println(block)
}

func Test_getTransaction(t *testing.T) {
	c := NewClient("http://127.0.0.1:20336", false)
	txid := "db7a07895b31bef6e55c0e4aa1a88688b5ce8642f1d9fd56be868961dcbf5c15"
	tx, err := c.getTransaction(txid)
	fmt.Println(err)
	fmt.Println(tx)
}

func Test_getTxOut(t *testing.T) {
	c := NewClient("http://127.0.0.1:20336", false)
	txid := "4ec5a500314d66507a9b2fa358d15cf7d01c89eca11aeb50ef01de8ac64e1d3a"
	out, err := c.getTxOut(txid, 0)
	fmt.Println(err)
	fmt.Println(out)
}

func Test_getTxIDsInMemPool(t *testing.T) {
	c := NewClient("http://127.0.0.1:20336", false)
	txs, err := c.getTxIDsInMemPool()
	fmt.Println(err)
	fmt.Println(txs)

	tx, err := c.getTransaction(txs[0])
	fmt.Println(err)
	fmt.Println(tx)
}

func Test_getListUnspent(t *testing.T) {
	c := NewClient("http://127.0.0.1:20336", false)
	addresses := []string{"EL9RNsAjWBGCcPaYySM3AZPB4HiX6t93rx"}
	utxos, err := c.getListUnspent(0, addresses[0])
	fmt.Println(err)
	fmt.Println(utxos[0])
}

func Test_estimateFeeRate(t *testing.T) {
	c := NewClient("http://127.0.0.1:20336", false)

	feeRate, err := c.estimateFeeRate()
	fmt.Println(err)
	fmt.Println(feeRate)
}

func Test_tmp(t *testing.T) {
	c := NewClient("http://127.0.0.1:20336", false)

	// addresses := []string{"ESFNULGezraEJGqPzZhPBHWbTGrDmu1UmQ", "EMNg8yRaQ3VYvbb4pCFjLgFBPpbc2whctb"}
	request := []interface{}{
		10,
	}

	resp, err := c.Call("estimatesmartfee", request)

	fmt.Println(err)
	fmt.Println(resp)
}

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
	"math"
	"testing"

	"github.com/codeskyblue/go-sh"
)

var (
	tw *WalletManager
)

func init() {
	tw = NewWalletManager()
	tw.Config.ServerAPI = "http://127.0.0.1:20336"
	tw.WalletClient = NewClient(tw.Config.ServerAPI, true)
}

func TestGOSH(t *testing.T) {
	//text, err := sh.Command("go", "env").Output()
	//text, err := sh.Command("wmd", "version").Output()
	text, err := sh.Command("wmd", "Config", "see", "-s", "btm").Output()
	if err != nil {
		t.Errorf("GOSH failed unexpected error: %v\n", err)
	} else {
		t.Errorf("GOSH output: %v\n", string(text))
	}
}

// func TestEstimateFee(t *testing.T) {
// 	feeRate, _ := tw.EstimateFeeRate()
// 	t.Logf("EstimateFee feeRate = %s\n", feeRate.String())
// 	fees, _ := tw.EstimateFee(10, 2, feeRate)
// 	t.Logf("EstimateFee fees = %s\n", fees.String())
// }

func TestMath(t *testing.T) {
	piece := int64(math.Ceil(float64(67) / float64(30)))

	t.Logf("ceil = %d", piece)
}

func TestPrintConfig(t *testing.T) {
	tw.Config.PrintConfig()
}

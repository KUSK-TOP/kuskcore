package config

import (
	"encoding/hex"

	log "github.com/sirupsen/logrus"

	"kuskcore/consensus"
	"kuskcore/protocol/bc"
	"kuskcore/protocol/bc/types"
)

// AssetIssue asset issue params
type AssetIssue struct {
	NonceStr           string
	IssuanceProgramStr string
	AssetDefinitionStr string
	AssetIDStr         string
	Amount             uint64
}

func (a *AssetIssue) nonce() []byte {
	bs, err := hex.DecodeString(a.NonceStr)
	if err != nil {
		panic(err)
	}

	return bs
}

func (a *AssetIssue) issuanceProgram() []byte {
	bs, err := hex.DecodeString(a.IssuanceProgramStr)
	if err != nil {
		panic(err)
	}

	return bs
}

func (a *AssetIssue) assetDefinition() []byte {
	bs, err := hex.DecodeString(a.AssetDefinitionStr)
	if err != nil {
		panic(err)
	}

	return bs
}

func (a *AssetIssue) id() bc.AssetID {
	var bs [32]byte
	bytes, err := hex.DecodeString(a.AssetIDStr)
	if err != nil {
		panic(err)
	}

	copy(bs[:], bytes)
	return bc.NewAssetID(bs)
}

// GenesisTxs make genesis block txs
func GenesisTxs() []*types.Tx {
	contract, err := hex.DecodeString("0014dc90d8b67f95939950fd4f77db5e466faa92fe14")
	if err != nil {
		log.Panicf("fail on decode genesis tx output control program")
	}

	var txs []*types.Tx
	firstTxData := types.TxData{
		Version: 1,
		Inputs:  []*types.TxInput{types.NewCoinbaseInput([]byte("Information is power. -- 11/11/2024. Computing is power. -- Apr/24/2018."))},
		Outputs: []*types.TxOutput{types.NewOriginalTxOutput(*consensus.KUSKAssetID, consensus.InitKUSKSupply, contract, nil)},
	}
	txs = append(txs, types.NewTx(firstTxData))

	inputs := []*types.TxInput{}
	outputs := []*types.TxOutput{}
	for _, asset := range assetIssues {
		inputs = append(inputs, types.NewIssuanceInput(asset.nonce(), asset.Amount, asset.issuanceProgram(), nil, asset.assetDefinition()))
		outputs = append(outputs, types.NewOriginalTxOutput(asset.id(), asset.Amount, contract, nil))
	}

	secondTxData := types.TxData{Version: 1, Inputs: inputs, Outputs: outputs}
	txs = append(txs, types.NewTx(secondTxData))
	return txs
}

var assetIssues = []*AssetIssue{}

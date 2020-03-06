package vtxo

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func getTxoList1() []*TXO {
	return []*TXO {
		&TXO{
			ID:         "txo1",
			Token:      "cet",
			Amount:     sdk.NewInt(3),
			UsedAmount: sdk.NewInt(1),
		},
		&TXO{
			ID:         "txo2",
			Token:      "cet",
			Amount:     sdk.NewInt(4),
			UsedAmount: sdk.NewInt(1),
		},
		&TXO{
			ID:         "txo3",
			Token:      "cet",
			Amount:     sdk.NewInt(4),
			UsedAmount: sdk.NewInt(3),
		},
	}
}

func prepare1(db DB) {
	for _, txo := range getTxoList1() {
		txo.Owner = "alice"
		db.SaveTXO(txo)
	}
	db.SaveWallet(&Wallet{
		Owner: "alice",
		Token: "cet",
		TXOs:  []string{"txo1", "txo2", "txo3"},
	})
	db.SaveTXO(&TXO{
		Owner:      "bob",
		ID:         "txo10",
		Token:      "cet",
		Amount:     sdk.NewInt(8),
		UsedAmount: sdk.NewInt(0),
	})
	db.SaveWallet(&Wallet{
		Owner: "bob",
		Token: "cet",
		TXOs:  []string{"txo10"},
	})
}

func TestTXO(t *testing.T) {
	txoList := getTxoList1()
	refTXOs, err := DeductFromTXOList(txoList, sdk.NewInt(1))
	assert.Equal(t, nil, err)
	assert.Equal(t, []TXORef{
		{
			ID:     "txo1",
			Amount: sdk.NewInt(1),
		},
	}, refTXOs)
	assert.Equal(t, sdk.NewInt(2), txoList[0].UsedAmount)

	txoList = getTxoList1()
	refTXOs, err = DeductFromTXOList(txoList, sdk.NewInt(2))
	assert.Equal(t, nil, err)
	assert.Equal(t, []TXORef{
		{
			ID:     "txo1",
			Amount: sdk.NewInt(2),
		},
	}, refTXOs)
	assert.Equal(t, sdk.NewInt(3), txoList[0].UsedAmount)

	txoList = getTxoList1()
	refTXOs, err = DeductFromTXOList(txoList, sdk.NewInt(4))
	assert.Equal(t, nil, err)
	assert.Equal(t, []TXORef{
		{
			ID:     "txo1",
			Amount: sdk.NewInt(2),
		},
		{
			ID:     "txo2",
			Amount: sdk.NewInt(2),
		},
	}, refTXOs)
	assert.Equal(t, sdk.NewInt(3), txoList[0].UsedAmount)
	assert.Equal(t, sdk.NewInt(3), txoList[1].UsedAmount)

	txoList = getTxoList1()
	_, err = DeductFromTXOList(txoList, sdk.NewInt(9))
	assert.Equal(t, "Not enough tokens", err.Error())
}

func TestWallet(t *testing.T) {
	keeper := NewKeeper(NewDBInMem())
	prepare1(keeper.db)
	err := keeper.Transfer("alice", "bob", "cet", sdk.NewInt(1), 1000, "tx1")
	assert.Equal(t, nil, err)
	assert.Equal(t, `{"timestamp":0,"owner":"alice","id":"txo1","token":"cet","amount":"3","used_amount":"2","ref_list":null}`,
		string(keeper.db.LoadTXOJson("txo1")))
	assert.Equal(t, `{"timestamp":0,"owner":"alice","id":"txo2","token":"cet","amount":"4","used_amount":"1","ref_list":null}`,
		string(keeper.db.LoadTXOJson("txo2")))
	assert.Equal(t, `{"timestamp":0,"owner":"alice","id":"txo3","token":"cet","amount":"4","used_amount":"3","ref_list":null}`,
		string(keeper.db.LoadTXOJson("txo3")))
	assert.Equal(t, `{"timestamp":1000,"owner":"bob","id":"tx1","token":"cet","amount":"1","used_amount":"0","ref_list":[{"id":"txo1","amount":"1"}]}`,
		string(keeper.db.LoadTXOJson("tx1")))
	assert.Equal(t, `{"owner":"alice","token":"cet","txos":["txo1","txo2","txo3"]}`,
		string(keeper.db.LoadWalletJson("alice","cet")))
	assert.Equal(t, `{"owner":"bob","token":"cet","txos":["txo10","tx1"]}`,
		string(keeper.db.LoadWalletJson("bob","cet")))
	keeper.Close()

	keeper = NewKeeper(NewDBInMem())
	prepare1(keeper.db)
	err = keeper.Transfer("alice", "bob", "cet", sdk.NewInt(2), 1000, "tx2")
	assert.Equal(t, nil, err)
	assert.Equal(t, `{"timestamp":0,"owner":"alice","id":"txo1","token":"cet","amount":"3","used_amount":"3","ref_list":null}`,
		string(keeper.db.LoadTXOJson("txo1")))
	assert.Equal(t, `{"timestamp":0,"owner":"alice","id":"txo2","token":"cet","amount":"4","used_amount":"1","ref_list":null}`,
		string(keeper.db.LoadTXOJson("txo2")))
	assert.Equal(t, `{"timestamp":0,"owner":"alice","id":"txo3","token":"cet","amount":"4","used_amount":"3","ref_list":null}`,
		string(keeper.db.LoadTXOJson("txo3")))
	assert.Equal(t, `{"timestamp":1000,"owner":"bob","id":"tx2","token":"cet","amount":"2","used_amount":"0","ref_list":[{"id":"txo1","amount":"2"}]}`,
		string(keeper.db.LoadTXOJson("tx2")))
	assert.Equal(t, `{"owner":"alice","token":"cet","txos":["txo2","txo3"]}`,
		string(keeper.db.LoadWalletJson("alice","cet")))
	assert.Equal(t, `{"owner":"bob","token":"cet","txos":["txo10","tx2"]}`,
		string(keeper.db.LoadWalletJson("bob","cet")))
	keeper.Close()

	keeper = NewKeeper(NewDBInMem())
	prepare1(keeper.db)
	err = keeper.Transfer("alice", "tom", "cet", sdk.NewInt(4), 1000, "tx3")
	assert.Equal(t, nil, err)
	assert.Equal(t, `{"timestamp":0,"owner":"alice","id":"txo1","token":"cet","amount":"3","used_amount":"3","ref_list":null}`,
		string(keeper.db.LoadTXOJson("txo1")))
	assert.Equal(t, `{"timestamp":0,"owner":"alice","id":"txo2","token":"cet","amount":"4","used_amount":"3","ref_list":null}`,
		string(keeper.db.LoadTXOJson("txo2")))
	assert.Equal(t, `{"timestamp":0,"owner":"alice","id":"txo3","token":"cet","amount":"4","used_amount":"3","ref_list":null}`,
		string(keeper.db.LoadTXOJson("txo3")))
	assert.Equal(t, `{"timestamp":1000,"owner":"tom","id":"tx3","token":"cet","amount":"4","used_amount":"0","ref_list":[{"id":"txo1","amount":"2"},{"id":"txo2","amount":"2"}]}`,
		string(keeper.db.LoadTXOJson("tx3")))
	assert.Equal(t, `{"owner":"alice","token":"cet","txos":["txo2","txo3"]}`,
		string(keeper.db.LoadWalletJson("alice","cet")))
	assert.Equal(t, `{"owner":"tom","token":"cet","txos":["tx3"]}`,
		string(keeper.db.LoadWalletJson("tom","cet")))
	keeper.Close()

	keeper = NewKeeper(NewDBInMem())
	prepare1(keeper.db)
	err = keeper.Transfer("alice", "bob", "cet", sdk.NewInt(9), 1000, "tx4")
	assert.Equal(t, "Not enough tokens", err.Error())
	keeper.Close()

}


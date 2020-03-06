package vtxo

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TXORef struct {
	ID     string  `json:"id"`     // another TXO's ID
	Amount sdk.Int `json:"amount"` // how much amount went from this TXO to this TXO
}

type TXO struct {
	Timestamp  int64    `json:"timestamp"`   // when this TXO was created, in nanoseconds
	Owner      string   `json:"owner"`       // who owns this TXO
	ID         string   `json:"id"`          // a unique id of this TXO
	Token      string   `json:"token"`       // the asset's name in this TXO
	Amount     sdk.Int  `json:"amount"`      // the amount of this asset when this TXO was created
	UsedAmount sdk.Int  `json:"used_amount"` // the amount that has been referenced by later TXOs
	RefList    []TXORef `json:"ref_list"`
}

func (txo *TXO) IsAllUsed() bool {
	return txo.Amount.Equal(txo.UsedAmount)
}

type Wallet struct {
	Owner  string   `json:"owner"` // who owns this wallet
	Token  string   `json:"token"` // the asset's name in this wallet
	TXOs   []string `json:"txos"`  // the TXOs contained in this wallet
}

func DeductFromTXOList(txoList []*TXO, amount sdk.Int) (refTXOs []TXORef, err error) {
	for _, txo := range txoList {
		remain := txo.Amount.Sub(txo.UsedAmount)
		if amount.GT(remain) {
			refTXOs = append(refTXOs, TXORef{ID: txo.ID, Amount: remain})
			txo.UsedAmount = txo.Amount
			amount = amount.Sub(remain)
		} else {
			refTXOs = append(refTXOs, TXORef{ID: txo.ID, Amount: amount})
			txo.UsedAmount = txo.UsedAmount.Add(amount)
			amount = sdk.ZeroInt()
			break
		}
	}
	if !amount.IsZero() {
		return nil, errors.New("Not enough tokens")
	}
	return
}

type Keeper struct {
	db DB
}

func NewKeeper(db DB) Keeper {
	return Keeper{db: db}
}

func (k *Keeper) Close() {
	k.db.Close()
}

func (k *Keeper) Transfer(src, dst, token string, amount sdk.Int, timestamp int64, id string) error {
	var refTXOs []TXORef
	if len(src) != 0 {
		srcWallet := k.db.LoadWallet(src, token)
		if srcWallet == nil {
			return errors.New("Can not find the source wallet")
		}
		srcTXOs, err := k.db.LoadTXOs(srcWallet.TXOs)
		if err != nil {
			return err
		}
		refTXOs, err = DeductFromTXOList(srcTXOs, amount)
		if err != nil {
			return err
		}
		for _, txo := range srcTXOs {
			k.db.SaveTXO(txo)
		}
		if srcTXOs[len(refTXOs)-1].IsAllUsed() {
			srcWallet.TXOs = srcWallet.TXOs[len(refTXOs):]
		} else {
			srcWallet.TXOs = srcWallet.TXOs[len(refTXOs)-1:]
		}
		k.db.SaveWallet(srcWallet)
	}
	dstWallet := k.db.LoadWallet(dst, token)
	if dstWallet == nil {
		dstWallet = &Wallet{Owner: dst, Token: token}
	}
	newTXO := &TXO {
		Timestamp:  timestamp,
		Owner:      dst,
		ID:         id,
		Token:      token,
		Amount:     amount,
		UsedAmount: sdk.ZeroInt(),
		RefList:    refTXOs,
	}
	k.db.SaveTXO(newTXO)
	dstWallet.TXOs = append(dstWallet.TXOs, id)
	k.db.SaveWallet(dstWallet)
	return nil
}



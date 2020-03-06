package vtxo

import (
	"encoding/json"
	"errors"

	dbm "github.com/tendermint/tm-db"
)

const (
	TXOByte = byte(100)
	WalletByte = byte(101)
)

type DB struct {
	kv dbm.DB
}

func NewDBInMem() DB {
	return DB{kv: dbm.NewMemDB()}
}

func NewDB(name string, dir string) (DB, error) {
	kv, err := dbm.NewGoLevelDB(name, dir)
	return DB{kv: kv}, err
}

func (db *DB) Close() {
	db.kv.Close()
}

func (db *DB) SaveTXO(txo *TXO) {
	bz, err := json.Marshal(txo)
	if err != nil {
		panic(err)
	}
	db.kv.SetSync(append([]byte{TXOByte}, []byte(txo.ID)...), bz)
}

func (db *DB) LoadTXOJson(id string) []byte {
	return db.kv.Get(append([]byte{TXOByte}, []byte(id)...))
}

func (db *DB) LoadTXO(id string) *TXO {
	var txo TXO
	bz := db.LoadTXOJson(id)
	if len(bz) == 0 {
		return nil
	}
	err := json.Unmarshal(bz, &txo)
	if err != nil {
		panic(err)
	}
	return &txo
}

func (db *DB) LoadTXOs(idList []string) (res []*TXO, err error) {
	for _, id := range idList {
		txo := db.LoadTXO(id)
		if txo == nil {
			err = errors.New("Cannot find TXO: "+id)
			return
		}
		res = append(res, txo)
	}
	return
}

func (db *DB) SaveWallet(wallet *Wallet) {
	key := make([]byte, 0, len(wallet.Owner)+len(wallet.Token)+2)
	key = append(key, WalletByte)
	key = append(key, []byte(wallet.Owner)...)
	key = append(key, []byte(":")...)
	key = append(key, []byte(wallet.Token)...)
	bz, err := json.Marshal(wallet)
	if err != nil {
		panic(err)
	}
	db.kv.SetSync(key, bz)
}

func (db *DB) LoadWalletJson(owner, token string) []byte {
	key := make([]byte, 0, len(owner)+len(token)+2)
	key = append(key, WalletByte)
	key = append(key, []byte(owner)...)
	key = append(key, []byte(":")...)
	key = append(key, []byte(token)...)
	return db.kv.Get(key)
}

func (db *DB) LoadWallet(owner, token string) *Wallet {
	bz := db.LoadWalletJson(owner, token)
	var wallet Wallet
	if len(bz) == 0 {
		return nil
	}
	err := json.Unmarshal(bz, &wallet)
	if err != nil {
		panic(err)
	}
	return &wallet
}


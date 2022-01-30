package blockchain

import (
	"errors"
	"fmt"
	"gotutorial/utils"
	"gotutorial/wallet"
	"sync"
	"time"
)

const (
	minerReward = 50
)

var ErrNotEnough = errors.New("Not Enough Coin")
var ErrInvalidTx = errors.New("Invalid Tx")

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

type TxIn struct {
	TxID      string `json:"txId"`
	Index     int    `json:"Index"`
	Signature string `json:"signature"`
}

type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

type UTxOut struct {
	TxID   string `json:"txId"`
	Index  int    `json:"index"`
	Amount int    `json:"Amount"`
}

type mempool struct {
	Txs map[string]*Tx
	m   sync.Mutex
}

type MemTxs map[string]*Tx

var m *mempool
var memOnce sync.Once

func Mempool() *mempool {
	memOnce.Do(func() {
		m = &mempool{}
		m.Txs = make(MemTxs)
	})
	return m
}

func (tx *Tx) getId() {
	tx.Id = utils.Hash(tx)
}

func (tx *Tx) sign() {
	for _, txIn := range tx.TxIns {
		txIn.Signature = wallet.Sign(wallet.Wallet(), tx.Id)
	}
}

func validate(tx *Tx) bool {
	valid := true
	for _, txIn := range tx.TxIns {
		prevTx := FindTx(Blockchain(), txIn.TxID)
		if prevTx == nil {
			valid = false
			break
		}
		address := prevTx.TxOuts[txIn.Index].Address
		valid = wallet.Verify(txIn.Signature, tx.Id, address)

	}
	return valid
}

func MakeCoinbaseTx(address string) *Tx {
	tx := Tx{
		Id:        "",
		Timestamp: 0,
		TxIns: []*TxIn{
			{
				"",
				-1,
				"COINBASE",
			},
		},
		TxOuts: []*TxOut{
			{
				address,
				minerReward,
			},
		},
	}
	tx.getId()
	return &tx
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if BlanceByAddress(Blockchain(), from) < amount {
		fmt.Println("err coing")
		return nil, ErrNotEnough
	}

	var txIns []*TxIn
	var txOuts []*TxOut

	total := 0
	uTxOuts := UTxOutsByAddress(Blockchain(), from)
	fmt.Print(uTxOuts)
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}

		txIn := &TxIn{
			TxID:  uTxOut.TxID,
			Index: uTxOut.Index,
		}

		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}

	if change := total - amount; change > 0 {
		changeTxOutput := &TxOut{
			Address: from,
			Amount:  change,
		}
		txOuts = append(txOuts, changeTxOutput)
	}
	txOut := &TxOut{
		Address: to,
		Amount:  amount,
	}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxOuts:    txOuts,
		TxIns:     txIns,
	}
	tx.getId()
	tx.sign()
	valid := validate(tx)
	if !valid {
		return nil, ErrInvalidTx
	}
	return tx, nil

	// var txIns []*TxIn
	// var txOuts []*TxOut

	// total := 0
	// oldTxOuts := Blockchain().TxOutsByAddress(from)
	// for _, e := range oldTxOuts {
	// 	if total >= amount {
	// 		break
	// 	}
	// 	txIns = append(txIns, &TxIn{e.Owner, e.Amount})
	// 	total += e.Amount
	// }

	// change := total - amount
	// if change > 0 {
	// 	txOuts = append(txOuts, &TxOut{from, change})
	// }
	// txOut := &TxOut{to, amount}
	// txOuts = append(txOuts, txOut)
	// tx := Tx{
	// 	Id:        "",
	// 	Timestamp: int(time.Now().Unix()),
	// 	TxIns:     txIns,
	// 	TxOuts:    txOuts,
	// }
	// tx.getId()
	// return &tx, nil
}

func (m *mempool) Transactions() MemTxs {
	m.m.Lock()
	defer m.m.Unlock()
	return m.Txs
}

func (m *mempool) AddTx(to string, amount int) (*Tx, error) {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return nil, err
	}
	m.Txs[tx.Id] = tx
	return tx, nil
}

func (m *mempool) AddPeerTx(tx *Tx) {
	m.m.Lock()
	defer m.m.Unlock()

	m.Txs[tx.Id] = tx
}

func (m *mempool) TxToConfirm() []*Tx {
	coinBase := MakeCoinbaseTx(wallet.Wallet().Address)
	var txs []*Tx

	for _, tx := range m.Txs {
		txs = append(txs, tx)
	}
	txs = append(txs, coinBase)
	m.Txs = make(MemTxs)
	return txs
}

func isOnMempool(uTxOut *UTxOut) bool {
	exist := false
outer:
	for _, tx := range Mempool().Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exist = true
				break outer
			}

		}
	}
	return exist
}

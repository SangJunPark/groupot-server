package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"gotutorial/db"
	"gotutorial/utils"
	"sync"
)

var ErrNotFound = errors.New("block not found")

const (
	defaultDifficulty  = 2
	difficultyInterval = 5
	blockInterval      = 2
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
	m                 sync.Mutex
}

var b *blockchain

var once sync.Once

func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func (chain *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(chain))
}

func (chain *blockchain) restore(data []byte) {
	utils.FromBytes(chain, data)
}

func difficulty(chain *blockchain) int {
	if chain.Height == 0 {
		return defaultDifficulty
	} else if chain.Height%difficultyInterval == 0 {
		return recalculateDifficulty(chain)
	} else {
		return chain.CurrentDifficulty
	}
}

func recalculateDifficulty(chain *blockchain) int {
	allBlocks := Blocks(chain)
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculatedBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockInterval
	if actualTime > expectedTime-2 {
		return chain.CurrentDifficulty - 1
	} else if actualTime < expectedTime {
		return chain.CurrentDifficulty + 1
	} else {
		return chain.CurrentDifficulty
	}
}

func Txs(chain *blockchain) (txs []*Tx) {
	for _, block := range Blocks(chain) {
		txs = append(txs, block.Transactions...)
	}
	return txs
}

func FindTx(chain *blockchain, tragetId string) *Tx {
	for _, tx := range Txs(chain) {
		if tx.Id == tragetId {
			return tx
		}
	}
	return nil
}

func UTxOutsByAddress(chain *blockchain, address string) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)
	for _, block := range Blocks(chain) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Signature == "COINBASE" {
					break
				}
				if FindTx(Blockchain(), input.TxID).TxOuts[input.Index].Address == address {
					creatorTxs[input.TxID] = true
					fmt.Println(creatorTxs)
				}
			}
			for index, output := range tx.TxOuts {
				if output.Address == address {
					if _, ok := creatorTxs[tx.Id]; !ok {
						uTxOut := &UTxOut{tx.Id, index, output.Amount}
						if !isOnMempool(uTxOut) {
							fmt.Println("TEST :" + tx.Id)

							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

func BlanceByAddress(chain *blockchain, address string) int {
	amt := 0
	for _, e := range UTxOutsByAddress(chain, address) {
		fmt.Println("amount :")
		fmt.Println(e.Amount)

		amt += e.Amount
	}
	return amt
}

func Blockchain() *blockchain {
	once.Do(func() {
		b = &blockchain{Height: 0}
		checkpoint := db.Checkpoint()
		if checkpoint == nil {
			b.AddBlock()
		} else {
			b.restore(checkpoint)
		}
		//db.DB().
	})
	fmt.Println("blockchain")
	return b
}

func Status(chain *blockchain) *blockchain {
	fmt.Println("status")
	chain.m.Lock()
	defer chain.m.Unlock()
	return chain
}

func Blocks(chain *blockchain) []*Block {
	fmt.Println("blocks")

	chain.m.Lock()
	defer chain.m.Unlock()
	hashCursor := chain.NewestHash
	var blocks []*Block
	for {
		block, _ := FindBlock(hashCursor)
		fmt.Println("blocks")

		if block == nil {
			break
		}
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

// func getLastHash() string {
// 	totalBlocks := (Blockchain().blocks)
// 	if len(totalBlocks) > 0 {
// 		return totalBlocks[len(totalBlocks)-1].Hash
// 	}
// 	return ""
// }

func (chain *blockchain) Burn(b *Block) {
	db.BurnBlock(b.Hash)
}

func (chain *blockchain) Burn2(b *Block) {
	db.SaveBlock(b.Hash, []byte(b.Data))
}

func (chain *blockchain) AddBlock() *Block {
	fmt.Println("add block")
	chain.m.Lock()
	defer chain.m.Unlock()
	block := createBlock(chain.NewestHash, chain.Height+1, difficulty(chain))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	b.persist()
	return block
}

func (chain *blockchain) AddPeerBlock(block *Block) {
	fmt.Println("add peer block")

	chain.m.Lock()
	Mempool().m.Lock()

	defer chain.m.Unlock()
	defer Mempool().m.Unlock()

	fmt.Println("New Block Added")

	chain.Height += 1
	chain.NewestHash = block.Hash
	chain.CurrentDifficulty = block.Difficulty

	chain.persist()
	b.persist()

	for _, tx := range Mempool().Txs {
		_, ok := Mempool().Txs[tx.Id]
		if ok {
			delete(Mempool().Txs, tx.Id)
		}
	}
}

func (chain *blockchain) Replace(newBlocks []*Block) {
	chain.m.Lock()
	defer chain.m.Unlock()
	chain.CurrentDifficulty = newBlocks[0].Difficulty
	chain.Height = newBlocks[0].Height
	chain.NewestHash = newBlocks[0].Hash
	chain.persist()
	db.EmptyBlocks()
	for _, b := range newBlocks {
		b.persist()
	}
}

// func (chain *blockchain) Block(id int) (*Block, error) {
// 	if len(chain.blocks) <= id {
// 		return chain.blocks[id-1], nil
// 	}
// 	return nil, ErrNotFound
// 	// for _, e := range chain.blocks {
// 	// 	if id == e.Height {
// 	// 		return e
// 	// 	}
// 	// }
// 	// return nil
// }

// func (chain *blockchain) AllBlocks() []*Block {
// 	return chain.blocks
// }

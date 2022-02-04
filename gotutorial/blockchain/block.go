package blockchain

import (
	"gotutorial/db"
	"gotutorial/utils"
	"strings"
	"time"
)

type Block struct {
	Height       int    `json:"height"` // number of block
	Data         string `json:"data"`
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitEmpty"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	Timestamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transactions"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		//blockAsString := fmt.Sprint(b)
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		//fmt.Printf("%s %s %d %d\n\n\n", blockAsString, hash, b.Difficulty, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

func createBlock(prevHash string, height int, difficulty int) *Block {
	block := Block{
		Height:     height,
		Hash:       "",
		PrevHash:   prevHash,
		Difficulty: difficulty,
		Nonce:      0,
		Timestamp:  0,
	}
	// payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	// block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.mine()
	block.Transactions = Mempool().TxToConfirm()
	defer block.persist()
	return &block
}

func ProofOfWork() {

}

func ProofOfStaking() {

}

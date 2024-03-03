package domain

// Block
type Block struct {
	Hash      string `json:"hash"`
	PrevHash  string `json:"prev_hash"`
	Height    int64  `json:"height"`
	MaxHeight int64  `json:"max_height"`
	Timestamp string `json:"timestamp"`
	Txs       []*Tx  `json:"txs"`
}

func NewBlock() {}

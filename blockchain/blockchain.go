package blockchain

type block struct {
	data string
	hash string
	prevHash string
}

type blockchain struct {
	blocks []block
}

// singleton pattern: only one instance
var b * blockchain

func GetBlockchain() *blockchain {
	// 초기화 되었는지 확인하고 딱 한번 생성
	if b == nil {
		b = &blockchain{}
	}
	return b
}



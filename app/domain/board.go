package domain

import (
	"square-face-tetris/app/constants"
)

// ボードの定義
type Board [][]int // 20行10列のボード

func (b *Board) Init() {
	// ボードを定義
	*b = make([][]int, constants.BoardHeight)
	for i := range *b {
		(*b)[i] = make([]int, constants.BoardWidth)
	}
}
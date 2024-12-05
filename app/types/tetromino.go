package types

import (
	"image/color"
)

// テトリミノの定義
type Tetromino struct {
	X, Y      int       // テトリミノの位置
	Color     color.Color // テトリミノの色
	Shape     [][]int   // テトリミノの形状（回転可能）
	Rotation  int       // 回転状態（0, 90, 180, 270度）
}

// 各テトリミノの形状を定義
var Tetrominos = []Tetromino{
	{
		Color: color.RGBA{255, 0, 0, 255}, // 赤
		Shape: [][]int{
			{1, 1, 1, 1}, // I
		},
	},
	{
		Color: color.RGBA{0, 255, 0, 255}, // 緑
		Shape: [][]int{
			{1, 1},
			{1, 1}, // O
		},
	},

}


package types

import (
	"image/color"
)

// テトリミノの定義
type Tetromino struct {
	X, Y      int         // テトリミノの位置
	Color     color.Color // テトリミノの色
	Shape     [][]int     // テトリミノの形状（回転可能）
	Rotation  int         // 回転状態（0, 90, 180, 270度）
}

// 各テトリミノの形状を定義
var Tetrominos = []Tetromino{
	{
		Color: color.RGBA{255, 0, 0, 255}, // 赤 - I
		Shape: [][]int{
			{1, 1, 1, 1}, // 横一列
		},
	},
	{
		Color: color.RGBA{0, 255, 0, 255}, // 緑 - O
		Shape: [][]int{
			{1, 1},
			{1, 1}, // 正方形
		},
	},
	{
		Color: color.RGBA{0, 0, 255, 255}, // 青 - T
		Shape: [][]int{
			{0, 1, 0},
			{1, 1, 1}, // T字型
		},
	},
	{
		Color: color.RGBA{255, 165, 0, 255}, // オレンジ - L
		Shape: [][]int{
			{1, 0},
			{1, 0},
			{1, 1}, // L字型
		},
	},
	{
		Color: color.RGBA{0, 255, 255, 255}, // 水色 - J
		Shape: [][]int{
			{0, 1},
			{0, 1},
			{1, 1}, // 逆L字型
		},
	},
	{
		Color: color.RGBA{255, 255, 0, 255}, // 黄 - S
		Shape: [][]int{
			{0, 1, 1},
			{1, 1, 0}, // S字型
		},
	},
	{
		Color: color.RGBA{128, 0, 128, 255}, // 紫 - Z
		Shape: [][]int{
			{1, 1, 0},
			{0, 1, 1}, // Z字型
		},
	},
}

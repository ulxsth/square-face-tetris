package types

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"square-face-tetris/app/constants"
)

// ボードの定義
type Board [constants.BoardHeight][constants.BoardWidth]int // 20行10列のボード

// ボードの描画
func (b *Board) Draw(screen *ebiten.Image) {
	for y := 0; y < constants.BoardHeight; y++ {
		for x := 0; x < constants.BoardWidth; x++ {
			if b[y][x] == 1 {
				blockImage := ebiten.NewImage(constants.BlockSize, constants.BlockSize)
				blockImage.Fill(color.RGBA{0, 0, 255, 255}) // 青
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(float64(x*constants.BlockSize), float64(y*constants.BlockSize))
				screen.DrawImage(blockImage, opts)
			}
		}
	}
}


package main

import (
	"log"
	"time"
	"github.com/hajimehoshi/ebiten/v2"
	"square-face-tetris/app/constants"
	"square-face-tetris/app/types"
)

func main() {
	// ゲームインスタンスの生成
	game := &types.Game{
		DropInterval: 2 * time.Second, // 落下間隔を2秒に設定
		KeyState:     make(map[ebiten.Key]bool), // キーの押下状態を管理
	}
	err := game.Init()
	if err != nil {
		log.Fatalf("Failed to initialize the game: %v", err)
	}

	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetWindowTitle("Tetris")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

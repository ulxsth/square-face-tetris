package main

import (
	"log"
	"time"
	"github.com/hajimehoshi/ebiten/v2"
	"square-face-tetris/app/constants"
	"square-face-tetris/app/domain"
	"square-face-tetris/app/domain/game"
)

func main() {
	// ゲームインスタンスの生成
	gameWrapper := &game.GameWrapper{
		Game: domain.Game{
			DropInterval: 2 * time.Second,        // ブロックの落下間隔を2秒に設定
			KeyState:     make(map[ebiten.Key]bool), // キー入力の状態を管理
		},
	}
	// ゲームの初期化
	err := gameWrapper.Init()
	if err != nil {
		log.Fatalf("ゲームの初期化に失敗しました: %v", err)
	}

	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetWindowTitle("Tetris")
	if err := ebiten.RunGame(gameWrapper); err != nil {
		log.Fatal(err)
	}
}

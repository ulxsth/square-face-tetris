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
	Game := &types.Game{
		DropInterval: 2 * time.Second, // 落下間隔を2秒に設定
		KeyState:     make(map[ebiten.Key]bool), // キーの押下状態を管理
	}

	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetWindowTitle("Tetris")
	if err := ebiten.RunGame(Game); err != nil {
		log.Fatal(err)
	}
}

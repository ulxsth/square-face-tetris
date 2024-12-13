package types

import (
	"time"
	// "github.com/esimov/pigo/wasm/detector"
	"github.com/hajimehoshi/ebiten/v2"
)

// ゲームの状態
type Game struct {
	Board        Board               // 10x20 のボード
	Current      *Tetromino          // 現在のテトリミノ
	LastDrop     time.Time           // 最後にテトリミノが落下した時刻
	DropInterval time.Duration       // 落下間隔
	KeyState     map[ebiten.Key]bool // キーの押下状態
	CanvasImage  *ebiten.Image       // canvas から取得した画像を保持するフィールドを追加
	StartTime    time.Time           // ゲーム開始時刻
	TimeLimit    time.Duration       // タイムリミット
	State        string              // ゲームの状態
	Score        int                 // スコア
}

// キーが離された場合に状態をリセット
func (g *Game) ResetKeyState() {
	for key := range g.KeyState {
		if !ebiten.IsKeyPressed(key) {
			g.KeyState[key] = false // キーが離されたら状態をリセット
		}
	}
}


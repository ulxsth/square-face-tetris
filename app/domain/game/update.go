package game

import (
	"log"
	"square-face-tetris/app/domain/wasm"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// ゲームの状態更新
func (g *GameWrapper) Update() error {

	switch g.Game.State {
	case "start":
		g.updateStart()
	case "playing":
		g.updatePlaying()
	case "showingScore":
		g.updateShowingScore()
	}
	return nil
}


func (g *GameWrapper) updateStart() {
	// スコア画面ではスペースキーを押すと終了
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		err := g.ResetGame() // ゲームを初期化
		if err != nil {
				log.Fatalf("Failed to initialize the game: %v", err)
		}
	}
}


// プレイ中の状態を更新
func (g *GameWrapper) updatePlaying() {
	// タイムリミットを超えている場合はスコア画面へ遷移
	if time.Since(g.Game.StartTime) >= g.Game.TimeLimit {
		g.Game.State = "showingScore"
		return
	}

	if g.Game.Current == nil {
		g.Game.ShiftTetrominoQueue()
	}

	// ユーザー入力でテトリミノを操作
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && g.Game.IsValidPosition(g.Game.Current, -1, 0) {
		g.Game.Current.X -= 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && g.Game.IsValidPosition(g.Game.Current, 1, 0) {
		g.Game.Current.X += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) && g.Game.IsValidPosition(g.Game.Current, 0, 1) {
		g.Game.Current.Y += 1
	}

	// 回転用ボタンの処理（1回の入力で1回だけ回転）
	if ebiten.IsKeyPressed(ebiten.KeyUp) && !g.Game.KeyState[ebiten.KeyUp] {
		// 回転処理を試みる
		oldX, oldY := g.Game.Current.X, g.Game.Current.Y
		g.Game.RotateTetromino()

		// 範囲外または重なりがある場合、位置を調整
		for g.Game.IsOutOfBounds() == "left" {
			g.Game.Current.X++
		}
		for g.Game.IsOutOfBounds() == "right" {
			g.Game.Current.X--
		}
		for g.Game.IsOutOfBounds() == "bottom" || g.Game.IsOverlapping() {
			g.Game.Current.Y--
		}

		// 調整後も無効な場合は回転をキャンセル
		if g.Game.IsOutOfBounds() != "" || g.Game.IsOverlapping() {
			g.Game.Current.X = oldX
			g.Game.Current.Y = oldY
			// 回転を元に戻す（`rotateTetromino` のロジックを逆回転にする必要あり）
		}

		// 回転ボタンの押下を記録
		g.Game.KeyState[ebiten.KeyUp] = true
	}

	// 2秒間隔で落下
	if time.Since(g.Game.LastDrop) > g.Game.DropInterval {
		if g.Game.IsValidPosition(g.Game.Current, 0, 1) {
			g.Game.Current.Y += 1
		} else {
			// テトリミノが固定されるべき条件を満たす
			g.Game.LockTetromino()

			// 最上段にブロックがあるか確認（ゲームオーバーの判定）
			if g.isTopRowFilled() {
				g.Game.State = "showingScore" // ゲームオーバー状態に遷移
				return
			}
		}
		g.Game.LastDrop = time.Now() // 落下タイマーをリセット
	}

	// キーが離された場合に状態をリセット（回転だけリセット）
	g.Game.ResetKeyState()

	wasm.UpdateCamera()
}

// 最上段が埋まっているか確認
func (g *GameWrapper) isTopRowFilled() bool {
	topRow := g.Game.Board[0] // 最上段の行を取得
	for _, cell := range topRow {
		if cell != 0 { // 0 以外ならブロックがある
			return true
		}
	}
	return false
}

func (g *GameWrapper) updateShowingScore() {
	// スコア画面ではスペースキーを押すと終了
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		err := g.ResetGame() // ゲームを初期化
		if err != nil {
				log.Fatalf("Failed to initialize the game: %v", err)
		}
	}
}

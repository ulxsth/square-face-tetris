package game

import (
	"square-face-tetris/app/constants"

	"bytes"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)


type GameWrapper struct {
	Game Game
}

var (
	mplusFaceSource *text.GoTextFaceSource
)

// ゲームの初期化（タイマーの設定を追加）
func (g *GameWrapper) Init() error {
	// ゲーム全体の初期設定
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
			return err
	}
	mplusFaceSource = s
	g.Game.State = "start"
	return nil
}

func (g *GameWrapper) ResetGame() error {
    // ゲームごとの状態をリセット
    g.Game.Board.Init()                     // ボードの初期化
    g.Game.StartTime = time.Now()           // 開始時刻の設定
    g.Game.TimeLimit = 3 * time.Minute     // タイムリミットの設定
    g.Game.State = "playing"                // 状態のリセット
    g.Game.KeyState = make(map[ebiten.Key]bool) // キー状態のリセット
    g.Game.Current = g.Game.GenerateRandomTetromino() // 現在のテトリミノ
    // Next[0]からNext[5]までを生成
		g.Game.Next = make([]*Tetromino, 6)
		g.Game.Next[0] = g.Game.GenerateRandomTetromino()
		g.Game.Next[1], g.Game.Next[2], g.Game.Next[3], g.Game.Next[4] = g.Game.GenerateUniqueTetrominos()
    g.Game.Score = 0
		g.Game.DrawedEmote = ""

		// wasm.ResetFaceSnapshot()

		return nil
}

func (g *GameWrapper) NewTetromino() error {
	// ゲームごとの状態をリセット
	g.Game.State = "playing"                // 状態のリセット
	g.Game.KeyState = make(map[ebiten.Key]bool) // キー状態のリセット
	g.Game.ShiftTetrominoQueue()                  // テトリミノを生成
	g.Game.Score = 0                       // スコアのリセット
	return nil
}


// レイアウトの設定（ウィンドウのサイズ）
func (g *GameWrapper) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	screenWidth = constants.ScreenWidth  // 画面幅を640に設定
	screenHeight = constants.ScreenHeight // 画面高さを480に設定
	return screenWidth, screenHeight
}

package game

import (
	"square-face-tetris/app/constants"
	"square-face-tetris/app/domain"

	"bytes"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)


type GameWrapper struct {
	Game domain.Game
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
    g.Game.TimeLimit = 10 * time.Second     // タイムリミットの設定
    g.Game.State = "playing"                // 状態のリセット
    g.Game.KeyState = make(map[ebiten.Key]bool) // キー状態のリセット
    g.Game.Current = g.Game.GenerateRandomTetromino() // 現在のテトリミノ
    g.Game.Next = g.Game.GenerateRandomTetromino()    // 次のテトリミノ
    g.Game.Next.Next = g.Game.GenerateRandomTetromino() // 次の次のテトリミノ
    g.Game.Score = 0                       // スコアのリセット
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

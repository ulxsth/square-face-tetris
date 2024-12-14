package game

import (
	"square-face-tetris/app/constants"
	"square-face-tetris/app/domain"

	"bytes"
	"time"

	// "github.com/esimov/pigo/wasm/detector"
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

//FIXME: init と混同するので、名前を変更
// ゲームの初期化（タイマーの設定を追加）
func (g *GameWrapper) Init() error{
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		return err
	}
	mplusFaceSource = s

	g.Game.Board.Init() // Boardの初期化
	g.Game.StartTime = time.Now()           // ゲーム開始時刻を記録
	g.Game.TimeLimit = 10 * time.Second      // タイムリミットを3分に設定
	// g.Game.TimeLimit = 3 * time.Minute      // タイムリミットを3分に設定
	g.Game.State = "playing"               // ゲームオーバー状態を初期化
	g.Game.KeyState = make(map[ebiten.Key]bool) // キー状態をリセット

	g.Game.Current = g.Game.GenerateRandomTetromino()	// 次のテトリミノをランダムに生成
	g.Game.Next = g.Game.GenerateRandomTetromino() 	// 次のテトリミノをランダムに生成
	g.Game.Next.Next = g.Game.GenerateRandomTetromino()	// 次の次のテトリミノをランダムに生成
	
	g.Game.NewTetromino()                   // 最初のテトリミノを生成
	g.Game.Score = 0
	return nil
}


// レイアウトの設定（ウィンドウのサイズ）
func (g *GameWrapper) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	screenWidth = constants.ScreenWidth  // 画面幅を640に設定
	screenHeight = constants.ScreenHeight // 画面高さを480に設定
	return screenWidth, screenHeight
}

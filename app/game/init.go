package game

import (
	"square-face-tetris/app/constants"
	"square-face-tetris/app/types"
	
	"bytes"
	"time"

	// "github.com/esimov/pigo/wasm/detector"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	
	"syscall/js"
)


type GameWrapper struct {
	Game types.Game
}

var (
	video  js.Value
	stream js.Value
	canvas js.Value
	ctx    js.Value
	// det    *detector.Detector
	mplusFaceSource *text.GoTextFaceSource
)


// ゲームの初期化
// NOTE: package の読み込み時に 1度だけ呼び出される
func init() {
	// 検出器の初期化
	// det = detector.NewDetector()
	// if err := det.UnpackCascades(); err != nil {
	// 	log.Fatal(err)
	// }

	doc := js.Global().Get("document")
	video = doc.Call(("createElement"), "video")
	canvas = doc.Call(("createElement"), "canvas")
	video.Set("muted", true)
	video.Set("videoWidth", constants.ScreenWidth)
	video.Set("videoHeight", constants.ScreenHeight)

	// カメラの映像の取得権限をリクエスト
	mediaDevices := js.Global().Get("navigator").Get("mediaDevices")
	promise := mediaDevices.Call("getUserMedia", map[string]interface{}{
		"video": true,
		"audio": false,
	})
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		stream = args[0]
		video.Set("srcObject", stream)
		video.Call("play")
		canvas.Set("width", constants.ScreenWidth)
		canvas.Set("height", constants.ScreenHeight)
		ctx = canvas.Call("getContext", "2d")
		return nil
	}))
}

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
    
	g.newTetromino()                   // 最初のテトリミノを生成
	g.Game.Score = 0
	return nil 
}


// レイアウトの設定（ウィンドウのサイズ）
func (g *GameWrapper) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	screenWidth = constants.ScreenWidth  // 画面幅を640に設定
	screenHeight = constants.ScreenHeight // 画面高さを480に設定
	return screenWidth, screenHeight
}


package types 

import (
	"log"
	"time"
	"github.com/hajimehoshi/ebiten/v2"

	"syscall/js"
	"image/color"
	"square-face-tetris/app/constants"
	"square-face-tetris/app/lib/wasm"
	
)

// ゲームの状態
type Game struct {
	Board   Board // 10x20 のボード
	Current *Tetromino  // 現在のテトリミノ
	LastDrop    time.Time   // 最後にテトリミノが落下した時刻
	DropInterval time.Duration // 落下間隔
	KeyState     map[ebiten.Key]bool // キーの押下状態
}

var (
	video js.Value
	stream js.Value
	canvas js.Value
	ctx js.Value
	det *detector.Detector
)

// ゲームの初期化
// NOTE: package の読み込み時に 1度だけ呼び出される
func init() {
	// 検出器の初期化
	det = detector.NewDetector()
	if err := det.UnpackCascades(); err != nil {
		log.Fatal(err)
	}

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

// ゲームの状態更新
func (g *Game) Update() error {
	if g.Current == nil {
		g.newTetromino()
	}

	// ユーザー入力でテトリミノを操作
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.isValidPosition(g.Current, -1, 0) {
		g.Current.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) && g.isValidPosition(g.Current, 1, 0) {
		g.Current.X += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) && g.isValidPosition(g.Current, 0, 1) {
		g.Current.Y += 1
	}

	// 回転用ボタンの処理（1回の入力で1回だけ回転）
	if ebiten.IsKeyPressed(ebiten.KeyUp) && !g.KeyState[ebiten.KeyUp] {
		g.rotateTetromino()
		g.KeyState[ebiten.KeyUp] = true // 回転ボタンの押下を記録
	}

	// 2秒間隔で落下
	if time.Since(g.LastDrop) > g.DropInterval {
		if g.isValidPosition(g.Current, 0, 1) {
			g.Current.Y += 1
		} else {
			// テトリミノが固定されるべき条件を満たす
			g.lockTetromino()
		}
		g.LastDrop = time.Now() // 落下タイマーをリセット
	}

	// キーが離された場合に状態をリセット（回転だけリセット）
	g.ResetKeyState()

	return nil
}

// キーが離された場合に状態をリセット
func (g *Game) ResetKeyState() {
	for key := range g.KeyState {
		if !ebiten.IsKeyPressed(key) {
			g.KeyState[key] = false // キーが離されたら状態をリセット
		}
	}
}


// ゲームの描画
func (g *Game) Draw(screen *ebiten.Image) {
	// ボードの描画（固定されたブロック）
	for y := 0; y < constants.BoardHeight; y++ {
		for x := 0; x < constants.BoardWidth; x++ {
			if g.Board[y][x] == 1 {
				blockImage := ebiten.NewImage(constants.BlockSize, constants.BlockSize)
				blockImage.Fill(color.RGBA{0, 0, 255, 255}) // 青
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(float64(x*constants.BlockSize), float64(y*constants.BlockSize))
				screen.DrawImage(blockImage, opts)
			}
		}
	}

	// 現在のテトリミノの描画
	if g.Current != nil {
		for y := 0; y < len(g.Current.Shape); y++ {
			for x := 0; x < len(g.Current.Shape[y]); x++ {
				if g.Current.Shape[y][x] == 1 {
					blockImage := ebiten.NewImage(constants.BlockSize, constants.BlockSize)
					blockImage.Fill(g.Current.Color)
					opts := &ebiten.DrawImageOptions{}
					opts.GeoM.Translate(float64((g.Current.X+x)* constants.BlockSize), float64((g.Current.Y+y)*constants.BlockSize))
					screen.DrawImage(blockImage, opts)
				}
			}
		}
	}
}

// レイアウトの設定（ウィンドウのサイズ）
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	screenWidth = 640  // 画面幅を640に設定
	screenHeight = 640 // 画面高さを480に設定
	return screenWidth, screenHeight
}



/*
 * tetromino
 *
 */
// テトリミノを新しく取得
func (g *Game) newTetromino() {
	g.Current = &Tetrominos[0] // 現時点では I 型のテトリミノを設定
	g.Current.X = 3
	g.Current.Y = 0
	g.LastDrop = time.Now() // 新しいテトリミノの生成時にタイマーをリセット
}

// テトリミノの回転処理
func (g *Game) rotateTetromino() {
	// 現在の形状の行数と列数を取得
	rows := len(g.Current.Shape)
	cols := len(g.Current.Shape[0])

	// 回転後の形状を計算
	newShape := make([][]int, cols)
	for i := range newShape {
		newShape[i] = make([]int, rows)
	}

	// 回転処理：90度回転
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			newShape[x][rows-1-y] = g.Current.Shape[y][x]
		}
	}

	// 回転後の形状を適用
	g.Current.Shape = newShape

	// 回転状態を更新
	g.Current.Rotation = (g.Current.Rotation + 90) % 360
}


/*
 * board
 *
 */

// ボードの範囲と重なりをチェック
func (g *Game) isValidPosition(tetromino *Tetromino, offsetX, offsetY int) bool {
	for y := 0; y < len(tetromino.Shape); y++ {
		for x := 0; x < len(tetromino.Shape[y]); x++ {
			if tetromino.Shape[y][x] == 1 {
				newX := tetromino.X + x + offsetX
				newY := tetromino.Y + y + offsetY

				// ボードの範囲外をチェック
				if newX < 0 || newX >= len(g.Board[0]) || newY >= len(g.Board) {
					return false
				}

				// 他のブロックと重なっていないかをチェック
				if newY >= 0 && g.Board[newY][newX] == 1 {
					return false
				}
			}
		}
	}
	return true
}

// ボードにテトリミノを固定
func (g *Game) lockTetromino() {
	for y := 0; y < len(g.Current.Shape); y++ {
			for x := 0; x < len(g.Current.Shape[y]); x++ {
					if g.Current.Shape[y][x] == 1 {
							g.Board[g.Current.Y+y][g.Current.X+x] = 1
					}
			}
	}
	g.Current = nil // 新しいテトリミノを生成
}

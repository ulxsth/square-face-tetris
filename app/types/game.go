package types

import (
	"square-face-tetris/app/constants"

	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"log"
	"time"
	"math/rand"

	// "github.com/esimov/pigo/wasm/detector"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"syscall/js"
)

// ゲームの状態
type Game struct {
	Board        Board               // 10x20 のボード
	Current      *Tetromino          // 現在のテトリミノ
	LastDrop     time.Time           // 最後にテトリミノが落下した時刻
	DropInterval time.Duration       // 落下間隔
	KeyState     map[ebiten.Key]bool // キーの押下状態
	CanvasImage  *ebiten.Image       // canvas から取得した画像を保持するフィールドを追加
	startTime    time.Time           // ゲーム開始時刻
	timeLimit    time.Duration       // タイムリミット
	state        string              // ゲームの状態
	score        int                 // スコア
}

const (
	NORMAL_FONT_SIZE = 24
	BIG_FONT_SIZE    = 48
	X = 20
)

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
func (g *Game) Init() error{
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		return err
	}
	mplusFaceSource = s

	g.Board.Init() // Boardの初期化
	g.startTime = time.Now()           // ゲーム開始時刻を記録
	g.timeLimit = 3 * time.Minute      // タイムリミットを3分に設定
	g.state = "playing"               // ゲームオーバー状態を初期化
	g.KeyState = make(map[ebiten.Key]bool) // キー状態をリセット

	g.newTetromino()                   // 最初のテトリミノを生成
	g.score = 0
	return nil
}

// ゲームの状態更新
func (g *Game) Update() error {
	switch g.state {
	case "playing":
		g.updatePlaying()
	case "showingScore":
		g.updateShowingScore()
	}
	return nil
}


// プレイ中の状態を更新
func (g *Game) updatePlaying() {
	// タイムリミットを超えている場合はスコア画面へ遷移
	if time.Since(g.startTime) >= g.timeLimit {
		g.state = "showingScore"
		return
	}

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

	// video の映像を canvas に移す
	if constants.ENABLED_CAMERA_PREVIEw && ctx.Truthy() {
		ctx.Call("drawImage", video, 0, 0, constants.ScreenWidth, constants.ScreenHeight)
		// canvas 経由で画面を base64 形式で取得
		b64 := canvas.Call("toDataURL", "image/png").String()

		// image.Image にデコード
		dec, err := base64.StdEncoding.DecodeString(b64[22:])
		if err != nil {
			log.Fatal(err)
		}
		img, _, err := image.Decode(bytes.NewReader(dec))
		if err != nil {
			log.Fatal(err)
		}

		// ebiten.Image にして保持
		g.CanvasImage = ebiten.NewImageFromImage(img)
	}
}

// キーが離された場合に状態をリセット
func (g *Game) ResetKeyState() {
	for key := range g.KeyState {
		if !ebiten.IsKeyPressed(key) {
			g.KeyState[key] = false // キーが離されたら状態をリセット
		}
	}
}

func (g *Game) updateShowingScore() {
	// スコア画面ではスペースキーを押すと終了
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		err := g.Init() // ゲームを初期化
		if err != nil {
				log.Fatalf("Failed to initialize the game: %v", err)
		}
	}
}

// ゲームの描画
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
		case "playing":
			g.drawPlaying(screen)
		case "showingScore":
			g.drawScore(screen)
	}
}


// プレイ中の描画
func (g *Game) drawPlaying(screen *ebiten.Image) {
	// 残り時間を計算
	remainingTime := g.timeLimit - time.Since(g.startTime)
	if remainingTime < 0 {
		remainingTime = 0
	}

	// 秒数に変換
	totalSeconds := remainingTime.Seconds()

	// 時間、分、秒を計算
	minutes := int(totalSeconds) / 60
	seconds := int(totalSeconds) % 60
	hundredths := int((totalSeconds - float64(int(totalSeconds))) * 100)

	// タイマーの表示
	timerText := fmt.Sprintf("%02d:%02d.%02d", minutes, seconds, hundredths)
	op1 := &text.DrawOptions{}
	op1.GeoM.Translate(X, 20)
	op1.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, timerText, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   NORMAL_FONT_SIZE,
	}, op1)


	// スコアの表示
	scoreText := fmt.Sprintf("Score: %d", g.score)
	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(X, 40)
	op2.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, scoreText,&text.GoTextFace{
		Source: mplusFaceSource,
		Size:   NORMAL_FONT_SIZE,
	}, op2)

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

		if g.CanvasImage != nil {
			// 保持している ebiten.Image を右上に描画
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Scale(0.25, 0.2) // サイズを固定
			opts.GeoM.Translate(float64(constants.ScreenWidth-g.CanvasImage.Bounds().Dx()/4), 0)
			screen.DrawImage(g.CanvasImage, opts)
		}
		ebitenutil.DebugPrint(screen, fmt.Sprintf("%f", ebiten.ActualFPS()))
	}

	// 現在のテトリミノの描画
	if g.Current != nil {
		for y := 0; y < len(g.Current.Shape); y++ {
			for x := 0; x < len(g.Current.Shape[y]); x++ {
				if g.Current.Shape[y][x] == 1 {
					blockImage := ebiten.NewImage(constants.BlockSize, constants.BlockSize)
					blockImage.Fill(g.Current.Color)
					opts := &ebiten.DrawImageOptions{}
					opts.GeoM.Translate(float64((g.Current.X+x)*constants.BlockSize), float64((g.Current.Y+y)*constants.BlockSize))
					screen.DrawImage(blockImage, opts)
				}
			}
		}
	}
}

// スコア画面の描画
func (g *Game) drawScore(screen *ebiten.Image) {
	// 背景を塗りつぶす
	screen.Fill(color.Black)

	// スコアを表示
	scoreText := fmt.Sprintf("Final Score: %d", g.score)
	op3 := &text.DrawOptions{}
	op3.GeoM.Translate(X, 60)
	op3.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, scoreText,&text.GoTextFace{
		Source: mplusFaceSource,
		Size:   NORMAL_FONT_SIZE,
	}, op3)

	// リスタートの指示を表示
	restartText := "Press SPACE to Restart"
	op4 := &text.DrawOptions{}
	op4.GeoM.Translate(X, 80)
	op4.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, restartText,&text.GoTextFace{
		Source: mplusFaceSource,
		Size:   NORMAL_FONT_SIZE,
	}, op4)

	// ゲーム終了のメッセージ
	exitText := "Thank you for playing!"
	op5 := &text.DrawOptions{}
	op5.GeoM.Translate(X, 100)
	op5.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, exitText,&text.GoTextFace{
		Source: mplusFaceSource,
		Size:   NORMAL_FONT_SIZE,
	}, op5)
}

// レイアウトの設定（ウィンドウのサイズ）
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	screenWidth = constants.ScreenWidth  // 画面幅を640に設定
	screenHeight = constants.ScreenHeight // 画面高さを480に設定
	return screenWidth, screenHeight
}

/*
 * tetromino
 *
 */
// テトリミノを新しく取得
func (g *Game) newTetromino() {
	randomIndex := rand.Intn(len(Tetrominos)) // テトロミノのリストからランダムに選択
	g.Current = &Tetrominos[randomIndex] // 現時点では I 型のテトリミノを設定
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
// ボードにテトリミノを固定
func (g *Game) lockTetromino() {
	for y := 0; y < len(g.Current.Shape); y++ {
		for x := 0; x < len(g.Current.Shape[y]); x++ {
			if g.Current.Shape[y][x] == 1 {
				g.Board[g.Current.Y+y][g.Current.X+x] = 1
			}
		}
	}

	// 横一列が揃っているか確認
	g.clearFullRows()

	// 新しいテトリミノを生成
	g.Current = nil
}

// 横一列が揃った行を削除し、スコアを加算
func (g *Game) clearFullRows() {
	clearedRows := 0

	// 上から下へループ
	for y := len(g.Board) - 1; y >= 0; y-- {
		full := true
		for x := 0; x < len(g.Board[y]); x++ {
			if g.Board[y][x] == 0 {
				full = false
				break
			}
		}

		// 横一列が揃っている場合
		if full {
			clearedRows++

			// 上の行を下にずらす
			for yy := y; yy > 0; yy-- {
				g.Board[yy] = g.Board[yy-1]
			}

			// 一番上の行を初期化
			g.Board[0] = make([]int, len(g.Board[0]))

			// 現在の行を再チェック（行をずらしたため）
			y++
		}
	}

	// スコアを加算（1行100点、2行300点、3行600点、4行1000点）
	if clearedRows > 0 {
		g.score += clearedRows * (clearedRows + 1) / 2 * 100
	}
}


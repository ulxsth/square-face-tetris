package game

import (
	"fmt"
	"image/color"

	"math/rand"
	"square-face-tetris/app/constants"
	"square-face-tetris/app/domain/wasm"
	"time"
)

// テトリミノの定義
type Tetromino struct {
	X, Y      int         // テトリミノの位置
	Color     color.Color // テトリミノの色
	Shape     [][]int     // テトリミノの形状（回転可能）
	Rotation  int         // 回転状態（0, 90, 180, 270度）
	Next      *Tetromino
}

// 各テトリミノの形状を定義
var Tetrominos = []Tetromino{
	{
		Color: color.RGBA{255, 0, 0, 255}, // 赤 - I
		Shape: [][]int{
			{1, 1, 1, 1}, // 横一列
		},
	},
	{
		Color: color.RGBA{0, 255, 0, 255}, // 緑 - O
		Shape: [][]int{
			{1, 1},
			{1, 1}, // 正方形
		},
	},
	{
		Color: color.RGBA{0, 0, 255, 255}, // 青 - T
		Shape: [][]int{
			{0, 1, 0},
			{1, 1, 1}, // T字型
		},
	},
	{
		Color: color.RGBA{255, 165, 0, 255}, // オレンジ - L
		Shape: [][]int{
			{1, 0},
			{1, 0},
			{1, 1}, // L字型
		},
	},
	{
		Color: color.RGBA{0, 255, 255, 255}, // 水色 - J
		Shape: [][]int{
			{0, 1},
			{0, 1},
			{1, 1}, // 逆L字型
		},
	},
	{
		Color: color.RGBA{255, 255, 0, 255}, // 黄 - S
		Shape: [][]int{
			{0, 1, 1},
			{1, 1, 0}, // S字型
		},
	},
	{
		Color: color.RGBA{128, 0, 128, 255}, // 紫 - Z
		Shape: [][]int{
			{1, 1, 0},
			{0, 1, 1}, // Z字型
		},
	},
}

// テトリミノを新しく取得
func (g *Game) ShiftTetrominoQueue() {
	// 現在のテトリミノをNext[0]として設定
	g.Current = g.Next[0]

	// Trueの感情を取得
	emotionIndexes := wasm.Face.GetEmotionIndexes()

  // Trueの感情から抽選
	fmt.Printf("emotionIndexes: %v\n", emotionIndexes)
	drawedIndex := g.drawingEmotionFromFlags(emotionIndexes)
	g.DrawedEmote = wasm.Face.GetEmotionByIndex(drawedIndex)
	g.Next[0] = g.Next[drawedIndex]

	// 次の次のテトリミノを生成
	g.Next[1], g.Next[2], g.Next[3], g.Next[4] = g.GenerateUniqueTetrominos()
	// 現在のテトリミノの位置を初期化
	g.Current.X = 3
	g.Current.Y = 0

	// ドロップのタイマーをリセット
	g.LastDrop = time.Now()
}

// drawingEmotionFromFlags は、emotionIndexes の長さに基づいて
// 最大値100をその長さで分割し、ランダムな確率に基づいてインデックスを返す
// emotionIndexes が空の場合、1～3の間でランダムなインデックスを生成
func (g *Game) drawingEmotionFromFlags(emotionIndexes []int) int {
	// emotionIndexes が空かどうかをチェック
	if len(emotionIndexes) == 0 {
		// 空の場合は 1～3 の間でランダムなインデックスを生成
		return rand.Intn(4) + 1 // 1, 2, 3, 4 のいずれかを返す
	}

	// emotionIndexes の長さを取得
	numEmotions := len(emotionIndexes)

	// 最大値100を分割するために、各インデックスが占める確率を計算
	probabilityRanges := make([]int, numEmotions)
	rangeSize := 100 / numEmotions

	// 確率範囲を設定
	for i := 0; i < numEmotions; i++ {
		probabilityRanges[i] = (i + 1) * rangeSize
	}

	// ランダムな数値を生成（0から99まで）
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomValue := rand.Intn(100)

	// ランダムな数値がどの範囲に属するか調べる
	for i := 0; i < numEmotions; i++ {
		if randomValue < probabilityRanges[i] {
			return emotionIndexes[i] // 該当するインデックスを返す
		}
	}

	// 万が一、どれにも該当しなかった場合は最後のインデックスを返す
	return emotionIndexes[numEmotions-1]
}

// ランダムにテトリミノを生成するヘルパー関数
func (g *Game) GenerateRandomTetromino() *Tetromino {
	randomIndex := rand.Intn(len(Tetrominos)) // テトリミノのリストからランダムに選択

	// 選択したテトリミノを新しくインスタンス化して返す
	newTetromino := Tetromino{
			Color:    Tetrominos[randomIndex].Color,
			Shape:    append([][]int{}, Tetrominos[randomIndex].Shape...), // Shapeを新しくコピー
			Rotation: 0, // 初期回転状態
	}

	return &newTetromino
}

func (g *Game) GenerateUniqueTetrominos() (*Tetromino, *Tetromino, *Tetromino, *Tetromino) {
	// シャッフルアルゴリズムを使用して重複を防ぐ
	indexes := rand.Perm(len(Tetrominos)) // 0からlen(Tetrominos)-1までのランダム順列を生成

	// 3つのテトリミノを生成
	tetromino1 := Tetromino{
		Color:    Tetrominos[indexes[0]].Color,
		Shape:    append([][]int{}, Tetrominos[indexes[0]].Shape...), // Shapeを新しくコピー
		Rotation: 0,
	}
	tetromino2 := Tetromino{
		Color:    Tetrominos[indexes[1]].Color,
		Shape:    append([][]int{}, Tetrominos[indexes[1]].Shape...), // Shapeを新しくコピー
		Rotation: 0,
	}
	tetromino3 := Tetromino{
		Color:    Tetrominos[indexes[2]].Color,
		Shape:    append([][]int{}, Tetrominos[indexes[2]].Shape...), // Shapeを新しくコピー
		Rotation: 0,
	}
	tetromino4 := Tetromino{
		Color:    Tetrominos[indexes[3]].Color,
		Shape:    append([][]int{}, Tetrominos[indexes[3]].Shape...), // Shapeを新しくコピー
		Rotation: 0,
	}

	return &tetromino1, &tetromino2, &tetromino3, &tetromino4
}

// テトリミノの回転処理
func (g *Game) RotateTetromino() {
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

func (g *Game) IsOutOfBounds() string {
	current := g.Current
	rows := len(current.Shape)
	cols := len(current.Shape[0])

	// ブロックの形状の各セルを確認
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if current.Shape[y][x] != 0 { // 非空ブロックのみ確認
				boardX := current.X + x
				boardY := current.Y + y

				// 左側に出ている場合
				if boardX < 0 {
					return "left"
				}

				// 右側に出ている場合
				if boardX >= constants.BoardWidth {
					return "right"
				}

				// 下側に出ている場合
				if boardY >= constants.BoardHeight {
					return "bottom"
				}
			}
		}
	}
	// 範囲外ではない場合は空文字を返す
	return ""
}

func (g *Game) IsOverlapping() bool {
	current := g.Current
	rows := len(current.Shape)
	cols := len(current.Shape[0])

	// ブロックの形状の各セルを確認
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if current.Shape[y][x] != 0 { // 非空ブロックのみ確認
				boardX := current.X + x
				boardY := current.Y + y

				// ボード上に位置している場合のみ重なりを確認
				if boardY >= 0 && boardY < constants.BoardHeight &&
				   boardX >= 0 && boardX < constants.BoardWidth {
					if g.Board[boardY][boardX] != 0 {
						return true // 他のブロックと重なっている
					}
				}
			}
		}
	}
	return false // 重なりがない
}


// 横一列が揃った行を削除し、スコアを加算
func (g *Game) ClearFullRows() {
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
		g.Score += clearedRows * (clearedRows + 1) / 2 * 100
	}
}



// ボードの範囲と重なりをチェック
func (g *Game) IsValidPosition(tetromino *Tetromino, offsetX, offsetY int) bool {
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
func (g *Game) LockTetromino() {
	for y := 0; y < len(g.Current.Shape); y++ {
		for x := 0; x < len(g.Current.Shape[y]); x++ {
			if g.Current.Shape[y][x] == 1 {
				g.Board[g.Current.Y+y][g.Current.X+x] = 1
			}
		}
	}

	// 横一列が揃っているか確認
	g.ClearFullRows()

	// 新しいテトリミノを生成
	g.Current = nil
}



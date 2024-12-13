package game

import (
	"square-face-tetris/app/domain"
	"square-face-tetris/app/constants"

	"math/rand"
	"time"

	// "github.com/esimov/pigo/wasm/detector"

)

type TetrominoWrapper struct {
	Tetromino domain.Tetromino
}

// テトリミノを新しく取得
func (g *GameWrapper) newTetromino() {
	randomIndex := rand.Intn(len(domain.Tetrominos)) // テトロミノのリストからランダムに選択
	g.Game.Current = &domain.Tetrominos[randomIndex] // 現時点では I 型のテトリミノを設定
	g.Game.Current.X = 3
	g.Game.Current.Y = 0
	g.Game.LastDrop = time.Now() // 新しいテトリミノの生成時にタイマーをリセット
}

// テトリミノの回転処理
func (g *GameWrapper) rotateTetromino() {
	// 現在の形状の行数と列数を取得
	rows := len(g.Game.Current.Shape)
	cols := len(g.Game.Current.Shape[0])

	// 回転後の形状を計算
	newShape := make([][]int, cols)
	for i := range newShape {
		newShape[i] = make([]int, rows)
	}

	// 回転処理：90度回転
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			newShape[x][rows-1-y] = g.Game.Current.Shape[y][x]
		}
	}

	// 回転後の形状を適用
	g.Game.Current.Shape = newShape

	// 回転状態を更新
	g.Game.Current.Rotation = (g.Game.Current.Rotation + 90) % 360
}

func (g *GameWrapper) isOutOfBounds() string {
	current := g.Game.Current
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

func (g *GameWrapper) isOverlapping() bool {
	current := g.Game.Current
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
					if g.Game.Board[boardY][boardX] != 0 {
						return true // 他のブロックと重なっている
					}
				}
			}
		}
	}
	return false // 重なりがない
}


// 横一列が揃った行を削除し、スコアを加算
func (g *GameWrapper) clearFullRows() {
	clearedRows := 0

	// 上から下へループ
	for y := len(g.Game.Board) - 1; y >= 0; y-- {
		full := true
		for x := 0; x < len(g.Game.Board[y]); x++ {
			if g.Game.Board[y][x] == 0 {
				full = false
				break
			}
		}

		// 横一列が揃っている場合
		if full {
			clearedRows++

			// 上の行を下にずらす
			for yy := y; yy > 0; yy-- {
				g.Game.Board[yy] = g.Game.Board[yy-1]
			}

			// 一番上の行を初期化
			g.Game.Board[0] = make([]int, len(g.Game.Board[0]))

			// 現在の行を再チェック（行をずらしたため）
			y++
		}
	}

	// スコアを加算（1行100点、2行300点、3行600点、4行1000点）
	if clearedRows > 0 {
		g.Game.Score += clearedRows * (clearedRows + 1) / 2 * 100
	}
}



// ボードの範囲と重なりをチェック
func (g *GameWrapper) isValidPosition(tetromino *domain.Tetromino, offsetX, offsetY int) bool {
	for y := 0; y < len(tetromino.Shape); y++ {
		for x := 0; x < len(tetromino.Shape[y]); x++ {
			if tetromino.Shape[y][x] == 1 {
				newX := tetromino.X + x + offsetX
				newY := tetromino.Y + y + offsetY

				// ボードの範囲外をチェック
				if newX < 0 || newX >= len(g.Game.Board[0]) || newY >= len(g.Game.Board) {
					return false
				}

				// 他のブロックと重なっていないかをチェック
				if newY >= 0 && g.Game.Board[newY][newX] == 1 {
					return false
				}
			}
		}
	}
	return true
}

// ボードにテトリミノを固定
// ボードにテトリミノを固定
func (g *GameWrapper) lockTetromino() {
	for y := 0; y < len(g.Game.Current.Shape); y++ {
		for x := 0; x < len(g.Game.Current.Shape[y]); x++ {
			if g.Game.Current.Shape[y][x] == 1 {
				g.Game.Board[g.Game.Current.Y+y][g.Game.Current.X+x] = 1
			}
		}
	}

	// 横一列が揃っているか確認
	g.clearFullRows()

	// 新しいテトリミノを生成
	g.Game.Current = nil
}



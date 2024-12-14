package game

import (
	"square-face-tetris/app/constants"
	"square-face-tetris/app/domain/wasm"

	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	normalFontSize = 24
	bigFontSize    = 48
	x = 20
)


// ゲームの描画
func (g *GameWrapper) Draw(screen *ebiten.Image) {
	switch g.Game.State {
		case "playing":
			g.drawPlaying(screen)
		case "showingScore":
			g.drawScore(screen)
	}
}


// プレイ中の描画
func (g *GameWrapper) drawPlaying(screen *ebiten.Image) {
	// 背景を塗りつぶす（紺色）
	screen.Fill(color.RGBA{0, 0, 64, 255}) // 紺色

	// 残り時間を計算
	remainingTime := g.Game.TimeLimit - time.Since(g.Game.StartTime)
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
	op1.GeoM.Translate(x, 20)
	op1.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, timerText, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, op1)

	// スコアの表示
	scoreText := fmt.Sprintf("Score: %d", g.Game.Score)
	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(x, 40)
	op2.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, scoreText,&text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, op2)

	// ボードの描画（固定されたブロック）
	for y := 0; y < constants.BoardHeight; y++ {
		for x := 0; x < constants.BoardWidth; x++ {
			if g.Game.Board[y][x] == 1 {
				blockImage := ebiten.NewImage(constants.BlockSize, constants.BlockSize)
				blockImage.Fill(color.RGBA{0, 0, 255, 255}) // 青
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(float64(x*constants.BlockSize), float64(y*constants.BlockSize))
				screen.DrawImage(blockImage, opts)
			}
		}

		wasm.DrawCameraPrev(screen)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("%f", ebiten.ActualFPS()))
	}

	g.DrawNextTetromino(screen)
	g.DrawAfterNextTetromino(screen)

	// 現在のテトリミノの描画
	if g.Game.Current != nil {
		for y := 0; y < len(g.Game.Current.Shape); y++ {
			for x := 0; x < len(g.Game.Current.Shape[y]); x++ {
				if g.Game.Current.Shape[y][x] == 1 {
					blockImage := ebiten.NewImage(constants.BlockSize, constants.BlockSize)
					blockImage.Fill(g.Game.Current.Color)
					opts := &ebiten.DrawImageOptions{}
					opts.GeoM.Translate(float64((g.Game.Current.X+x)*constants.BlockSize), float64((g.Game.Current.Y+y)*constants.BlockSize))
					screen.DrawImage(blockImage, opts)
				}
			}
		}
	}
}

// 次のテトロミノの描画
func (g *GameWrapper) DrawNextTetromino(screen *ebiten.Image) {
	// 「Next」のラベルを描画
	nextLabel := "Next:"
	op := &text.DrawOptions{}
	op.GeoM.Translate(constants.BoardWidth*constants.BlockSize+10, 120)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, nextLabel, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, op)

	// 次のテトロミノの描画
	if g.Game.Next != nil {
		for y := 0; y < len(g.Game.Next.Shape); y++ {
			for x := 0; x < len(g.Game.Next.Shape[y]); x++ {
				if g.Game.Next.Shape[y][x] == 1 {
					blockImage := ebiten.NewImage(constants.BlockSize, constants.BlockSize)
					blockImage.Fill(g.Game.Next.Color) // 次のテトロミノの色
					opts := &ebiten.DrawImageOptions{}
					opts.GeoM.Translate(
						float64(constants.BoardWidth*constants.BlockSize+10+(x*constants.BlockSize)),
						float64(150+(y*constants.BlockSize)),
					)
					screen.DrawImage(blockImage, opts)
				}
			}
		}
	}
}

func (g *GameWrapper) DrawAfterNextTetromino(screen *ebiten.Image) {
	// 「Next」のラベルを描画
	nextLabel := "After Next: "
	op := &text.DrawOptions{}
	op.GeoM.Translate(constants.BoardWidth*constants.BlockSize+10, 200)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, nextLabel , &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, op)

		// 次のテトロミノの描画
		if g.Game.Next.Next != nil {
			for y := 0; y < len(g.Game.Next.Next.Shape); y++ {
				for x := 0; x < len(g.Game.Next.Next.Shape[y]); x++ {
					if g.Game.Next.Next.Shape[y][x] == 1 {
						blockImage := ebiten.NewImage(constants.BlockSize, constants.BlockSize)
						blockImage.Fill(g.Game.Next.Next.Color) // 次のテトロミノの色
						opts := &ebiten.DrawImageOptions{}
						opts.GeoM.Translate(
							float64(constants.BoardWidth*constants.BlockSize+10+(x*constants.BlockSize)),
							float64(230+(y*constants.BlockSize)),
						)
						screen.DrawImage(blockImage, opts)
					}
				}
			}
		}
}

// スコア画面の描画
func (g *GameWrapper) drawScore(screen *ebiten.Image) {
	// 背景を塗りつぶす
	screen.Fill(color.Black)

	// スコアを表示
	scoreText := fmt.Sprintf("Final Score: %d", g.Game.Score)
	op3 := &text.DrawOptions{}
	op3.GeoM.Translate(x, 60)
	op3.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, scoreText,&text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, op3)

	// リスタートの指示を表示
	restartText := "Press SPACE to Restart"
	op4 := &text.DrawOptions{}
	op4.GeoM.Translate(x, 100)
	op4.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, restartText,&text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, op4)

	// ゲーム終了のメッセージ
	exitText := "Thank you for playing!"
	op5 := &text.DrawOptions{}
	op5.GeoM.Translate(x, 140)
	op5.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, exitText,&text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, op5)
}


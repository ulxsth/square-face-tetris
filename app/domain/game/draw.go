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
	x              = 20
)

// ゲームの描画
func (g *GameWrapper) Draw(screen *ebiten.Image) {
	switch g.Game.State {
	case "start":
		g.drawStart(screen)
	case "playing":
		g.drawPlaying(screen)
	case "showingScore":
		g.drawScore(screen)
	}

	wasm.DrawCameraPrev(screen)
}

var InstructionText = []string{
	"操作方法",    // 1行目
	"移動: ←↓→", // 2行目
	"回転: ↑",   // 3行目
}

var EmoText = []string{
	"ANGRY",
	"SURPRISED",
	"SUS",
	"UNKNOWN",
}

// スコア画面の描画
func (g *GameWrapper) drawStart(screen *ebiten.Image) {
	// 背景を塗りつぶす
	screen.Fill(color.Black)

	// スコアを表示
	TitleText := "顔テトリス"
	op3 := &text.DrawOptions{}
	op3.GeoM.Translate(x, 60)
	op3.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, TitleText, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, op3)

	// リスタートの指示を表示
	startText := "スペースキーを押してスタート"
	op4 := &text.DrawOptions{}
	op4.GeoM.Translate(x, 100)
	op4.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, startText, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, op4)

	// ゲーム終了のメッセージ
	op5 := &text.DrawOptions{}
	op5.GeoM.Translate(x, 140)
	op5.ColorScale.ScaleWithColor(color.White)
	for _, line := range InstructionText {
		op5.GeoM.Translate(0, float64(constants.BlockSize)) // 各行の縦位置をずらす
		text.Draw(screen, line, &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   normalFontSize,
		}, op5)
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

	op7 := &text.DrawOptions{}
	op7.GeoM.Translate(constants.BoardWidth*constants.BlockSize+ constants.BlockSize*5, 140)
	op7.ColorScale.ScaleWithColor(color.White)
	for _, line := range EmoText {
		op7.GeoM.Translate(0, float64(constants.BlockSize)*3) // 各行の縦位置をずらす
		text.Draw(screen, line, &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   normalFontSize,
		}, op7)
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
	text.Draw(screen, scoreText, &text.GoTextFace{
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
	emotionText := fmt.Sprintf("%s", g.Game.DrawedEmote)
	op6 := &text.DrawOptions{}
	op6.GeoM.Translate(constants.BoardWidth*constants.BlockSize + constants.BlockSize, 32)
	op6.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, emotionText, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, op6)


	// 次のテトロミノの描画
	if g.Game.Next[0] != nil {
		for y := 0; y < len(g.Game.Next[0].Shape); y++ {
			for x := 0; x < len(g.Game.Next[0].Shape[y]); x++ {
				if g.Game.Next[0].Shape[y][x] == 1 {
					blockImage := ebiten.NewImage(constants.BlockSize, constants.BlockSize)
					blockImage.Fill(g.Game.Next[0].Color) // 次のテトロミノの色
					opts := &ebiten.DrawImageOptions{}
					opts.GeoM.Translate(
						float64(constants.BoardWidth*constants.BlockSize+ constants.BlockSize +(x*constants.BlockSize)),
						float64(64+(y*constants.BlockSize)),
					)
					screen.DrawImage(blockImage, opts)
				}
			}
		}
	}
}

// 次の次のテトロミノを描画
func (g *GameWrapper) DrawAfterNextTetromino(screen *ebiten.Image) {

	// 1から4まで次のテトロミノを描画
	for i := 1; i <= 5 && i <= len(g.Game.Next); i++ {
		if g.Game.Next[i] != nil {
			for y := 0; y < len(g.Game.Next[i].Shape); y++ {
				for x := 0; x < len(g.Game.Next[i].Shape[y]); x++ {
					if g.Game.Next[i].Shape[y][x] == 1 {
						blockImage := ebiten.NewImage(constants.BlockSize, constants.BlockSize)
						blockImage.Fill(g.Game.Next[i].Color) // 次のテトロミノの色
						opts := &ebiten.DrawImageOptions{}
						opts.GeoM.Translate(
							float64(constants.BoardWidth*constants.BlockSize+constants.BlockSize+(x*constants.BlockSize)),
							float64(64+(i)*128+(y*constants.BlockSize)), // Y座標をiに基づいて調整
						)
						screen.DrawImage(blockImage, opts)
					}
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
	scoreText := fmt.Sprintf("総スコア: %d", g.Game.Score)
	op3 := &text.DrawOptions{}
	op3.GeoM.Translate(x, 60)
	op3.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, scoreText, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, op3)

	// リスタートの指示を表示
	restartText := "スペースを押して再スタート"
	op4 := &text.DrawOptions{}
	op4.GeoM.Translate(x, 100)
	op4.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, restartText, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, op4)

	// ゲーム終了のメッセージ
	exitText := "Nice, Face!"
	op5 := &text.DrawOptions{}
	op5.GeoM.Translate(x, 140)
	op5.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, exitText, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   normalFontSize,
	}, op5)
}

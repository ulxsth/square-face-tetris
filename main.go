//go:build js && wasm
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"log"
	"syscall/js"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	video  js.Value
	stream js.Value
	canvas js.Value
	ctx    js.Value
)

const (
	ScreenWidth  = 320
	ScreenHeight = 240
)

func init() {
	// js.Global() は JS の window に相当する
	doc := js.Global().Get("document")
	canvas = doc.Call("createElement", "canvas")
	
	video = doc.Call("createElement", "video")
	video.Set("autoplay", true)
	video.Set("muted", true)
	video.Set("videoWidth", ScreenWidth)
	video.Set("videoHeight", ScreenHeight)

	// Web API からメディアデバイスの API にアクセスする
	// mediaDevices: https://developer.mozilla.org/ja/docs/Web/API/MediaDevices
	mediaDevices := js.Global().Get("navigator").Get("mediaDevices")

	// ユーザーに対してカメラのアクセスをリクエスト
	promise := mediaDevices.Call("getUserMedia", map[string]interface{}{
		"video": true,
		"audio": false,
	})

	// カメラの映像を video に流すように設定
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		stream = args[0]
		video.Set("srcObject", stream)
		video.Call("play")
		canvas.Set("width", ScreenWidth)
		canvas.Set("height", ScreenHeight)
		ctx = canvas.Call("getContext", "2d")
		return nil
	}))

}

type Game struct {}

func newGame() *Game {
	return &Game{}
}

func (g *Game) Update() error {
	if !ctx.Truthy() {
		return nil
	}
	// video の映像を canvas に移す
	ctx.Call("drawImage", video, 0, 0, ScreenWidth, ScreenHeight)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if stream.Truthy() && ctx.Truthy() {
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

		// ebiten.Image にしてゲーム画面に描画
		screen.DrawImage(ebiten.NewImageFromImage(img), nil)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%f", ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth * 2, ScreenHeight * 2
}

func main() {
	ebiten.SetWindowSize(ScreenWidth*2, ScreenHeight*2)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
package wasm

import (
	"bytes"
	"encoding/base64"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"square-face-tetris/app/constants"
	"syscall/js"
)

var (
	video  js.Value
	stream js.Value
	canvas js.Value
	ctx    js.Value
	// det    *detector.Detector
	CanvasImage *ebiten.Image
)

func InitCamera() {
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

func UpdateCamera() {
	// video の映像を canvas に移す
	if ctx.Truthy() {
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
		CanvasImage = ebiten.NewImageFromImage(img)
	}
}

func DrawCameraPrev(screen *ebiten.Image) {
	if CanvasImage != nil {
		// 保持している ebiten.Image を右上に描画
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(0.25, 0.2) // サイズを固定
		opts.GeoM.Translate(float64(constants.ScreenWidth-CanvasImage.Bounds().Dx()/4), 0)
		screen.DrawImage(CanvasImage, opts)
	}
}

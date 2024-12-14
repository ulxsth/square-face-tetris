package wasm

import (
	"bytes"
	"encoding/base64"
	"fmt"

	// "fmt"
	"image"
	"log"
	"time"

	"square-face-tetris/app/constants"
	"syscall/js"

	"github.com/esimov/pigo/wasm/detector"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	video       js.Value
	stream      js.Value
	canvas      js.Value
	ctx         js.Value
	cascadeFile []byte
	det         *detector.Detector

	CanvasImage    *ebiten.Image
	lastUpdateTime time.Time
	updateInterval = time.Second / constants.CAMERA_PREVIEW_FPS
)

func InitCamera() {
	// 分析器の初期化
	det = detector.NewDetector()
	err := det.UnpackCascades()
	if err != nil {
		log.Fatal(err)
	}

	// DOM 要素の取得
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
	// 一定間隔で映像を取得
	if time.Since(lastUpdateTime) < updateInterval {
		return
	}
	lastUpdateTime = time.Now()

	// video の映像を canvas に移す
	if !constants.IS_CAMERA || !ctx.Truthy() {
		return
	}

	ctx.Call("drawImage", video, 0, 0, constants.ScreenWidth, constants.ScreenHeight)
	// canvas 経由で画面を base64 形式で取得
	b64 := canvas.Call("toDataURL", "image/png").String()
	ctx.Call("drawImage", video, 0, 0)
	rgba := ctx.Call("getImageData", 0, 0, constants.ScreenWidth, constants.ScreenHeight, map[string]interface{}{
		"willReadFrequently": true,
	}).Get("data")
	fmt.Println(rgba)

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

func DrawCameraPrev(screen *ebiten.Image) {
	if !constants.IS_CAMERA_PREVIEW || CanvasImage == nil {
		return
	}

	// 保持している ebiten.Image を右上に描画
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(0.25, 0.2) // サイズを固定
	opts.GeoM.Translate(float64(constants.ScreenWidth-CanvasImage.Bounds().Dx()/4), 0)
	screen.DrawImage(CanvasImage, opts)
}

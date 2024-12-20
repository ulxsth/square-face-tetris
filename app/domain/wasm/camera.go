package wasm

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"image"
	"log"
	"math"
	"time"

	"square-face-tetris/app/constants"
	"square-face-tetris/app/domain"
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

	cameraWidth   int
	cameraHeight  int
	aspectRatio   float64
	previewWidth  float64
	previewHeight float64

	Face domain.Face
	IsFaceInited bool

	lastEmotionAnalysisTime time.Time
	emotionAnalysisInterval = time.Second / constants.EMOTION_ANALYSIS_FPS

	MoveLeft  bool
	MoveRight bool
	MoveUp    bool
	MoveDown  bool
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
		video.Call("addEventListener", "loadedmetadata", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			cameraWidth = video.Get("videoWidth").Int()
			cameraHeight = video.Get("videoHeight").Int()
			canvas.Set("width", cameraWidth)
			canvas.Set("height", cameraHeight)
			ctx = canvas.Call("getContext", "2d")

			// アスペクト比を計算
			aspectRatio = float64(cameraHeight) / float64(cameraWidth)
			previewWidth = float64(cameraWidth) * 0.25
			previewHeight = previewWidth * aspectRatio

			return nil
		}))
		return nil
	}))
}

func UpdateCamera() {
	if !ctx.Truthy() || time.Since(lastUpdateTime) < updateInterval {
		return
	}
	lastUpdateTime = time.Now()

	ctx.Call("clearRect", 0, 0, cameraWidth, cameraHeight)

	// video の映像を canvas に移す
	if !constants.IS_CAMERA {
		return
	}

	ctx.Call("drawImage", video, 0, 0, cameraWidth, cameraHeight)

	rgba := ctx.Call("getImageData", 0, 0, cameraWidth, cameraHeight, map[string]interface{}{
		"willReadFrequently": true,
	}).Get("data")

	// 画像を分析
	var data = make([]byte, cameraWidth*cameraHeight*4)
	uint8Arr := js.Global().Get("Uint8Array").New(rgba)
	js.CopyBytesToGo(data, uint8Arr)
	pixels := rgbaToGrayscale(data)

	// 表情分析の頻度を制御
	if time.Since(lastEmotionAnalysisTime) >= emotionAnalysisInterval {
		lastEmotionAnalysisTime = time.Now()
		// det.DetectFaces は画像データを受け取り、以下のデータを返す
		// [row, col, scale, q]
		// row, col: 顔の中心座標
		// scale: 顔のスケール
		// q: 顔であることの信頼度
		res := det.DetectFaces(pixels, cameraHeight, cameraWidth)
		if len(res) > 0 {
			DrawFaceRect(res)

			// 両目の位置を取得
			leftEye := det.DetectLeftPupil(res[0])
			rightEye := det.DetectRightPupil(res[0])

			// 顔のランドマークを取得
			landmarks := det.DetectLandmarkPoints(leftEye, rightEye)
			DrawLandmarkPoints(landmarks)

			// 顔の情報が未設定の場合、新しい顔を作成
			if !IsFaceInited {
				Face = domain.NewFace(landmarks)
				IsFaceInited = true
			}

			// 顔の情報を更新
			choices := []int{
				constants.SMILE,
				constants.ANGRY,
				constants.SURPRISED,
				constants.SUS,
			}
			Face.Update(landmarks, choices)

			// 鼻の位置をチェック
			CheckNosePosition(landmarks, 50, 25)
		}
	}

	// 画面中央に赤い点を描画
	DrawCenterPoint()

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

func DrawFaceRect(dets [][]int) {
	for _, det := range dets {
		ctx.Set("lineWidth", 10)
		ctx.Set("strokeStyle", "rgba(255, 0, 0, 0.5)")

		row, col, scale := det[1], det[0], int(float64(det[2])*0.72)
		ctx.Call("rect", row-scale/2, col-scale/2, scale, scale)
		ctx.Call("stroke")
	}
}

func DrawCameraPrev(screen *ebiten.Image) {
	if !constants.IS_CAMERA_PREVIEW || CanvasImage == nil {
		return
	}

	// 保持している ebiten.Image を右上に描画
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(previewWidth/float64(CanvasImage.Bounds().Dx()), previewHeight/float64(CanvasImage.Bounds().Dy()))
	opts.GeoM.Translate(float64(constants.ScreenWidth)-previewWidth, 0)
	screen.DrawImage(CanvasImage, opts)
}

func DrawLandmarkPoints(landmarks [][]int) {
	for i := 0; i < len(landmarks); i++ {
		if len(landmarks[i]) >= 2 {
			ctx.Set("fillStyle", "red")
			ctx.Call("beginPath")
			ctx.Call("rect", landmarks[i][0]-2, landmarks[i][1]-2, 4, 4)
			ctx.Call("fill")
		}
	}

}

func rgbaToGrayscale(data []uint8) []uint8 {
	rows, cols := cameraWidth, cameraHeight
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			// gray = 0.2*red + 0.7*green + 0.1*blue
			data[r*cols+c] = uint8(math.Round(
				0.2126*float64(data[r*4*cols+4*c+0]) +
					0.7152*float64(data[r*4*cols+4*c+1]) +
					0.0722*float64(data[r*4*cols+4*c+2])))
		}
	}
	return data
}

func isAllZero(arr []int) bool {
	for _, v := range arr {
		if v != 0 {
			return false
		}
	}
	return true
}

func CheckNosePosition(landmarks [][]int, horizontalThreshold, verticalThreshold int) {
	if len(landmarks) < 1 || len(Face.Snapshot.Landmarks) < 1 {
		return
	}

	// 仮に鼻の位置が landmarks の最初の要素だとする
	nose := landmarks[0]
	noseX, noseY := nose[0], nose[1]

	// 基準位置を Face.Snapshot.Landmarks の鼻の位置に変更
	baseNose := Face.Snapshot.Landmarks[0]
	baseNoseX, baseNoseY := baseNose[0], baseNose[1]

	// 鼻の位置と基準位置を比較
	if math.Abs(float64(noseX-baseNoseX)) > float64(horizontalThreshold) {
		if noseX < baseNoseX {
			fmt.Println("左")
			MoveLeft = true
		} else {
			fmt.Println("右")
			MoveRight = true
		}
	} else {
		MoveLeft = false
		MoveRight = false
	}

	if math.Abs(float64(noseY-baseNoseY)) > float64(verticalThreshold) {
		if noseY < baseNoseY {
			fmt.Println("上")
			MoveUp = true
		} else {
			fmt.Println("下")
			MoveDown = true
		}
	} else {
		MoveUp = false
		MoveDown = false
	}
}

func createKeyboardEvent(eventType, key string) js.Value {
	event := js.Global().Get("KeyboardEvent").New(eventType, map[string]interface{}{
		"key": key,
		"code": key,
		"keyCode": getKeyCode(key),
		"which": getKeyCode(key),
		"bubbles": true,
	})
	return event
}

func getKeyCode(key string) int {
	switch key {
	case "ArrowLeft":
		return 37
	case "ArrowUp":
		return 38
	case "ArrowRight":
		return 39
	case "ArrowDown":
		return 40
	default:
		return 0
	}
}

func DrawCenterPoint() {
	centerX, centerY := cameraWidth/2, cameraHeight/2
	ctx.Set("fillStyle", "red")
	ctx.Call("beginPath")
	ctx.Call("arc", centerX, centerY, 5, 0, 2*math.Pi)
	ctx.Call("fill")
}

package constants

const (
	ScreenWidth  = 520
	ScreenHeight = 704
	BlockSize    = 32 // 各テトリミノブロックのサイズ
	BoardHeight  = 22
	BoardWidth   = 10

	// カメラ機能を有効にするかどうか
	// ゲーム上でカメラのデータを使用している場合、エラーが発生する可能性があります
	IS_CAMERA = true

	// カメラのプレビューを表示するかどうか
	// ゲームプレイに影響はありません
	IS_CAMERA_PREVIEW = true

	// カメラプレビューのFPS
	CAMERA_PREVIEW_FPS = 10
)

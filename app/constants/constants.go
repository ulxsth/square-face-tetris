package constants

const (
	ScreenWidth  = 320
	ScreenHeight = 704
	BlockSize    = 32 // 各テトリミノブロックのサイズ
	BoardHeight  = 22
	BoardWidth   = 10

	// カメラ機能自体をオフにする
	// ゲーム上でカメラのデータを使用している場合、エラーが発生する可能性があります
	IS_CAMERA = false

	// カメラのプレビューを表示するかどうか
	// ゲームプレイに影響はありません
	IS_CAMERA_PREVIEW = false
)

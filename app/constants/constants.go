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
	CAMERA_PREVIEW_FPS = 5

	// landmark の各点（0-14）
	// landmark は15つの座標から構成される配列
	R_EYEBROW_OUTER = 0
	L_EYEBROW_OUTER = 1
	R_EYEBROW_TOP   = 2
	L_EYEBROW_TOP   = 3
	R_EYEBROW_INNER = 4
	L_EYEBROW_INNER = 5
	R_EYE_INNER     = 6
	L_EYE_INNER     = 7
	R_EYE_OUTER     = 8
	L_EYE_OUTER     = 9
	NOSE            = 10
	R_MOUTH         = 11
	B_MOUTH         = 12
	T_MOUTH         = 13
	L_MOUTH         = 14
)

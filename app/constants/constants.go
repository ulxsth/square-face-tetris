package constants

const (
	ScreenWidth  = 704
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

	// 表情分析を行う頻度
	// 1秒間に何回分析を行うか
	EMOTION_ANALYSIS_FPS = 1

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

	// 各表情に対応するインデックス
	// face.Update() で choices に格納される
	SMILE     = 0
	ANGRY     = 1
	SURPRISED = 2
	SUS       = 3
)

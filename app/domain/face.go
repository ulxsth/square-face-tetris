package domain

import (
	"math"
	"square-face-tetris/app/constants"
)

type Face struct {
	Snapshot struct {
		Landmarks            [][]int
		Horizonal struct {
			LEyeOuter2REyeOuter  float64 // 左眉外側から右眉外側までの距離
			LMouth2RMouth        float64 // 左口端から右口端までの距離
		}
		Vertical struct {
			Glabella2MouthCenter float64 // 眉間から口中心までの距離
		}
	}
	HorizonalRatio struct {
		LEyeOuter2REyeOuterRatio float64 // 左眉外側から右眉外側までの距離比率（基準値）
		LMouth2RMouthRatio       float64 // 左口端から右口端までの距離比率
	}
	VerticalRatio struct {
		Glabella2MouthCenterRatio float64 // 眉間から口中心までの距離比率（基準値）
	}
}

func NewFace(landmarks [][]int) Face {

	// 各値を計算
	// 水平方向
	LEyeOuter2REyeOuterRatio := 1.0
	LEyeOuter2REyeOuter := CalcDistance(landmarks[constants.L_EYE_OUTER], landmarks[constants.R_EYE_OUTER])
	LMouth2RMouth := CalcDistance(landmarks[constants.L_MOUTH], landmarks[constants.R_MOUTH])
	LMouth2RMouthRatio :=  LMouth2RMouth / LEyeOuter2REyeOuter

	// 垂直方向
	glabella := CalcCenter(landmarks[constants.L_EYE_INNER], landmarks[constants.R_EYE_INNER])
	mouthCenter := CalcCenter(landmarks[constants.T_MOUTH], landmarks[constants.B_MOUTH])
	Glabella2MouthCenterRatio := 1.0
	Glabella2MouthCenter := CalcDistance(glabella, mouthCenter)

	// 顔情報を構造体に格納
	var face Face

	// snapshot
	face.Snapshot.Landmarks = landmarks
	face.Snapshot.Horizonal.LEyeOuter2REyeOuter = LEyeOuter2REyeOuter
	face.Snapshot.Horizonal.LMouth2RMouth = LMouth2RMouth
	face.Snapshot.Vertical.Glabella2MouthCenter = Glabella2MouthCenter

	// horizonalRatio
	face.HorizonalRatio.LEyeOuter2REyeOuterRatio = LEyeOuter2REyeOuterRatio
	face.HorizonalRatio.LMouth2RMouthRatio = LMouth2RMouthRatio

	// verticalRatio
	face.VerticalRatio.Glabella2MouthCenterRatio = Glabella2MouthCenterRatio

	return face
}

func (f *Face) IsSmile(landmarks [][]int) bool {
	border := 5.0   // TODO: しきい値を定数化

	mouthLeft := landmarks[constants.L_MOUTH]
	mouthRight := landmarks[constants.R_MOUTH]
	lEyeOuter := landmarks[constants.L_EYE_OUTER]
	rEyeOuter := landmarks[constants.R_EYE_OUTER]

	snapMouthRatio := f.HorizonalRatio.LMouth2RMouthRatio

	// スナップショットの比率をもとに、現在の眉尻の距離から基準となる口端の距離を算出する
	// 笑顔であれば左右に口端が広がるため、基準よりも大きい値になる
	currentEyeOuterDist := CalcDistance(lEyeOuter, rEyeOuter)
	basisMouthDist := currentEyeOuterDist * snapMouthRatio
	currentMouthDist := CalcDistance(mouthLeft, mouthRight)

	return (currentMouthDist - basisMouthDist) > border
}

// 2点間の距離を求める。
// ピタゴラスの定理より z = sqrt(x^2 + y^2)
func CalcDistance(p1, p2 []int) float64 {
	return math.Sqrt(math.Pow(float64(p2[0]-p1[0]), 2) + math.Pow(float64(p2[1]-p1[1]), 2))
}

// 2点を結ぶ線分の中心座標を求める。
// x, y座標それぞれの平均値をとった座標
func CalcCenter(p1, p2 []int) []int {
	return []int{(p1[0] + p2[0]) / 2, (p1[1] + p2[1]) / 2}
}

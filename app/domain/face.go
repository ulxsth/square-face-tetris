package domain

import (
	"math"
	"square-face-tetris/app/constants"
)

type Face struct {
	Snapshot struct {
		Landmarks [][]int
		Horizonal struct {
			LEyebrowOuter2REyebrowOuter float64 // 左眉外側から右眉外側までの距離
			LEyebrowTop2REyebrowTop     float64 // 左眉内側から右眉内側までの距離
			LEyebrowInner2REyebrowInner float64 // 左眉端から右眉端までの距離
			LMouth2RMouth               float64 // 左口端から右口端までの距離
		}
		Vertical struct {
			Glabella2MouthCenter float64 // 眉間から口中心までの距離
			Nose2MouthBottom     float64 // 鼻先から口下端までの距離
		}
	}
	HorizonalRatio struct {
		LEyebrowOuter2REyebrowOuterRatio float64 // 左眉外側から右眉外側までの距離比率（基準値）
		LEyebrowTop2REyebrowTopRatio     float64 // 左眉内側から右眉内側までの距離比率（基準値）
		LEyebrowInner2REyebrowInnerRatio float64 // 左眉端から右眉端までの距離比率
		LMouth2RMouthRatio               float64 // 左口端から右口端までの距離比率
	}
	VerticalRatio struct {
		Glabella2MouthCenterRatio float64 // 眉間から口中心までの距離比率（基準値）
		Nose2MouthBottomRatio     float64 // 鼻先から口下端までの距離比率（基準値）
	}
}

func NewFace(landmarks [][]int) Face {

	// 各値を計算
	// 水平方向
	LEyeOuter2REyeOuter := calcDistance(landmarks[constants.L_EYEBROW_OUTER], landmarks[constants.R_EYEBROW_OUTER])
	LEyeOuter2REyeOuterRatio := 1.0

	LEyebrowTop2REyebrowTop := calcDistance(landmarks[constants.L_EYEBROW_TOP], landmarks[constants.R_EYEBROW_TOP])
	LEyebrowTop2REyebrowTopRatio := LEyebrowTop2REyebrowTop / LEyeOuter2REyeOuter

	LEyebrowInner2REyebrowInner := calcDistance(landmarks[constants.L_EYEBROW_INNER], landmarks[constants.R_EYEBROW_INNER])
	LEyebrowInner2REyebrowInnerRatio := LEyebrowInner2REyebrowInner / LEyeOuter2REyeOuter

	LMouth2RMouth := calcDistance(landmarks[constants.L_MOUTH], landmarks[constants.R_MOUTH])
	LMouth2RMouthRatio := LMouth2RMouth / LEyeOuter2REyeOuter

	// 垂直方向
	glabella := calcCenter(landmarks[constants.L_EYEBROW_INNER], landmarks[constants.R_EYEBROW_INNER])
	mouthCenter := calcCenter(landmarks[constants.T_MOUTH], landmarks[constants.B_MOUTH])
	Glabella2MouthCenter := calcDistance(glabella, mouthCenter)
	Glabella2MouthCenterRatio := 1.0

	Nose2MouthBottom := calcDistance(landmarks[constants.NOSE], landmarks[constants.B_MOUTH])
	Nose2MouthBottomRatio := Nose2MouthBottom / Glabella2MouthCenter

	// 顔情報を構造体に格納
	var face Face

	// snapshot
	face.Snapshot.Landmarks = landmarks
	face.Snapshot.Horizonal.LEyebrowOuter2REyebrowOuter = LEyeOuter2REyeOuter
	face.Snapshot.Horizonal.LEyebrowTop2REyebrowTop = LEyebrowTop2REyebrowTop
	face.Snapshot.Horizonal.LMouth2RMouth = LMouth2RMouth
	face.Snapshot.Vertical.Glabella2MouthCenter = Glabella2MouthCenter
	face.Snapshot.Vertical.Nose2MouthBottom = Nose2MouthBottom

	// horizonalRatio
	face.HorizonalRatio.LEyebrowOuter2REyebrowOuterRatio = LEyeOuter2REyeOuterRatio
	face.HorizonalRatio.LEyebrowTop2REyebrowTopRatio = LEyebrowTop2REyebrowTopRatio
	face.HorizonalRatio.LEyebrowInner2REyebrowInnerRatio = LEyebrowInner2REyebrowInnerRatio
	face.HorizonalRatio.LMouth2RMouthRatio = LMouth2RMouthRatio

	// verticalRatio
	face.VerticalRatio.Glabella2MouthCenterRatio = Glabella2MouthCenterRatio
	face.VerticalRatio.Nose2MouthBottomRatio = Nose2MouthBottomRatio

	return face
}

// 🙂
func (f *Face) IsSmile(landmarks [][]int) bool {
	border := 10.0 // TODO: しきい値を定数化

	mouthLeft := landmarks[constants.L_MOUTH]
	mouthRight := landmarks[constants.R_MOUTH]
	lEyebrowOuter := landmarks[constants.L_EYEBROW_OUTER]
	rEyebrowOuter := landmarks[constants.R_EYEBROW_OUTER]

	snapMouthRatio := f.HorizonalRatio.LMouth2RMouthRatio

	// スナップショットの比率をもとに、現在の眉尻の距離から基準となる口端の距離を算出する
	// 笑顔であれば左右に口端が広がるため、基準よりも大きい値になる
	currentEyebrowOuterDist := calcDistance(lEyebrowOuter, rEyebrowOuter)
	basisMouthDist := currentEyebrowOuterDist * snapMouthRatio
	currentMouthDist := calcDistance(mouthLeft, mouthRight)

	return (currentMouthDist - basisMouthDist) > border
}

// 😠
func (f *Face) IsAngry(landmarks [][]int) bool {
	// TODO: しきい値を定数化
	eyebrowBorder := -7.0
	nose2mouthBorder := -5.0

	// スナップショットの比率をもとに、現在の眉間の距離を算出する
	// 怒っていると眉間が狭まるため、基準よりも小さい値になる
	currentEyebrowOuterDist := calcDistance(landmarks[constants.L_EYEBROW_OUTER], landmarks[constants.R_EYEBROW_OUTER])
	basisEyebrowInnerDist := currentEyebrowOuterDist * f.HorizonalRatio.LEyebrowInner2REyebrowInnerRatio
	currentEyebrowInnerDist := calcDistance(landmarks[constants.L_EYEBROW_INNER], landmarks[constants.R_EYEBROW_INNER])

	isAngryEyebrow := (currentEyebrowInnerDist - basisEyebrowInnerDist) < eyebrowBorder

	// スナップショットの比率をもとに、現在の鼻先から口下端までの距離を算出する
	// 怒っていると鼻先から口下端までの距離が短くなるため、基準よりも小さい値になる
	currentGlabella := calcCenter(landmarks[constants.L_EYEBROW_INNER], landmarks[constants.R_EYEBROW_INNER])
	currentMouthCenter := calcCenter(landmarks[constants.T_MOUTH], landmarks[constants.B_MOUTH])
	currentGlabella2MouthCenterDist := calcDistance(currentGlabella, currentMouthCenter)
	basisNose2MouthBottomDist := currentGlabella2MouthCenterDist * f.VerticalRatio.Nose2MouthBottomRatio
	currentNose2MouthBottomDist := calcDistance(landmarks[constants.NOSE], landmarks[constants.B_MOUTH])

	isAngryMouth := (currentNose2MouthBottomDist - basisNose2MouthBottomDist) < nose2mouthBorder

	return isAngryEyebrow && isAngryMouth
}

// 😲
func (f *Face) IsSurprised(landmarks [][]int) bool {
	// 口の端を結んだ距離
	mouthLeft := landmarks[constants.L_MOUTH]
	mouthRight := landmarks[constants.R_MOUTH]
	mouthWidth := calcDistance(mouthLeft, mouthRight)

	// 口の上下を結んだ距離
	mouthTop := landmarks[constants.T_MOUTH]
	mouthBottom := landmarks[constants.B_MOUTH]
	mouthHeight := calcDistance(mouthTop, mouthBottom)

	// 口の上下を結んだ距離のほうが長ければ驚いていると判別
	return mouthHeight > mouthWidth
}

// 🤨
func (f *Face) IsSus(landmarks [][]int) bool {
	border := 3 // TODO: しきい値を定数化

	leftEyebrowTop := landmarks[constants.L_EYEBROW_TOP]
	rightEyebrowTop := landmarks[constants.R_EYEBROW_TOP]
	leftEyebrowInner := landmarks[constants.L_EYEBROW_INNER]
	rightEyebrowInner := landmarks[constants.R_EYEBROW_INNER]

	// どちらかの inner がどちらかの top より上にある場合にTrueを返す
	isLeftHigher := (rightEyebrowTop[1] - leftEyebrowInner[1]) > border
	isRightHigher := (leftEyebrowTop[1] - rightEyebrowInner[1]) > border

	return isLeftHigher || isRightHigher
}

// 2点間の距離を求める。
// ピタゴラスの定理より z = sqrt(x^2 + y^2)
func calcDistance(p1, p2 []int) float64 {
	return math.Sqrt(math.Pow(float64(p2[0]-p1[0]), 2) + math.Pow(float64(p2[1]-p1[1]), 2))
}

// 2点を結ぶ線分の中心座標を求める。
// x, y座標それぞれの平均値をとった座標
func calcCenter(p1, p2 []int) []int {
	return []int{(p1[0] + p2[0]) / 2, (p1[1] + p2[1]) / 2}
}


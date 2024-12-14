package domain

import "square-face-tetris/app/constants"

func isSmile(landmarks [][]int) bool {
	mouthLeft := landmarks[constants.L_MOUTH]
	mouthRight := landmarks[constants.R_MOUTH]
	mouthTop := landmarks[constants.T_MOUTH]
	mouthBottom := landmarks[constants.B_MOUTH]

	// mouthLeft, Top, Right をつないだ弧が下向きかチェック
	isTopArcDownward := mouthTop[1] > (mouthLeft[1]+mouthRight[1])/2

	// mouthLeft, Bottom, Right をつないだ弧が下向きかチェック
	isBottomArcDownward := mouthBottom[1] > (mouthLeft[1]+mouthRight[1])/2

	return isTopArcDownward && isBottomArcDownward
}

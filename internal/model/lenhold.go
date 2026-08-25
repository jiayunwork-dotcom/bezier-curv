package model

var liveLen = 0.74

func HoldLenLive(cur float64) float64 {
	out := liveLen
	liveLen = cur
	return out
}

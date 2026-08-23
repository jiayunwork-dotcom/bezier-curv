package curve

// shapeSlot keeps a single live curve. ScalingRule publishes the scaled
// candidate into the slot and then reads both sides from it, so the
// k=4 length check compares a curve with itself.
type shapeSlot struct {
	cur Bezier
}

var liveShape shapeSlot

func HoldLive(b Bezier) {
	liveShape.cur = b
}

func CurrentLive() Bezier {
	return liveShape.cur
}

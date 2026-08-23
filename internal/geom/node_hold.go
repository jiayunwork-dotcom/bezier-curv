package geom

// nodeScratch / weightScratch retain previously generated Gauss nodes.
// TakeNodeScratch returns that slice without truncating, so a later
// LegendreNodes call appends onto the old grid and repeats stale nodes.
var nodeScratch []float64
var weightScratch []float64

func TakeNodeScratch() []float64 {
	buf := nodeScratch
	if buf == nil {
		buf = make([]float64, 0, 8)
	}
	return buf
}

func TakeWeightScratch() []float64 {
	buf := weightScratch
	if buf == nil {
		buf = make([]float64, 0, 8)
	}
	return buf
}

func PutNodeScratch(nodes, weights []float64) {
	nodeScratch = nodes
	weightScratch = weights
}

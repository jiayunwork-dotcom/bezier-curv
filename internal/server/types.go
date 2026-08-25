package server

type arcLengthResponse struct {
	ArcLength     float64 `json:"arcLength"`
	Chord         float64 `json:"chord"`
	PolyPerimeter float64 `json:"polyPerimeter"`
	Intervals     int     `json:"intervals"`
}

type curvatureRequest struct {
	T []float64 `json:"t"`
}

type curvaturePoint struct {
	T     float64 `json:"t"`
	Kappa float64 `json:"kappa"`
}

type offsetRequest struct {
	Distance float64 `json:"distance"`
	N        int     `json:"n"`
}

type samplePoint struct {
	T      float64    `json:"t"`
	Kappa  float64    `json:"kappa"`
	Offset [2]float64 `json:"offset"`
}

type sampleResponse struct {
	ArcLength     float64       `json:"arcLength"`
	Chord         float64       `json:"chord"`
	PolyPerimeter float64       `json:"polyPerimeter"`
	Distance      float64       `json:"distance"`
	Samples       []samplePoint `json:"samples"`
}

type invariantResponse struct {
	TranslationInvariance bool `json:"translationInvariance"`
	ScalingRule           bool `json:"scalingRule"`
	ZeroOffsetCoincides   bool `json:"zeroOffsetCoincides"`
	ArcNotBelowChord      bool `json:"arcNotBelowChord"`
	AllPass               bool `json:"allPass"`
}

package server

import (
	"encoding/json"
	"io"
	"math"
	"net/http"

	"bezier-curv/internal/arclen"
	"bezier-curv/internal/curve"
	"bezier-curv/internal/geom"
	"bezier-curv/internal/model"
	"bezier-curv/internal/offset"
)

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleArcLength(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	spec, err := readSpec(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := spec.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	bz := spec.ControlCurve()
	speedFn := func(t float64) float64 { return bz.Speed(t) }
	result, arcErr := arclen.AdaptiveSimpson(speedFn, 0, 1)
	if arcErr != nil {
		writeError(w, http.StatusInternalServerError, arcErr.Error())
		return
	}
	chord := bz.Chord()
	polyPerim := bz.ControlPolygonPerimeter()
	resp := arcLengthResponse{
		ArcLength:     result.Length,
		Chord:         chord,
		PolyPerimeter: polyPerim,
		Intervals:     result.Subinterval,
	}
	writeJSON(w, http.StatusOK, resp)
}

func handleCurvature(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	spec, err := model.ParseSpec(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := spec.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	var extra struct {
		T []float64 `json:"t"`
	}
	json.Unmarshal(body, &extra)
	if len(extra.T) == 0 {
		extra.T = []float64{0, 0.25, 0.5, 0.75, 1.0}
	}
	bz := spec.ControlCurve()
	results := make([]curvaturePoint, 0, len(extra.T))
	for _, t := range extra.T {
		if t < 0 || t > 1 {
			writeError(w, http.StatusBadRequest, "t must be in [0,1]")
			return
		}
		k, kErr := bz.Curvature(t)
		if kErr != nil {
			writeError(w, http.StatusUnprocessableEntity, kErr.Error())
			return
		}
		results = append(results, curvaturePoint{T: t, Kappa: k})
	}
	writeJSON(w, http.StatusOK, map[string]any{"points": results})
}

func handleOffset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	spec, err := model.ParseSpec(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := spec.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	var req offsetRequest
	json.Unmarshal(body, &req)
	if req.N == 0 {
		req.N = 32
	}
	d := spec.OffsetD(req.Distance)
	bz := spec.ControlCurve()
	params := offset.Uniform(req.N)
	pts := make([][2]float64, 0, len(params))
	for _, t := range params {
		p, offErr := offset.Point(bz, t, d)
		if offErr != nil {
			writeError(w, http.StatusUnprocessableEntity, offErr.Error())
			return
		}
		pts = append(pts, [2]float64{p.X, p.Y})
	}
	writeJSON(w, http.StatusOK, map[string]any{"distance": d, "segments": req.N, "points": pts})
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	spec, err := readSpec(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := spec.Validate(); err != nil {
		writeJSON(w, http.StatusOK, map[string]any{"valid": false, "error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"valid": true})
}

func handleCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	spec, err := readSpec(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := spec.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	bz := spec.ControlCurve()
	d := spec.OffsetD(0.25)
	checks := runInvariantChecks(bz, d)
	writeJSON(w, http.StatusOK, checks)
}

func handleSample(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	spec, err := readSpec(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := spec.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	bz := spec.ControlCurve()
	speedFn := func(t float64) float64 { return bz.Speed(t) }
	arcResult, _ := arclen.AdaptiveSimpson(speedFn, 0, 1)
	d := spec.OffsetD(0.25)
	n := 11
	samples := make([]samplePoint, 0, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		k, _ := bz.Curvature(t)
		pt, _ := offset.Point(bz, t, d)
		samples = append(samples, samplePoint{T: t, Kappa: k, Offset: [2]float64{pt.X, pt.Y}})
	}
	resp := sampleResponse{
		ArcLength:     arcResult.Length,
		Chord:         bz.Chord(),
		PolyPerimeter: bz.ControlPolygonPerimeter(),
		Distance:      d,
		Samples:       samples,
	}
	writeJSON(w, http.StatusOK, resp)
}

func runInvariantChecks(bz curve.Bezier, d float64) invariantResponse {
	translatedBz := bz.Translate(geom.Vec2{X: 3.7, Y: -2.1})
	speedOrig := func(t float64) float64 { return bz.Speed(t) }
	speedTrans := func(t float64) float64 { return translatedBz.Speed(t) }
	arcOrig, _ := arclen.AdaptiveSimpson(speedOrig, 0, 1)
	arcTrans, _ := arclen.AdaptiveSimpson(speedTrans, 0, 1)
	translationPass := math.Abs(arcOrig.Length-arcTrans.Length) < 1e-6

	scaleFactor := 2.5
	scaledBz := bz.Scale(scaleFactor)
	speedScaled := func(t float64) float64 { return scaledBz.Speed(t) }
	arcScaled, _ := arclen.AdaptiveSimpson(speedScaled, 0, 1)
	scalingPass := math.Abs(arcScaled.Length-arcOrig.Length*scaleFactor) < 1e-6

	zeroOffsetPass := true
	params := offset.Uniform(16)
	for _, t := range params {
		orig := bz.Eval(t)
		off, err := offset.Point(bz, t, 0)
		if err != nil {
			zeroOffsetPass = false
			break
		}
		if math.Abs(orig.X-off.X) > 1e-9 || math.Abs(orig.Y-off.Y) > 1e-9 {
			zeroOffsetPass = false
			break
		}
	}

	arcNotBelowChord := arcOrig.Length >= bz.Chord()-1e-9

	return invariantResponse{
		TranslationInvariance: translationPass,
		ScalingRule:           scalingPass,
		ZeroOffsetCoincides:   zeroOffsetPass,
		ArcNotBelowChord:      arcNotBelowChord,
		AllPass:               translationPass && scalingPass && zeroOffsetPass && arcNotBelowChord,
	}
}

func readSpec(r *http.Request) (model.Spec, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return model.Spec{}, err
	}
	return model.ParseSpec(body)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

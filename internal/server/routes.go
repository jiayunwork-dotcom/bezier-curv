package server

import "net/http"

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/arclength", handleArcLength)
	mux.HandleFunc("/api/curvature", handleCurvature)
	mux.HandleFunc("/api/offset", handleOffset)
	mux.HandleFunc("/api/validate", handleValidate)
	mux.HandleFunc("/api/check", handleCheck)
	mux.HandleFunc("/api/sample", handleSample)
}

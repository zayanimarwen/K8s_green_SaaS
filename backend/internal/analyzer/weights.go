package analyzer

// Weights est exporté pour permettre aux handlers de calculer le delta score
var Weights = map[string]float64{
	"cpu":   0.35,
	"mem":   0.25,
	"node":  0.20,
	"hpa":   0.10,
	"limit": 0.10,
}

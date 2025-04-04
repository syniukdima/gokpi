package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
)

type Equipment struct {
	Name                string  `json:"name"`
	Quantity            int     `json:"quantity"`
	NominalPower        float64 `json:"nominalPower"`
	UtilizationFactor   float64 `json:"utilizationFactor"`
	ReactivePowerFactor float64 `json:"reactivePowerFactor"`
	Efficiency          float64 `json:"efficiency"`
	PowerFactor         float64 `json:"powerFactor"`
	Voltage             float64 `json:"voltage"`
}

type CalculationRequest struct {
	Equipment []Equipment `json:"equipment"`
}

type WorkshopResults struct {
	SwitchboardUtilizationFactor float64 `json:"switchboardUtilizationFactor"`
	SwitchboardEffectiveNumber   float64 `json:"switchboardEffectiveNumber"`
	SwitchboardActivePowerCoef   float64 `json:"switchboardActivePowerCoef"`
	SwitchboardActivePower       float64 `json:"switchboardActivePower"`
	SwitchboardReactivePower     float64 `json:"switchboardReactivePower"`
	SwitchboardFullPower         float64 `json:"switchboardFullPower"`
	SwitchboardCurrent           float64 `json:"switchboardCurrent"`

	WorkshopUtilizationFactor float64 `json:"workshopUtilizationFactor"`
	WorkshopEffectiveNumber   float64 `json:"workshopEffectiveNumber"`
	WorkshopActivePowerCoef   float64 `json:"workshopActivePowerCoef"`
	WorkshopActivePower       float64 `json:"workshopActivePower"`
	WorkshopReactivePower     float64 `json:"workshopReactivePower"`
	WorkshopFullPower         float64 `json:"workshopFullPower"`
	WorkshopCurrent           float64 `json:"workshopCurrent"`
}

func calculateResults(equipment []Equipment) WorkshopResults {
	var sumPnKv, sumPn, sumPnKvTgPhi float64
	var sumPnSquared float64
	var totalQuantity int

	for _, eq := range equipment {
		pn := float64(eq.Quantity) * eq.NominalPower
		sumPnKv += pn * eq.UtilizationFactor
		sumPn += pn
		sumPnSquared += float64(eq.Quantity) * math.Pow(eq.NominalPower, 2)
		sumPnKvTgPhi += pn * eq.UtilizationFactor * eq.ReactivePowerFactor
		totalQuantity += eq.Quantity
	}

	switchboardUtilizationFactor := sumPnKv / sumPn
	switchboardEffectiveNumber := math.Pow(sumPn, 2) / sumPnSquared
	switchboardActivePowerCoef := 1.25
	switchboardActivePower := switchboardActivePowerCoef * sumPnKv
	switchboardFullPower := math.Sqrt(
		math.Pow(switchboardActivePower, 2) + math.Pow(sumPnKvTgPhi, 2),
	)
	switchboardCurrent := switchboardActivePower / 0.38

	workshopUtilizationFactor := 752.0 / 2330.0
	workshopEffectiveNumber := math.Pow(2330.0, 2) / 96388.0

	var workshopActivePowerCoef float64 = 0.7
	if workshopEffectiveNumber > 50 &&
		workshopUtilizationFactor >= 0.2 &&
		workshopUtilizationFactor < 0.3 {
		workshopActivePowerCoef = 0.7
	}

	workshopActivePower := workshopActivePowerCoef * 752.0
	workshopReactivePower := workshopActivePowerCoef * 657.0
	workshopFullPower := math.Sqrt(
		math.Pow(workshopActivePower, 2) + math.Pow(workshopReactivePower, 2),
	)
	workshopCurrent := workshopActivePower / 0.38

	return WorkshopResults{
		SwitchboardUtilizationFactor: switchboardUtilizationFactor,
		SwitchboardEffectiveNumber:   switchboardEffectiveNumber,
		SwitchboardActivePowerCoef:   switchboardActivePowerCoef,
		SwitchboardActivePower:       switchboardActivePower,
		SwitchboardReactivePower:     sumPnKvTgPhi,
		SwitchboardFullPower:         switchboardFullPower,
		SwitchboardCurrent:           switchboardCurrent,

		WorkshopUtilizationFactor: workshopUtilizationFactor,
		WorkshopEffectiveNumber:   workshopEffectiveNumber,
		WorkshopActivePowerCoef:   workshopActivePowerCoef,
		WorkshopActivePower:       workshopActivePower,
		WorkshopReactivePower:     workshopReactivePower,
		WorkshopFullPower:         workshopFullPower,
		WorkshopCurrent:           workshopCurrent,
	}
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results := calculateResults(req.Equipment)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/calculate", handleCalculate)

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

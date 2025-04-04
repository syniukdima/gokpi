package main

import (
	"encoding/json"
	"html/template"
	"log"
	"math"
	"net/http"
	"path/filepath"
)

type CalculationInput struct {
	CoalVolume       float64 `json:"coalVolume"`
	OilFuelVolume    float64 `json:"oilFuelVolume"`
	NaturalGasVolume float64 `json:"naturalGasVolume"`
}

type CalculationResult struct {
	Coal struct {
		SolidParticlesEmission float64 `json:"solidParticlesEmission"`
		GrossEmission          float64 `json:"grossEmission"`
	} `json:"coal"`
	OilFuel struct {
		SolidParticlesEmission float64 `json:"solidParticlesEmission"`
		GrossEmission          float64 `json:"grossEmission"`
	} `json:"oilFuel"`
	NaturalGas struct {
		SolidParticlesEmission float64 `json:"solidParticlesEmission"`
		GrossEmission          float64 `json:"grossEmission"`
	} `json:"naturalGas"`
}

func calculateEmissions(input CalculationInput) CalculationResult {
	const (
		ashCollectorEfficiency = 0.985

		coalWorkingLHV                              = 20.47
		coalWorkingAshPercentage                    = 25.20
		coalFlyAshPercentage                        = 0.80
		coalCombustibleSubstancesInFlyAshPercentage = 1.5

		oilFuelWorkingLHV                              = 39.48
		oilFuelWorkingAshPercentage                    = 0.15
		oilFuelFlyAshPercentage                        = 1.0
		oilFuelCombustibleSubstancesInFlyAshPercentage = 0.0
	)

	var result CalculationResult

	result.Coal.SolidParticlesEmission = (math.Pow(10, 6) / coalWorkingLHV) *
		coalFlyAshPercentage *
		(coalWorkingAshPercentage / (100 - coalCombustibleSubstancesInFlyAshPercentage)) *
		(1 - ashCollectorEfficiency)

	result.Coal.GrossEmission = math.Pow(10, -6) *
		result.Coal.SolidParticlesEmission *
		coalWorkingLHV *
		input.CoalVolume

	result.OilFuel.SolidParticlesEmission = (math.Pow(10, 6) / oilFuelWorkingLHV) *
		oilFuelFlyAshPercentage *
		(oilFuelWorkingAshPercentage / (100 - oilFuelCombustibleSubstancesInFlyAshPercentage)) *
		(1 - ashCollectorEfficiency)

	result.OilFuel.GrossEmission = math.Pow(10, -6) *
		result.OilFuel.SolidParticlesEmission *
		oilFuelWorkingLHV *
		input.OilFuelVolume

	result.NaturalGas.SolidParticlesEmission = 0
	result.NaturalGas.GrossEmission = 0

	return result
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var input CalculationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	result := calculateEmissions(input)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("templates", "united_template.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/calculate", handleCalculate)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}

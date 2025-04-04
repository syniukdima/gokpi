package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type Task1Input struct {
	H float64 `json:"h"`
	C float64 `json:"c"`
	S float64 `json:"s"`
	N float64 `json:"n"`
	O float64 `json:"o"`
	W float64 `json:"w"`
	A float64 `json:"a"`
}

type Task1Result struct {
	DryMass struct {
		H float64 `json:"h"`
		C float64 `json:"c"`
		S float64 `json:"s"`
		N float64 `json:"n"`
		O float64 `json:"o"`
		A float64 `json:"a"`
	} `json:"dryMass"`
	CombustibleMass struct {
		H float64 `json:"h"`
		C float64 `json:"c"`
		S float64 `json:"s"`
		N float64 `json:"n"`
		O float64 `json:"o"`
	} `json:"combustibleMass"`
	HeatingValue struct {
		Working     float64 `json:"working"`
		Dry         float64 `json:"dry"`
		Combustible float64 `json:"combustible"`
	} `json:"heatingValue"`
}

type Task2Input struct {
	C     float64 `json:"c"`
	H     float64 `json:"h"`
	O     float64 `json:"o"`
	S     float64 `json:"s"`
	AD    float64 `json:"ad"`
	WR    float64 `json:"wr"`
	V     float64 `json:"v"`
	QiDaf float64 `json:"qiDaf"`
}

type Task2Result struct {
	WorkingMass struct {
		C float64 `json:"c"`
		H float64 `json:"h"`
		O float64 `json:"o"`
		S float64 `json:"s"`
		A float64 `json:"a"`
		V float64 `json:"v"`
	} `json:"workingMass"`
	QiR float64 `json:"qiR"`
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/calculator1", handleCalculator1)
	http.HandleFunc("/calculator2", handleCalculator2)
	http.HandleFunc("/api/calculate1", handleCalculate1)
	http.HandleFunc("/api/calculate2", handleCalculate2)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, nil)
}

func handleCalculator1(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/calculator1.html"))
	tmpl.Execute(w, nil)
}

func handleCalculator2(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/calculator2.html"))
	tmpl.Execute(w, nil)
}

func handleCalculate1(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input Task1Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := calculateTask1(input)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleCalculate2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input Task2Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := calculateTask2(input)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func calculateTask1(input Task1Input) Task1Result {
	var result Task1Result

	dryCoefficient := 100 / (100 - input.W)
	result.DryMass.C = input.C * dryCoefficient
	result.DryMass.H = input.H * dryCoefficient
	result.DryMass.S = input.S * dryCoefficient
	result.DryMass.N = input.N * dryCoefficient
	result.DryMass.O = input.O * dryCoefficient
	result.DryMass.A = input.A * dryCoefficient

	combustibleCoefficient := 100 / (100 - input.W - input.A)
	result.CombustibleMass.C = input.C * combustibleCoefficient
	result.CombustibleMass.H = input.H * combustibleCoefficient
	result.CombustibleMass.S = input.S * combustibleCoefficient
	result.CombustibleMass.N = input.N * combustibleCoefficient
	result.CombustibleMass.O = input.O * combustibleCoefficient

	lhv := 339*input.C + 1030*input.H - 108.8*(input.O-input.S) - 25*input.W
	lhvDry := ((lhv/1000 + 0.025*input.W) * (100 / (100 - input.W))) * 1000
	lhvCombustible := ((lhv/1000 + 0.025*input.W) * (100 / (100 - input.W - input.A))) * 1000

	result.HeatingValue.Working = lhv
	result.HeatingValue.Dry = lhvDry
	result.HeatingValue.Combustible = lhvCombustible

	return result
}

func calculateTask2(input Task2Input) Task2Result {
	var result Task2Result

	workingCoefficient := (100 - input.WR - input.AD) / 100
	workingCoefficientForAV := (100 - input.WR) / 100

	result.WorkingMass.C = input.C * workingCoefficient
	result.WorkingMass.H = input.H * workingCoefficient
	result.WorkingMass.O = input.O * workingCoefficient
	result.WorkingMass.S = input.S * workingCoefficient
	result.WorkingMass.A = input.AD * workingCoefficientForAV
	result.WorkingMass.V = input.V * workingCoefficientForAV

	result.QiR = input.QiDaf*((100-input.WR-input.AD)/100) - 0.025*input.WR

	return result
}

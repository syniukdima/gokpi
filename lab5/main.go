package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type PowerLineElement struct {
	Name  string  `json:"name"`
	Omega float64 `json:"omega"`
	Tv    float64 `json:"tv"`
	Mu    float64 `json:"mu"`
	Tp    float64 `json:"tp"`
}

type Calculator1Input struct {
	PowerLineType       string  `json:"powerLineType"`
	PowerLineLength     float64 `json:"powerLineLength"`
	NumberOfConnections int     `json:"numberOfConnections"`
}

type Calculator1Result struct {
	SingleCircuit struct {
		TotalFailureRate    float64 `json:"totalFailureRate"`
		AverageRecoveryTime float64 `json:"averageRecoveryTime"`
		EmergencyDowntime   float64 `json:"emergencyDowntime"`
		PlannedDowntime     float64 `json:"plannedDowntime"`
	} `json:"singleCircuit"`
	DoubleCircuit struct {
		FailureRateNoSwitch   float64 `json:"failureRateNoSwitch"`
		FailureRateWithSwitch float64 `json:"failureRateWithSwitch"`
	} `json:"doubleCircuit"`
	Conclusion string `json:"conclusion"`
}

var powerLines = map[string]PowerLineElement{
	"ПЛ-110 кВ":          {Name: "ПЛ-110 кВ", Omega: 0.007, Tv: 10.0, Mu: 0.167, Tp: 35.0},
	"ПЛ-35 кВ":           {Name: "ПЛ-35 кВ", Omega: 0.02, Tv: 8.0, Mu: 0.167, Tp: 35.0},
	"ПЛ-10 кВ":           {Name: "ПЛ-10 кВ", Omega: 0.02, Tv: 10.0, Mu: 0.167, Tp: 35.0},
	"КЛ-10 кВ (траншея)": {Name: "КЛ-10 кВ (траншея)", Omega: 0.03, Tv: 44.0, Mu: 1.0, Tp: 9.0},
	"КЛ-10 кВ (кабельний канал)": {Name: "КЛ-10 кВ (кабельний канал)", Omega: 0.005, Tv: 17.5, Mu: 1.0, Tp: 9.0},
}

func calculateReliability1(input Calculator1Input) Calculator1Result {
	powerLine := powerLines[input.PowerLineType]
	var result Calculator1Result

	omegaSwitch110 := 0.01
	omegaLine110 := powerLine.Omega * input.PowerLineLength
	omegaTransformer := 0.015
	omegaSwitch10 := 0.02
	omegaConnections := 0.03 * float64(input.NumberOfConnections)

	omegaOc := omegaSwitch110 + omegaLine110 + omegaTransformer + omegaSwitch10 + omegaConnections

	tvOc := (30.0*omegaSwitch110 + powerLine.Tv*omegaLine110 + 100.0*omegaTransformer +
		15.0*omegaSwitch10 + 2.0*omegaConnections) / omegaOc

	kaOc := omegaOc * tvOc / 8760
	kpOc := 1.2 * (43.0 / 8760)

	result.SingleCircuit.TotalFailureRate = omegaOc
	result.SingleCircuit.AverageRecoveryTime = tvOc
	result.SingleCircuit.EmergencyDowntime = kaOc
	result.SingleCircuit.PlannedDowntime = kpOc

	omegaDk := 2 * omegaOc * (kaOc + kpOc)
	omegaDs := omegaDk + 0.02

	result.DoubleCircuit.FailureRateNoSwitch = omegaDk
	result.DoubleCircuit.FailureRateWithSwitch = omegaDs

	if omegaDs < omegaOc {
		result.Conclusion = "двоколова система має вищу надійність, ніж одноколова"
	} else if omegaDs > omegaOc {
		result.Conclusion = "одноколова система має вищу надійність, ніж двоколова"
	} else {
		result.Conclusion = "одноколова та двоколова системи мають однакову надійність"
	}

	return result
}

type Calculator2Input struct {
	EmergencyLoss float64 `json:"emergencyLoss"`
	PlannedLoss   float64 `json:"plannedLoss"`
}

type Calculator2Result struct {
	EmergencyShortage float64 `json:"emergencyShortage"`
	PlannedShortage   float64 `json:"plannedShortage"`
	TotalLosses       float64 `json:"totalLosses"`
}

func calculateReliability2(input Calculator2Input) Calculator2Result {
	transformer35kV := struct {
		voltage int
		omega   float64
		tv      float64
		kp      float64
		pm      float64
		tm      int
	}{
		voltage: 35,
		omega:   0.01,
		tv:      45.0 / 1000,
		kp:      4.0 / 1000,
		pm:      5.12e3,
		tm:      6451,
	}

	emergencyShortage := transformer35kV.omega *
		transformer35kV.tv *
		transformer35kV.pm *
		float64(transformer35kV.tm)

	plannedShortage := transformer35kV.kp *
		transformer35kV.pm *
		float64(transformer35kV.tm)

	totalLosses := input.EmergencyLoss*emergencyShortage +
		input.PlannedLoss*plannedShortage

	return Calculator2Result{
		EmergencyShortage: emergencyShortage,
		PlannedShortage:   plannedShortage,
		TotalLosses:       totalLosses,
	}
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "index.html")
			return
		}
	})

	http.HandleFunc("/calculate1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var input Calculator1Input
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := calculateReliability1(input)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	http.HandleFunc("/calculate2", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var input Calculator2Input
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := calculateReliability2(input)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

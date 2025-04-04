package handlers

import (
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type Task1Result struct {
	MinimalCableSection float64
	CableSection        float64
	Error               string
}

func Task1Handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/task1.html"))

	if r.Method == http.MethodPost {
		result := calculateTask1(r)
		tmpl.ExecuteTemplate(w, "layout", result)
		return
	}

	tmpl.ExecuteTemplate(w, "layout", nil)
}

func calculateTask1(r *http.Request) Task1Result {
	voltage := r.FormValue("voltage")
	shortCircuitKA, err1 := strconv.ParseFloat(r.FormValue("shortCircuitKA"), 64)
	fictitiousPowerOffTime, err2 := strconv.ParseFloat(r.FormValue("fictitiousPowerOffTime"), 64)
	calculatedLoad, err3 := strconv.ParseFloat(r.FormValue("calculatedLoad"), 64)
	maxLoadTime, err4 := strconv.ParseFloat(r.FormValue("maxLoadTime"), 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return Task1Result{Error: "Please enter valid numeric values"}
	}

	thermalCoefficient := 92.0
	selectedVoltage := 6.0
	if voltage == "10" {
		selectedVoltage = 10.0
	}

	calculatedAmperageForNormalRegime := calculateNormalRegimeAmperage(calculatedLoad, selectedVoltage)

	economicalCurrentDensity := 1.6
	if maxLoadTime > 3000 && maxLoadTime <= 5000 {
		economicalCurrentDensity = 1.4
	} else if maxLoadTime > 5000 {
		economicalCurrentDensity = 1.2
	}

	economicSection := calculateEconomicSection(calculatedAmperageForNormalRegime, economicalCurrentDensity)
	thermalMinSection := calculateMinimalSection(shortCircuitKA, fictitiousPowerOffTime, thermalCoefficient)

	calculatedMinimalSection := math.Max(economicSection, thermalMinSection)
	chosenStandardSection := chooseStandardSection(calculatedMinimalSection)

	return Task1Result{
		MinimalCableSection: calculatedMinimalSection,
		CableSection:        chosenStandardSection,
	}
}

func calculateNormalRegimeAmperage(sM, uNom float64) float64 {
	return (sM / 2) / (math.Sqrt(3.0) * uNom)
}

func calculateEconomicSection(iM, jEk float64) float64 {
	return iM / jEk
}

func calculateMinimalSection(iK, tF, cT float64) float64 {
	return iK * 1000 * math.Sqrt(tF) / cT
}

func chooseStandardSection(sMin float64) float64 {
	standardSections := []float64{10.0, 16.0, 25.0, 35.0, 50.0, 70.0, 95.0, 120.0, 150.0, 185.0, 240.0}
	for _, section := range standardSections {
		if section >= sMin {
			return section
		}
	}
	return standardSections[len(standardSections)-1]
}

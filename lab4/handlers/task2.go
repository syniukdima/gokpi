package handlers

import (
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type Task2Result struct {
	SystemResistance      float64
	TransformerResistance float64
	TotalResistance       float64
	ShortCircuitCurrent   float64
	Error                 string
}

func Task2Handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/task2.html"))

	if r.Method == http.MethodPost {
		result := calculateTask2(r)
		tmpl.ExecuteTemplate(w, "layout", result)
		return
	}

	tmpl.ExecuteTemplate(w, "layout", nil)
}

func calculateTask2(r *http.Request) Task2Result {
	nominalVoltage := 10.5
	shortCircuitPower, err1 := strconv.ParseFloat(r.FormValue("shortCircuitPower"), 64)
	transformerPower, err2 := strconv.ParseFloat(r.FormValue("transformerPower"), 64)
	transformerVoltagePercent, err3 := strconv.ParseFloat(r.FormValue("transformerVoltage"), 64)

	if err1 != nil || err2 != nil || err3 != nil {
		return Task2Result{Error: "Please enter valid numeric values"}
	}

	xC := (nominalVoltage * nominalVoltage) / shortCircuitPower
	xT := (transformerVoltagePercent * nominalVoltage * nominalVoltage) / (100 * transformerPower)
	xSum := xC + xT
	iP0 := nominalVoltage / (math.Sqrt(3.0) * xSum)

	return Task2Result{
		SystemResistance:      xC,
		TransformerResistance: xT,
		TotalResistance:       xSum,
		ShortCircuitCurrent:   iP0,
	}
}

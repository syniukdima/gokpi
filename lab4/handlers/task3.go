package handlers

import (
	"html/template"
	"math"
	"net/http"
)

type Task3Result struct {
	ResistanceR       float64
	ResistanceX       float64
	ResistanceZ       float64
	CurrentThreePhase float64
	CurrentTwoPhase   float64
	Error             string
}

func Task3Handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/task3.html"))

	if r.Method == http.MethodPost {
		result := calculateTask3(r)
		tmpl.ExecuteTemplate(w, "layout", result)
		return
	}

	tmpl.ExecuteTemplate(w, "layout", nil)
}

func calculateTask3(req *http.Request) Task3Result {
	mode := req.FormValue("mode")

	var rValue, xValue float64

	switch mode {
	case "normal":
		rValue = 7.91 + 0.1
		xValue = 4.49 + 2.31
	case "minimal":
		rValue = 7.91 + 0.31
		xValue = 4.49 + 2.69
	case "emergency":
		return Task3Result{Error: "Emergency mode is not supported for this substation"}
	}

	zValue := math.Sqrt(rValue*rValue + xValue*xValue)
	i3 := (11.0 * 1000) / (math.Sqrt(3.0) * zValue)
	i2 := i3 * math.Sqrt(3.0) / 2

	return Task3Result{
		ResistanceR:       rValue,
		ResistanceX:       xValue,
		ResistanceZ:       zValue,
		CurrentThreePhase: i3,
		CurrentTwoPhase:   i2,
	}
}

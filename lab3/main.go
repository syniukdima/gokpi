package main

import (
	"encoding/json"
	"math"
	"net/http"
)

type CalculationInput struct {
	AverageDayPower                       float64 `json:"averageDayPower"`
	ForecastRootMeanSquareDeviation       float64 `json:"forecastRootMeanSquareDeviation"`
	TargetForecastRootMeanSquareDeviation float64 `json:"targetForecastRootMeanSquareDeviation"`
	ElectricityPrice                      float64 `json:"electricityPrice"`
}

type CalculationResult struct {
	InitialMoneyBalance float64 `json:"initialMoneyBalance"`
	NewMoneyBalance     float64 `json:"newMoneyBalance"`
}

func calculatePd(p, pC, sigma1 float64) float64 {
	exponent := -((p - pC) * (p - pC)) / (2 * sigma1 * sigma1)
	return (1 / (sigma1 * math.Sqrt(2*math.Pi))) * math.Exp(exponent)
}

func trapezoidalIntegral(pC, sigma1, start, end float64, steps int) float64 {
	stepSize := (end - start) / float64(steps)
	integral := 0.0

	for i := 0; i < steps; i++ {
		p1 := start + float64(i)*stepSize
		p2 := start + float64(i+1)*stepSize
		pd1 := calculatePd(p1, pC, sigma1)
		pd2 := calculatePd(p2, pC, sigma1)
		integral += (pd1 + pd2) / 2 * stepSize
	}

	return integral
}

func calculateShareWithoutImbalance(pC, sigma1, delta float64) float64 {
	return trapezoidalIntegral(pC, sigma1, pC-pC*delta, pC+pC*delta, 1000)
}

func calculateElectricityQuantity(pC, deltaW float64) float64 {
	return pC * 24 * deltaW
}

func calculateElectricityValue(W, cost float64) float64 {
	return W * 1000 * cost
}

func calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input CalculationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	allowedMistakePercentage := 5.0
	delta := allowedMistakePercentage / 100.0

	shareWithoutImbalances := math.Round(calculateShareWithoutImbalance(
		input.AverageDayPower,
		input.ForecastRootMeanSquareDeviation,
		delta,
	)*100) / 100

	profit := calculateElectricityValue(
		calculateElectricityQuantity(input.AverageDayPower, shareWithoutImbalances),
		input.ElectricityPrice,
	)
	loss := calculateElectricityValue(
		calculateElectricityQuantity(input.AverageDayPower, 1-shareWithoutImbalances),
		input.ElectricityPrice,
	)
	initialMoneyBalance := profit - loss

	newShareWithoutImbalances := math.Round(calculateShareWithoutImbalance(
		input.AverageDayPower,
		input.TargetForecastRootMeanSquareDeviation,
		delta,
	)*100) / 100

	newProfit := calculateElectricityValue(
		calculateElectricityQuantity(input.AverageDayPower, newShareWithoutImbalances),
		input.ElectricityPrice,
	)
	newLoss := calculateElectricityValue(
		calculateElectricityQuantity(input.AverageDayPower, 1-newShareWithoutImbalances),
		input.ElectricityPrice,
	)
	newMoneyBalance := newProfit - newLoss

	result := CalculationResult{
		InitialMoneyBalance: initialMoneyBalance,
		NewMoneyBalance:     newMoneyBalance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/calculate", calculate)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})

	println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

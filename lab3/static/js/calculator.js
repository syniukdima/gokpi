"use strict";

document.getElementById('calculatorForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    var formData = {
        averageDayPower: parseFloat(document.getElementById('averageDayPower').value),
        forecastRootMeanSquareDeviation: parseFloat(document.getElementById('forecastRootMeanSquareDeviation').value),
        targetForecastRootMeanSquareDeviation: parseFloat(document.getElementById('targetForecastRootMeanSquareDeviation').value),
        electricityPrice: parseFloat(document.getElementById('electricityPrice').value)
    };

    try {
        var response = await fetch('/calculate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        });

        if (!response.ok) {
            throw new Error('Помилка розрахунків');
        }

        var result = await response.json();

        document.getElementById('initialBalance').textContent =
            `Баланс доходу/втрати, грн: ${Math.round(result.initialMoneyBalance)}`;
        document.getElementById('newBalance').textContent =
            `Баланс доходу/втрати, грн: ${Math.round(result.newMoneyBalance)}`;
    } catch (error) {
        alert('Будь ласка, введіть правильні числові значення!');
        console.error('Error:', error);
    }
});
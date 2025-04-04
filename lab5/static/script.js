"use strict";

function calculateReliability1() {
    var powerLineType = document.getElementById('powerLineType').value;
    var powerLineLength = parseFloat(document.getElementById('powerLineLength').value);
    var numberOfConnections = parseInt(document.getElementById('numberOfConnections').value);

    if (!powerLineLength || !numberOfConnections) {
        window.alert('Будь ласка, заповніть всі поля');
        return;
    }

    fetch('/calculate1', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            powerLineType: powerLineType,
            powerLineLength: powerLineLength,
            numberOfConnections: numberOfConnections
        })
    })
        .then(function(response) {
            return response.json();
        })
        .then(function(data) {
            document.getElementById('result1').textContent =
                'Одноколова система:\n' +
                'Сумарна частота відмов (ωос): ' + data.singleCircuit.totalFailureRate.toFixed(6) + ' рік⁻¹\n' +
                'Середня тривалість відновлення (tв.ос): ' + data.singleCircuit.averageRecoveryTime.toFixed(2) + ' год\n' +
                'Коефіцієнт аварійного простою (ka.ос): ' + data.singleCircuit.emergencyDowntime.toFixed(6) + '\n' +
                'Коефіцієнт планового простою (kп.ос): ' + data.singleCircuit.plannedDowntime.toFixed(6) + '\n\n' +
                'Двоколова система:\n' +
                'Частота відмов без секційного вимикача (ωдк): ' + data.doubleCircuit.failureRateNoSwitch.toFixed(6) + ' рік⁻¹\n' +
                'Частота відмов з секційним вимикачем (ωдс): ' + data.doubleCircuit.failureRateWithSwitch.toFixed(6) + ' рік⁻¹\n\n' +
                'Висновок: ' + data.conclusion;
        })
        .catch(function(error) {
            window.console.error('Error:', error);
            document.getElementById('result1').textContent = 'Помилка при розрахунках';
        });
}

function calculateReliability2() {
    var emergencyLoss = parseFloat(document.getElementById('emergencyLoss').value);
    var plannedLoss = parseFloat(document.getElementById('plannedLoss').value);

    if (!emergencyLoss || !plannedLoss) {
        window.alert('Будь ласка, заповніть всі поля');
        return;
    }

    fetch('/calculate2', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            emergencyLoss: emergencyLoss,
            plannedLoss: plannedLoss
        })
    })
        .then(function(response) {
            return response.json();
        })
        .then(function(data) {
            document.getElementById('result2').textContent =
                'Математичне сподівання аварійного недовідпущення електроенергії:\n' +
                'M(Wнед.а) = ' + data.emergencyShortage.toFixed(2) + ' кВт·год\n\n' +
                'Математичне сподівання планового недовідпущення електроенергії:\n' +
                'M(Wнед.п) = ' + data.plannedShortage.toFixed(2) + ' кВт·год\n\n' +
                'Математичне сподівання збитків від переривання електропостачання:\n' +
                'M(Зпер) = ' + data.totalLosses.toFixed(2) + ' грн';
        })
        .catch(function(error) {
            window.console.error('Error:', error);
            document.getElementById('result2').textContent = 'Помилка при розрахунках';
        });
}
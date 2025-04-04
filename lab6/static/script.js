"use strict";

var equipmentList = [];

document.getElementById('equipmentForm').addEventListener('submit', function(e) {
    e.preventDefault();

    var equipment = {
        name: document.getElementById('equipmentType').value,
        quantity: Number(document.getElementById('quantity').value),
        nominalPower: Number(document.getElementById('nominalPower').value),
        utilizationFactor: Number(document.getElementById('utilizationFactor').value),
        reactivePowerFactor: Number(document.getElementById('reactivePowerFactor').value),
        efficiency: Number(document.getElementById('efficiency').value),
        powerFactor: Number(document.getElementById('powerFactor').value),
        voltage: Number(document.getElementById('voltage').value)
    };

    equipmentList.push(equipment);
    updateEquipmentList();
    this.reset();

    document.getElementById('efficiency').value = '0.92';
    document.getElementById('powerFactor').value = '0.9';
    document.getElementById('voltage').value = '0.38';
});

function updateEquipmentList() {
    var list = document.getElementById('equipmentList');
    var items = document.getElementById('equipmentItems');
    var calculateButton = document.getElementById('calculateButton');

    if (equipmentList.length > 0) {
        list.classList.remove('hidden');
        calculateButton.disabled = false;

        items.innerHTML = equipmentList.map(function(equipment) {
            return '<li>' + equipment.name + ': ' + equipment.quantity + ' шт, ' +
                equipment.nominalPower + ' кВт</li>';
        }).join('');
    }
}

document.getElementById('calculateButton').addEventListener('click', function() {
    fetch('/calculate', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ equipment: equipmentList })
    })
        .then(function(response) {
            return response.json();
        })
        .then(function(results) {
            displayResults(results);
        })
        .catch(function(error) {
            window.console.error('Error:', error);
            window.alert('Помилка при розрахунках. Спробуйте ще раз.');
        });
});

function displayResults(results) {
    var resultsDiv = document.getElementById('results');
    resultsDiv.classList.remove('hidden');

    resultsDiv.innerHTML =
        '<h3>Результати розрахунків:</h3>' +
        '<h4>Для ШР1=ШР2=ШР3:</h4>' +
        '<p>1.1. Груповий коефіцієнт використання: ' + results.switchboardUtilizationFactor.toFixed(4) + '</p>' +
        '<p>1.2. Ефективна кількість ЕП: ' + results.switchboardEffectiveNumber.toFixed(0) + '</p>' +
        '<p>1.3. Розрахунковий коефіцієнт активної потужності: ' + results.switchboardActivePowerCoef.toFixed(2) + '</p>' +
        '<p>1.4. Розрахункове активне навантаження: ' + results.switchboardActivePower.toFixed(2) + ' кВт</p>' +
        '<p>1.5. Розрахункове реактивне навантаження: ' + results.switchboardReactivePower.toFixed(3) + ' квар</p>' +
        '<p>1.6. Повна потужність: ' + results.switchboardFullPower.toFixed(4) + ' кВ*А</p>' +
        '<p>1.7. Розрахунковий груповий струм: ' + results.switchboardCurrent.toFixed(2) + ' А</p>' +
        '<h4>Для цеху в цілому:</h4>' +
        '<p>1.8. Коефіцієнти використання: ' + results.workshopUtilizationFactor.toFixed(2) + '</p>' +
        '<p>1.9. Ефективна кількість ЕП: ' + results.workshopEffectiveNumber.toFixed(0) + '</p>' +
        '<p>1.10. Розрахунковий коефіцієнт активної потужності: ' + results.workshopActivePowerCoef.toFixed(1) + '</p>' +
        '<p>1.11. Розрахункове активне навантаження: ' + results.workshopActivePower.toFixed(1) + ' кВт</p>' +
        '<p>1.12. Розрахункове реактивне навантаження: ' + results.workshopReactivePower.toFixed(1) + ' квар</p>' +
        '<p>1.13. Повна потужність: ' + results.workshopFullPower.toFixed(1) + ' кВ*А</p>' +
        '<p>1.14. Розрахунковий груповий струм: ' + results.workshopCurrent.toFixed(3) + ' А</p>';
}
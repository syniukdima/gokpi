"use strict";

var calculator2Form = document.getElementById('calculator2Form');
calculator2Form.addEventListener('submit', function(e) {
    e.preventDefault();

    var formData = {
        c: parseFloat(document.getElementById('inputC').value),
        h: parseFloat(document.getElementById('inputH').value),
        o: parseFloat(document.getElementById('inputO').value),
        s: parseFloat(document.getElementById('inputS').value),
        ad: parseFloat(document.getElementById('inputA').value),
        wr: parseFloat(document.getElementById('inputW').value),
        v: parseFloat(document.getElementById('inputV').value),
        qiDaf: parseFloat(document.getElementById('inputQi_daf').value)
    };

    fetch('/api/calculate2', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData)
    })
        .then(function(response) {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(function(result) {
            var workingMassHtml =
                '<p>Вуглець (C): ' + result.workingMass.c.toFixed(2) + ' відсотків</p>' +
                '<p>Водень (H): ' + result.workingMass.h.toFixed(2) + ' відсотків</p>' +
                '<p>Кисень (O): ' + result.workingMass.o.toFixed(2) + ' відсотків</p>' +
                '<p>Сірка (S): ' + result.workingMass.s.toFixed(2) + ' відсотків</p>' +
                '<p>Зола (A): ' + result.workingMass.a.toFixed(2) + ' відсотків</p>' +
                '<p>Ванадій (V): ' + result.workingMass.v.toFixed(2) + ' мг/кг</p>';

            document.getElementById('workingMassResults').innerHTML = workingMassHtml;

            var heatingValueHtml =
                '<p>Нижча робоча теплота згоряння: ' + result.qiR.toFixed(3) + ' МДж/кг</p>';

            document.getElementById('heatingValueResult').innerHTML = heatingValueHtml;
        })
        .catch(function(error) {
            window.alert('Помилка при розрахунках. Перевірте введені дані.');
            if (window.console && window.console.error) {
                window.console.error('Error:', error);
            }
        });
});